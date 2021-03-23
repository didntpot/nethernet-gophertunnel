package packet

const (
	IDLogin = iota + 1
	IDPlayStatus
	IDServerToClientHandshake
	IDClientToServerHandshake
	IDDisconnect
	IDResourcePacksInfo
	IDResourcePackStack
	IDResourcePackClientResponse
	IDText
	IDSetTime
	IDStartGame
	IDAddPlayer
	IDAddActor
	IDRemoveActor
	IDAddItemActor
	_
	IDTakeItemActor
	IDMoveActorAbsolute
	IDMovePlayer
	IDRiderJump
	IDUpdateBlock
	IDAddPainting
	IDTickSync
	_
	IDLevelEvent
	IDBlockEvent
	IDActorEvent
	IDMobEffect
	IDUpdateAttributes
	IDInventoryTransaction
	IDMobEquipment
	IDMobArmourEquipment
	IDInteract
	IDBlockPickRequest
	IDActorPickRequest
	IDPlayerAction
	_
	IDHurtArmour
	IDSetActorData
	IDSetActorMotion
	IDSetActorLink
	IDSetHealth
	IDSetSpawnPosition
	IDAnimate
	IDRespawn
	IDContainerOpen
	IDContainerClose
	IDPlayerHotBar
	IDInventoryContent
	IDInventorySlot
	IDContainerSetData
	IDCraftingData
	IDCraftingEvent
	IDGUIDataPickItem
	IDAdventureSettings
	IDBlockActorData
	IDPlayerInput
	IDLevelChunk
	IDSetCommandsEnabled
	IDSetDifficulty
	IDChangeDimension
	IDSetPlayerGameType
	IDPlayerList
	IDSimpleEvent
	IDEvent
	IDSpawnExperienceOrb
	IDClientBoundMapItemData
	IDMapInfoRequest
	IDRequestChunkRadius
	IDChunkRadiusUpdated
	IDItemFrameDropItem
	IDGameRulesChanged
	IDCamera
	IDBossEvent
	IDShowCredits
	IDAvailableCommands
	IDCommandRequest
	IDCommandBlockUpdate
	IDCommandOutput
	IDUpdateTrade
	IDUpdateEquip
	IDResourcePackDataInfo
	IDResourcePackChunkData
	IDResourcePackChunkRequest
	IDTransfer
	IDPlaySound
	IDStopSound
	IDSetTitle
	IDAddBehaviourTree
	IDStructureBlockUpdate
	IDShowStoreOffer
	IDPurchaseReceipt
	IDPlayerSkin
	IDSubClientLogin
	IDAutomationClientConnect
	IDSetLastHurtBy
	IDBookEdit
	IDNPCRequest
	IDPhotoTransfer
	IDModalFormRequest
	IDModalFormResponse
	IDServerSettingsRequest
	IDServerSettingsResponse
	IDShowProfile
	IDSetDefaultGameType
	IDRemoveObjective
	IDSetDisplayObjective
	IDSetScore
	IDLabTable
	IDUpdateBlockSynced
	IDMoveActorDelta
	IDSetScoreboardIdentity
	IDSetLocalPlayerAsInitialised
	IDUpdateSoftEnum
	IDNetworkStackLatency
	_
	IDScriptCustomEvent
	IDSpawnParticleEffect
	IDAvailableActorIdentifiers
	_
	IDNetworkChunkPublisherUpdate
	IDBiomeDefinitionList
	IDLevelSoundEvent
	IDLevelEventGeneric
	IDLecternUpdate
	_
	IDAddEntity
	IDRemoveEntity
	IDClientCacheStatus
	IDMapCreateLockedCopy
	IDOnScreenTextureAnimation
	IDStructureTemplateDataRequest
	IDStructureTemplateDataResponse
	_
	IDClientCacheBlobStatus
	IDClientCacheMissResponse
	IDEducationSettings
	IDEmote
	IDMultiPlayerSettings
	IDSettingsCommand
	IDAnvilDamage
	IDCompletedUsingItem
	IDNetworkSettings
	IDPlayerAuthInput
	IDCreativeContent
	IDPlayerEnchantOptions
	IDItemStackRequest
	IDItemStackResponse
	IDPlayerArmourDamage
	IDCodeBuilder
	IDUpdatePlayerGameType
	IDEmoteList
	IDPositionTrackingDBServerBroadcast
	IDPositionTrackingDBClientRequest
	IDDebugInfo
	IDPacketViolationWarning
	IDMotionPredictionHints
	IDAnimateEntity
	IDCameraShake
	IDPlayerFog
	IDCorrectPlayerMovePrediction
	IDItemComponent
	IDFilterText
	IDClientBoundDebugRenderer
)
