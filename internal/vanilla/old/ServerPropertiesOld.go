package old

const (
	GameModeSurvival  GameMode = "survival"
	GameModeCreative  GameMode = "creative"
	GameModeAdventure GameMode = "adventure"

	DifficultyPeaceful Difficulty = "peaceful"
	DifficultyEasy     Difficulty = "easy"
	DifficultyNormal   Difficulty = "normal"
	DifficultyHard     Difficulty = "hard"

	PlayerPermissionLevelVisitor  PlayerPermissionLevel = "visitor"
	PlayerPermissionLevelMember   PlayerPermissionLevel = "member"
	PlayerPermissionLevelOperator PlayerPermissionLevel = "operator"

	CompressionAlgorithmZlib   CompressionAlgorithm = "zlib"
	CompressionAlgorithmSnappy CompressionAlgorithm = "snappy"

	ServerAuthoritativeMovementClientAuth           ServerAuthoritativeMovement = "client-auth"
	ServerAuthoritativeMovementServerAuth           ServerAuthoritativeMovement = "server-auth"
	ServerAuthoritativeMovementServerAuthWithRewind ServerAuthoritativeMovement = "server-auth-with-rewind"

	ChatRestrictionNone     ChatRestriction = "None"
	ChatRestrictionDropped  ChatRestriction = "Dropped"
	ChatRestrictionDisabled ChatRestriction = "Disabled"
)

type GameMode string

type Difficulty string

type PlayerPermissionLevel string

type CompressionAlgorithm string

type ServerAuthoritativeMovement string

type ChatRestriction string

type ServerPropertiesOld struct {
	// ServerName used as the server name
	// Allowed values: Any string without semicolon symbol.
	ServerName string `properties:"server-name"`

	// GameMode sets the game mode for new players.
	//Allowed values: "survival", "creative", or "adventure"
	GameMode GameMode `property:"gamemode"`

	// ForceGameMode force-gamemode=false (or force-gamemode is not defined in the server.properties)
	// prevents the server from sending to the client gamemode values other
	// than the gamemode value saved by the server during world creation
	// even if those values are set in server.properties after world creation.
	//
	// force-gamemode=true forces the server to send to the client gamemode values
	// other than the gamemode value saved by the server during world creation
	// if those values are set in server.properties after world creation.
	ForceGameMode bool `property:"force-gamemode"`

	// Difficulty sets the difficulty of the world.
	// Allowed values: "peaceful", "easy", "normal", or "hard"
	Difficulty Difficulty `property:"difficulty"`

	// AllowCheats if true then cheats like commands can be used.
	// Allowed values: "true" or "false"
	AllowCheats bool `property:"allow-cheats"`

	// MaxPlayers the maximum number of players that can play on the server.
	// Allowed values: Any positive integer
	MaxPlayers int `property:"max-players"`

	// OnlineMode if true then all connected players must be authenticated to Xbox Live.
	// Clients connecting to remote (non-LAN) servers will always require Xbox Live authentication regardless of this setting.
	// If the server accepts connections from the Internet, then it's highly recommended to enable online-mode.
	// Allowed values: "true" or "false"
	OnlineMode bool `property:"online-mode"`

	// AllowList if true then all connected players must be listed in the separate allowlist.json file.
	// Allowed values: "true" or "false"
	AllowList bool `property:"allow-list"`

	// ServerPort Which IPv4 port the server should listen to.
	// Allowed values: Integers in the range [1, 65535]
	ServerPort int `property:"server-port"`

	// ServerPortV6 Which IPv6 port the server should listen to.
	// Allowed values: Integers in the range [1, 65535]
	ServerPortV6 int `property:"server-portv6"`

	// EnableLanVisibility listen and respond to clients that are looking for servers on the LAN. This will cause the server
	// to bind to the default ports (19132, 19133) even when `server-port` and `server-portv6`
	// have non-default values. Consider turning this off if LAN discovery is not desirable, or when
	// running multiple servers on the same host may lead to port conflicts.
	// Allowed values: "true" or "false"
	EnableLanVisibility bool `property:"enable-lan-visibility"`

	// ViewDistance the maximum allowed view distance in number of chunks.
	// Allowed values: Positive integer equal to 5 or greater.
	ViewDistance int `property:"view-distance"`

	// TickDistance the world will be ticked this many chunks away from any player.
	// Allowed values: Integers in the range [4, 12]
	TickDistance int `property:"tick-distance"`

	// PlayerIdleTimeout after a player has idled for this many minutes they will be kicked. If set to 0 then players can idle indefinitely.
	// Allowed values: Any non-negative integer.
	PlayerIdleTimeout int `property:"player-idle-timeout"`

	// MaxThreads maximum number of threads the server will try to use. If set to 0 or removed then it will use as many as possible.
	// Allowed values: Any positive integer.
	MaxThreads int `property:"max-threads"`

	// LevelName allowed values: Any string without semicolon symbol or symbols illegal for file name: /\n\r\t\f`?*\\<>|\":
	LevelName string `property:"level-name"`

	// LevelSeed use to randomize the world
	// Allowed values: Any string
	LevelSeed string `property:"level-seed"`

	// DefaultPlayerPermissionLevel permission level for new players joining for the first time.
	// Allowed values: "visitor", "member", "operator"
	DefaultPlayerPermissionLevel PlayerPermissionLevel `property:"default-player-permission-level"`

	// TexturepackRequired force clients to use texture packs in the current world
	// Allowed values: "true" or "false"
	TexturepackRequired bool `property:"texturepack-required"`

	// ContentLogFileEnabled enables logging content errors to a file
	// Allowed values: "true" or "false"
	ContentLogFileEnabled bool `property:"content-log-file-enabled"`

	// CompressionThreshold determines the smallest size of raw network payload to compress
	// Allowed values: 0-65535
	CompressionThreshold int `property:"compression-threshold"`

	// CompressionAlgorithm determines the compression algorithm to use for networking
	// Allowed values: "zlib", "snappy"
	CompressionAlgorithm CompressionAlgorithm `property:"compression-algorithm"`

	// ServerAuthoritativeMovement allowed values: "client-auth", "server-auth", "server-auth-with-rewind"
	// Enables server authoritative movement. If "server-auth", the server will replay local user input on
	// the server and send down corrections when the client's position doesn't match the server's.
	// If "server-auth-with-rewind" is enabled and the server sends a correction, the clients will be instructed
	// to rewind time back to the correction time, apply the correction, then replay all the player's inputs since then. This results in smoother and more frequent corrections.
	// Corrections will only happen if correct-player-movement is set to true.
	ServerAuthoritativeMovement ServerAuthoritativeMovement `property:"server-authoritative-movement"`

	// PlayerMovementScoreThreshold the number of incongruent time intervals needed before abnormal behavior is reported.
	// Disabled by server-authoritative-movement.
	PlayerMovementScoreThreshold int `property:"player-movement-score-threshold"`

	// PlayerMovementActionDirectionThreshold the amount that the player's attack direction and look direction can differ.
	// Allowed values: Any value in the range of [0, 1] where 1 means that the
	// direction of the players view and the direction the player is attacking
	// must match exactly and a value of 0 means that the two directions can
	// differ by up to and including 90 degrees.
	PlayerMovementActionDirectionThreshold float64 `property:"player-movement-action-direction-threshold"`

	// PlayerMovementDistanceThreshold the difference between server and client positions that needs to be exceeded before abnormal behavior is detected.
	// Disabled by server-authoritative-movement.
	PlayerMovementDistanceThreshold float64 `property:"player-movement-distance-threshold"`

	// PlayerMovementDurationThresholdInMS the duration of time the server and client positions can be out of sync (as defined by player-movement-distance-threshold)
	// before the abnormal movement score is incremented. This value is defined in milliseconds.
	// Disabled by server-authoritative-movement.
	PlayerMovementDurationThresholdInMS int `property:"player-movement-duration-threshold-in-ms"`

	// CorrectPlayerMovement if true, the client position will get corrected to the server position if the movement score exceeds the threshold.
	CorrectPlayerMovement bool `property:"correct-player-movement"`

	// ServerAuthoritativeBlockBreaking if true, the server will compute block mining operations in sync with the client so it can verify that the client should be able to break blocks when it thinks it can.
	ServerAuthoritativeBlockBreaking bool `property:"server-authoritative-block-breaking"`

	// ChatRestriction allowed values: "None", "Dropped", "Disabled"
	// This represents the level of restriction applied to the chat for each player that joins the server.
	// "None" is the default and represents regular free chat.
	// "Dropped" means the chat messages are dropped and never sent to any client. Players receive a message to let them know the feature is disabled.
	// "Disabled" means that unless the player is an operator, the chat UI does not even appear. No information is displayed to the player.
	ChatRestriction ChatRestriction `property:"chat-restriction"`

	// DisablePlayerInteraction if true, the server will inform clients that they should ignore other players when interacting with the world. This is not server authoritative.
	DisablePlayerInteraction bool `property:"disable-player-interaction"`

	// ClientSideChunkGenerationEnabled if true, the server will inform clients that they have the ability to generate visual level chunks outside of player interaction distances.
	ClientSideChunkGenerationEnabled bool `property:"client-side-chunk-generation-enabled"`

	// BlockNetworkIDsAreHashes If true, the server will send hashed block network ID's instead of id's that start from 0 and go up.  These id's are stable and won't change regardless of other block changes.
	BlockNetworkIDsAreHashes bool `property:"block-network-ids-are-hashes"`

	// DisablePersona internal Use Only
	DisablePersona bool `property:"disable-persona=false"`

	// DisableCustomSkins if true, disable players customized skins that were customized outside the Minecraft store assets or in game assets.  This is used to disable possibly offensive custom skins players make.
	DisableCustomSkins bool `property:"disable-custom-skins"`

	// ServerBuildRadiusRatio allowed values: "Disabled" or any value in range [0.0, 1.0]
	// If "Disabled" the server will dynamically calculate how much of the player's view it will generate, assigning the rest to the client to build.
	// Otherwise, from the overridden ratio tell the server how much of the player's view to generate, disregarding client hardware capability.
	// Only valid if client-side-chunk-generation-enabled is enabled
	ServerBuildRadiusRatio string `property:"server-build-radius-ratio=Disabled"`

	RCONPort           int    `property:"rcon.port"`
	BroadcastRCONToOps bool   `property:"broadcast-rcon-to-ops"`
	EnableRCON         bool   `property:"enable-rcon=false"`
	RCONPassword       string `property:"rcon.password"`
}

type ConfigureServer func(cfg *ServerPropertiesOld) *ServerPropertiesOld

func WithServerPort(p int) ConfigureServer {
	return func(cfg *ServerPropertiesOld) *ServerPropertiesOld {
		cfg.ServerPort = p
		return cfg
	}
}

func WithRCONEnabled(e bool) ConfigureServer {
	return func(cfg *ServerPropertiesOld) *ServerPropertiesOld {
		cfg.EnableRCON = e
		return cfg
	}
}

func WithRCONPass(p string) ConfigureServer {
	return func(cfg *ServerPropertiesOld) *ServerPropertiesOld {
		cfg.RCONPassword = p
		return cfg
	}
}
