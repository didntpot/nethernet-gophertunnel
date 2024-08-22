package minecraft_test

import (
	"context"
	rand2 "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/auth"
	"github.com/sandertv/gophertunnel/minecraft/franchise"
	"github.com/sandertv/gophertunnel/minecraft/franchise/signaling"
	"github.com/sandertv/gophertunnel/minecraft/nethernet"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/room"
	"github.com/sandertv/gophertunnel/xsapi/xal"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/xsapi"
	"github.com/sandertv/gophertunnel/xsapi/mpsd"
	"golang.org/x/oauth2"
	"math/rand"
	"strings"
	"testing"
)

// TestWorldListen demonstrates a world displayed in the friend list.
func TestWorldListen(t *testing.T) {
	discovery, err := franchise.Discover(protocol.CurrentVersion)
	if err != nil {
		t.Fatalf("error retrieving discovery: %s", err)
	}
	a := new(franchise.AuthorizationEnvironment)
	if err := discovery.Environment(a, franchise.EnvironmentTypeProduction); err != nil {
		t.Fatalf("error reading environment for authorization: %s", err)
	}
	s := new(signaling.Environment)
	if err := discovery.Environment(s, franchise.EnvironmentTypeProduction); err != nil {
		t.Fatalf("error reading environment for signaling: %s", err)
	}

	tok, err := readToken("franchise/internal/test/auth.tok", auth.TokenSource)
	if err != nil {
		t.Fatalf("error reading token: %s", err)
	}
	src := auth.RefreshTokenSource(tok)

	refresh, cancel := context.WithCancel(context.Background())
	defer cancel()
	prov := franchise.PlayFabXBLIdentityProvider{
		Environment: a,
		TokenSource: xal.RefreshTokenSourceContext(refresh, src, "http://playfab.xboxlive.com/"),
	}

	d := signaling.Dialer{
		NetworkID: rand.Uint64(),
	}

	dial, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	conn, err := d.DialContext(dial, prov, s)
	if err != nil {
		t.Fatalf("error dialing signaling: %s", err)
	}
	t.Cleanup(func() {
		if err := conn.Close(); err != nil {
			t.Fatalf("error closing signaling: %s", err)
		}
	})

	// A token source that refreshes a token used for generic Xbox Live services.
	x := xal.RefreshTokenSourceContext(refresh, src, "http://xboxlive.com")
	xt, err := x.Token()
	if err != nil {
		t.Fatalf("error refreshing xbox live token: %s", err)
	}
	claimer, ok := xt.(xsapi.DisplayClaimer)
	if !ok {
		t.Fatalf("xbox live token %T does not implement xsapi.DisplayClaimer", xt)
	}
	displayClaims := claimer.DisplayClaims()

	// The name of the session being published. This seems always to be generated
	// randomly, referenced as "GUID" of the session.
	name := strings.ToUpper(uuid.NewString())

	levelID := make([]byte, 8)
	_, _ = rand2.Read(levelID)

	custom, err := json.Marshal(room.Status{
		Joinability: room.JoinabilityJoinableByFriends,
		HostName:    displayClaims.GamerTag,
		OwnerID:     displayClaims.XUID,
		RakNetGUID:  "",
		// This is displayed as the suffix of the world name.
		Version:   protocol.CurrentVersion,
		LevelID:   base64.StdEncoding.EncodeToString(levelID),
		WorldName: "TestWorldListen: " + name,
		WorldType: room.WorldTypeCreative,
		// The game seems checking this field before joining a session, causes
		// RequestNetworkSettings packet not being even sent to the remote host.
		Protocol:                protocol.CurrentProtocol,
		MemberCount:             1,
		MaxMemberCount:          8,
		BroadcastSetting:        room.BroadcastSettingFriendsOfFriends,
		LanGame:                 true,
		IsEditorWorld:           false,
		TransportLayer:          2,
		WebRTCNetworkID:         d.NetworkID,
		OnlineCrossPlatformGame: true,
		CrossPlayDisabled:       false,
		TitleID:                 0,
		SupportedConnections: []room.Connection{
			{
				ConnectionType:  3, // WebSocketsWebRTCSignaling
				HostIPAddress:   "",
				HostPort:        0,
				NetherNetID:     d.NetworkID,
				WebRTCNetworkID: d.NetworkID,
				RakNetGUID:      "UNASSIGNED_RAKNET_GUID",
			},
		},
	})
	if err != nil {
		t.Fatalf("error encoding custom properties: %s", err)
	}
	cfg := mpsd.PublishConfig{
		Description: &mpsd.SessionDescription{
			Properties: &mpsd.SessionProperties{
				System: &mpsd.SessionPropertiesSystem{
					JoinRestriction: mpsd.SessionRestrictionFollowed,
					ReadRestriction: mpsd.SessionRestrictionFollowed,
				},
				Custom: custom,
			},
		},
	}

	publish, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	session, err := cfg.PublishContext(publish, x, mpsd.SessionReference{
		ServiceConfigID: serviceConfigID,
		TemplateName:    "MinecraftLobby",
		Name:            name,
	})
	if err != nil {
		t.Fatalf("error publishing session: %s", err)
	}
	t.Cleanup(func() {
		if err := session.Close(); err != nil {
			t.Fatalf("error closing session: %s", err)
		}
	})

	t.Logf("Session Name: %q", name)
	t.Logf("Network ID: %d", d.NetworkID)

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	minecraft.RegisterNetwork("nethernet", &nethernet.Network{
		Signaling: conn,
	})

	l, err := minecraft.Listen("nethernet", nethernet.NetworkAddress(d.NetworkID))
	if err != nil {
		t.Fatalf("error listening: %s", err)
	}
	t.Cleanup(func() {
		if err := l.Close(); err != nil {
			t.Fatalf("error closing listener: %s", err)
		}
	})

	for {
		netConn, err := l.Accept()
		if err != nil {
			if !errors.Is(err, net.ErrClosed) {
				t.Fatalf("error accepting conn: %s", err)
			}
		}
		c := netConn.(*minecraft.Conn)
		if err := c.StartGame(minecraft.GameData{
			WorldName:         "NetherNet",
			WorldSeed:         0,
			Difficulty:        0,
			EntityUniqueID:    rand.Int63(),
			EntityRuntimeID:   rand.Uint64(),
			PlayerGameMode:    1,
			PlayerPosition:    mgl32.Vec3{},
			WorldSpawn:        protocol.BlockPos{},
			WorldGameMode:     1,
			Time:              rand.Int63(),
			PlayerPermissions: 2,
			// Allow inviting player into the world.
			GamePublishSetting: 3,
		}); err != nil {
			t.Fatalf("error starting game: %s", err)
		}
	}
}

var serviceConfigID = uuid.MustParse("4fc10100-5f7a-4470-899b-280835760c07")

// TestWorldDial connects to a world. Before running the test, you need to capture the network ID of
// the world to join, and fill in the constant below.
func TestWorldDial(t *testing.T) {
	// TODO: Implement looking up sessions and find a network ID from the response.
	const remoteNetworkID = 9511338490860978050

	discovery, err := franchise.Discover(protocol.CurrentVersion)
	if err != nil {
		t.Fatalf("error retrieving discovery: %s", err)
	}
	a := new(franchise.AuthorizationEnvironment)
	if err := discovery.Environment(a, franchise.EnvironmentTypeProduction); err != nil {
		t.Fatalf("error reading environment for authorization: %s", err)
	}
	s := new(signaling.Environment)
	if err := discovery.Environment(s, franchise.EnvironmentTypeProduction); err != nil {
		t.Fatalf("error reading environment for signaling: %s", err)
	}

	tok, err := readToken("franchise/internal/test/auth.tok", auth.TokenSource)
	if err != nil {
		t.Fatalf("error reading token: %s", err)
	}
	src := auth.RefreshTokenSource(tok)

	refresh, cancel := context.WithCancel(context.Background())
	defer cancel()
	prov := franchise.PlayFabXBLIdentityProvider{
		Environment: a,
		TokenSource: xal.RefreshTokenSourceContext(refresh, src, "http://playfab.xboxlive.com/"),
	}

	d := signaling.Dialer{
		NetworkID: rand.Uint64(),
	}

	dial, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	sig, err := d.DialContext(dial, prov, s)
	if err != nil {
		t.Fatalf("error dialing signaling: %s", err)
	}
	t.Cleanup(func() {
		if err := sig.Close(); err != nil {
			t.Fatalf("error closing signaling: %s", err)
		}
	})

	// TODO: Implement joining a session.
	// A token source that refreshes a token used for generic Xbox Live services.
	//x := xal.RefreshTokenSourceContext(refresh, src, "http://xboxlive.com")

	t.Logf("Network ID: %d", d.NetworkID)

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	minecraft.RegisterNetwork("nethernet", &nethernet.Network{
		Signaling: sig,
	})

	conn, err := minecraft.Dialer{
		TokenSource: src,
	}.DialTimeout("nethernet", nethernet.NetworkAddress(remoteNetworkID), time.Second*15)
	if err != nil {
		t.Fatalf("error dialing: %s", err)
	}
	t.Cleanup(func() {
		if err := conn.Close(); err != nil {
			t.Fatalf("error closing session: %s", err)
		}
	})

	if err := conn.DoSpawn(); err != nil {
		t.Fatalf("error spawning: %s", err)
	}
	if err := conn.WritePacket(&packet.Text{
		TextType:   packet.TextTypeChat,
		SourceName: conn.IdentityData().DisplayName,
		Message:    "Successful",
		XUID:       conn.IdentityData().XUID,
	}); err != nil {
		t.Fatalf("error writing packet: %s", err)
	}

	// Try decoding deferred packets received from the connection.
	go func() {
		for {
			pk, err := conn.ReadPacket()
			if err != nil {
				if !strings.Contains(err.Error(), "use of closed network connection") {
					t.Errorf("error reading packet: %s", err)
				}
				return
			}
			_ = pk
		}
	}()

	time.Sleep(time.Second * 15)
}

func readToken(path string, src oauth2.TokenSource) (t *oauth2.Token, err error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t, err = src.Token()
		if err != nil {
			return nil, fmt.Errorf("obtain token: %w", err)
		}
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		if err := json.NewEncoder(f).Encode(t); err != nil {
			return nil, fmt.Errorf("encode: %w", err)
		}
		return t, nil
	} else if err != nil {
		return nil, fmt.Errorf("stat: %w", err)
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&t); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return t, nil
}