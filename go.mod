module github.com/sandertv/gophertunnel

go 1.23.0

toolchain go1.23.4

require (
	github.com/coder/websocket v1.8.12
	github.com/df-mc/go-playfab v0.0.0-20240902102459-2f8b5cd02173
	github.com/df-mc/go-xsapi v0.0.0-20240902102602-e7c4bffb955f
	github.com/go-gl/mathgl v1.1.0
	github.com/go-jose/go-jose/v3 v3.0.3
	github.com/golang/snappy v0.0.4
	github.com/google/uuid v1.6.0
	github.com/klauspost/compress v1.17.9
	github.com/muhammadmuzzammil1998/jsonc v1.0.0
	github.com/pelletier/go-toml v1.9.5
	github.com/sandertv/go-raknet v1.14.1
	golang.org/x/net v0.33.0
	golang.org/x/oauth2 v0.21.0
	golang.org/x/text v0.21.0
)

require (
	github.com/stretchr/testify v1.9.0 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/image v0.17.0 // indirect
)

replace (
	github.com/df-mc/go-nethernet => github.com/lactyy/go-nethernet v0.0.0-20240918151603-8274a4680204
	github.com/df-mc/go-playfab => github.com/lactyy/go-playfab v0.0.0-20240911042657-037f6afe426f
	github.com/df-mc/go-xsapi => github.com/lactyy/go-xsapi v0.0.0-20240911052022-1b9dffef64ab
)
