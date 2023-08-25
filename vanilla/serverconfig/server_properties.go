package serverconfig

import (
	"fmt"
	"github.com/eldius/properties"
)

const (
	// GameDifficultyPeaceful represents a peaceful world
	GameDifficultyPeaceful GameDifficulty = "peaceful"
	// GameDifficultyEasy represents a easy world
	GameDifficultyEasy GameDifficulty = "easy"
	// GameDifficultyNormal represents a normal world
	GameDifficultyNormal GameDifficulty = "normal"
	// GameDifficultyHard represents a hard world
	GameDifficultyHard GameDifficulty = "hard"

	// GameModeSurvival represents a survival (0) game mode
	GameModeSurvival GameMode = "survival"
	// GameModeCreative represents a creative (1) game mode
	GameModeCreative GameMode = "creative"
	// GameModeAdventure represents a adventure (2) game mode
	GameModeAdventure GameMode = "adventure"
	// GameModeSpectator represents a spectator (3) game mode
	GameModeSpectator GameMode = "spectator"

	// LevelTypeNormal Standard world with hills, valleys, water, etc.
	LevelTypeNormal LevelType = "minecraft:normal"
	// LevelTypeFlat A flat world with no features, can be modified with generator-settings.
	LevelTypeFlat LevelType = "minecraft:flat"
	// LevelTypeLargeBiomes Same as default but all biomes are larger.
	LevelTypeLargeBiomes LevelType = "minecraft:large_biomes"
	// LevelTypeAmplified Same as default but world-generation height limit is increased.
	LevelTypeAmplified LevelType = "minecraft:amplified"
	// LevelTypeSingleBiomeSurface A buffet world which the entire overworld consists of one biome, can be modified with generator-settings.
	LevelTypeSingleBiomeSurface LevelType = "minecraft:single_biome_surface"
	// LevelTypeBuffet Only for 1.15 or before. Same as default unless generator-settings is set.
	LevelTypeBuffet LevelType = "buffet"
	// LevelTypeDefault11 Only for 1.15 or before. Same as default, but counted as a different world type.
	LevelTypeDefault11 LevelType = "default_1_1"
	// LevelTypeCustomized Only for 1.15 or before. After 1.13, this value is no different than default, but in 1.12 and before, it could be used to create a completely custom world.
	LevelTypeCustomized LevelType = "customized"
)

// GameDifficulty is the game difficulty configuration
type GameDifficulty string

// GameMode is the game mode configuration
type GameMode string

// LevelType is the level type configuration
type LevelType string

// ServerConfig is configuration function to help setting up server
type ServerConfig func(s *ServerProperties) *ServerProperties

// ServerProperties Represents the server.properties content
// references from: https://minecraft.fandom.com/wiki/Server.properties
type ServerProperties struct {
	// AllowFlight
	// Allows users to use flight on the server while in Survival mode, if they have a mod that provides flight installed.
	// With allow-flight enabled, griefers may become more common, because it makes their work easier. In Creative mode, this has no effect.
	// allow-flight
	//
	// type: boolean
	// default: false
	// false - Flight is not allowed (players in air for at least 5 seconds get kicked).
	// true - Flight is allowed, and used if the player has a fly mod installed.
	AllowFlight bool `properties:"allow-flight"`

	// AllowNether
	// Allows players to travel to the Nether.
	// allow-nether
	//
	// type: boolean
	// default: true
	// false - Nether portals do not work.
	// true - The server allows portals to send players to the Nether.
	AllowNether bool `properties:"allow-nether"`

	// BroadcastToOps
	// Send console command outputs to all online operators.
	// broadcast-console-to-ops
	// type: boolean
	// default: true
	BroadcastToOps bool `properties:"broadcast-console-to-ops"`

	// BroadcastRconToOps
	// Send rcon console command outputs to all online operators.
	// broadcast-rcon-to-ops
	// type: boolean
	// default: true
	BroadcastRconToOps bool `properties:"broadcast-rcon-to-ops"`

	// Difficulty
	// Defines the difficulty (such as damage dealt by mobs and the way hunger and poison affects players) of the server.
	// difficulty
	// type: string
	// default: easy
	//
	// If a legacy difficulty number is specified, it is silently converted to a difficulty name.
	//
	// peaceful (0)
	// easy (1)
	// normal (2)
	// hard (3)
	Difficulty GameDifficulty `properties:"difficulty"`

	// EnableCommandBlock
	// Enables command blocks.
	// enable-command-block
	// type: boolean
	// default: false
	EnableCommandBlock bool `properties:"enable-command-block"`

	// EnableJmxMonitoring
	// Exposes an MBean with the Object name net.minecraft.server:type=Server and two attributes averageTickTime and tickTimes exposing the tick times in milliseconds.
	// enable-jmx-monitoring
	// type: boolean
	// default: false
	EnableJmxMonitoring bool `properties:"enable-jmx-monitoring"`

	// EnableRcon
	// Enables remote access to the server console.
	// In order for enabling JMX on the Java runtime you also need to add a couple of JVM flags to the startup as documented here.
	// enable-rcon
	// type: boolean
	// default: false
	EnableRcon bool `properties:"enable-rcon"`

	// EnableStatus
	// Makes the server appear as "online" on the server list.
	//It's not recommended to expose RCON to the Internet, because RCON protocol transfers everything without encryption. Everything (including RCON password) communicated between the RCON server and client can be leaked to someone listening in on your connection.
	// enable-status
	// type: boolean
	// default: true
	EnableStatus bool `properties:"enable-status"`

	// EnableQuery
	// Enables GameSpy4 protocol server listener. Used to get information about server.
	// If set to false, it will suppress replies from clients. This means it will appear as offline, but will still accept connections.
	// enable-query
	// type: boolean
	// default: false
	EnableQuery bool `properties:"enable-query"`

	// EnforceSecureProfile
	// If set to true, players without a Mojang-signed public key will not be able to connect to the server.
	// enforce-secure-profile
	// type: boolean
	// default: true
	EnforceSecureProfile bool `properties:"enforce-secure-profile"`

	// EnforceWhitelist
	// Enforces the whitelist on the server.
	// enforce-whitelist
	// type: boolean
	// default: false
	EnforceWhitelist bool `properties:"enforce-whitelist"`

	// EntityBroadcastRangePercentage
	// Controls how close entities need to be before being sent to clients. Higher values means they'll be rendered from farther away, potentially causing more lag. This is expressed the percentage of the default value. For example, setting to 50 will make it half as usual. This mimics the function on the client video settings (not unlike Render Distance, which the client can customize so long as it's under the server's setting).
	// When this option is enabled, users who are not present on the whitelist (if it's enabled) get kicked from the server after the server reloads the whitelist file.
	// false - No user gets kicked if not on the whitelist.
	// true - Online users not on the whitelist get kicked.
	//
	// entity-broadcast-range-percentage
	// type: integer (10-1000)
	// default: 100
	EntityBroadcastRangePercentage int `properties:"entity-broadcast-range-percentage"`

	// ForceGameMode
	// Force players to join in the default game mode.
	// Sets the default permission level for [functions](https://minecraft.fandom.com/wiki/Function_(Java_Edition).
	// false - Players join in the gamemode they left in.
	// true - Players always join in the default gamemode.
	//
	// force-gamemode
	// type: boolean
	// default: false
	ForceGameMode bool `properties:"force-gamemode"`

	// FunctionPermissionLevel
	// See [permission level](https://minecraft.fandom.com/wiki/Permission_level) for the details on the 4 levels.
	// function-permission-level
	// type: integer (1-4)
	// default: 2
	FunctionPermissionLevel int `properties:"function-permission-level"`

	// GameMode
	// Defines the mode of gameplay.
	// gamemode
	// type: string
	// default: survival
	//
	// If a legacy gamemode number is specified, it is silently converted to a gamemode name.
	//
	//survival (0)
	//creative (1)
	//adventure (2)
	//spectator (3)
	GameMode GameMode `properties:"gamemode"`

	// GenerateStructures
	// Defines whether structures (such as villages) can be generated.
	// generate-structures
	// type: boolean
	// default: true
	//
	// false - Structures are not generated in new chunks.
	// true - Structures are generated in new chunks.
	//
	// Note: Dungeons still generate if this is set to false.
	GenerateStructures bool `properties:"generate-structures"`

	// GeneratorSettings
	// The settings used to customize world generation. Follow its format and write the corresponding JSON string. Remember to escape all : with \:.
	// generator-settings
	// type: string
	// default: {}
	GeneratorSettings string `properties:"generator-settings"`

	// Hardcore
	// If set to true, server difficulty is ignored and set to hard and players are set to spectator mode if they die.
	// hardcore
	// type: boolean
	// default: false
	Hardcore bool `properties:"hardcore"`

	// HideOnlinePlayers
	// If set to true, a player list is not sent on status requests.
	// hide-online-players
	// type: boolean
	// default: false
	HideOnlinePlayers bool `properties:"hide-online-players"`

	// InitialDisabledPacks
	// Comma-separated list of datapacks to not be auto-enabled on world creation.
	// initial-disabled-packs
	// type: string
	// default: blank
	InitialDisabledPacks string `properties:"initial-disabled-packs"`

	// InitialEnabledPacks
	// Comma-separated list of datapacks to be enabled during world creation. Feature packs need to be explicitly enabled.
	// initial-enabled-packs
	// type: string
	// default: vanilla
	InitialEnabledPacks string `properties:"initial-enabled-packs"`

	// LevelName
	// The "level-name" value is used as the world name and its folder name. The player may also copy their saved game folder here, and change the name to the same as that folder's to load it instead.
	// level-name
	// type: string
	// default: world
	//
	// Characters such as ' (apostrophe) may need to be escaped by adding a backslash before them.
	LevelName string `properties:"level-name"`

	// LevelSeed
	// Sets a world seed for the player's world, as in Singleplayer. The world generates with a random seed if left blank.
	// level-seed
	// type: string
	// default: blank
	//
	//Some examples are: minecraft, 404, 1a2b3c.
	LevelSeed string `properties:"level-seed"`

	// LevelType
	// Determines the world preset that is generated.
	// level-type
	// type: string
	// default: minecraft:normal
	//
	// Escaping ":" is required when using a world preset ID, and the vanilla world preset ID's namespace (minecraft:) can be omitted.
	//
	// minecraft:normal - Standard world with hills, valleys, water, etc.
	// minecraft:flat - A flat world with no features, can be modified with generator-settings.
	// minecraft:large_biomes - Same as default but all biomes are larger.
	// minecraft:amplified - Same as default but world-generation height limit is increased.
	// minecraft:single_biome_surface - A buffet world which the entire overworld consists of one biome, can be modified with generator-settings.
	// buffet - Only for 1.15 or before. Same as default unless generator-settings is set.
	// default_1_1 - Only for 1.15 or before. Same as default, but counted as a different world type.
	// customized - Only for 1.15 or before. After 1.13, this value is no different than default, but in 1.12 and before, it could be used to create a completely custom world.
	LevelType LevelType `properties:"level-type"`

	// Limiting the amount of consecutive neighbor updates before skipping additional ones. Negative values remove the limit.
	// max-chained-neighbor-updates
	// type: integer
	// default: 1000000
	MaxChainedNeighborUpdates int `properties:"max-chained-neighbor-updates"`

	// The maximum number of players that can play on the server at the same time. Note that more players on the server consume more resources. Note also, op player connections are not supposed to count against the max players, but ops currently cannot join a full server. However, this can be changed by going to the file called ops.json in the player's server directory, opening it, finding the op that the player wants to change, and changing the setting called bypassesPlayerLimit to true (the default is false). This means that that op does not have to wait for a player to leave in order to join. Extremely large values for this field result in the client-side user list being broken.
	// max-players
	// type: integer (0-(2^31 - 1))
	// default: 20
	MaxPlayers int `properties:"max-players"`

	// The maximum number of milliseconds a single tick may take before the server watchdog stops the server with the message, A single server tick took 60.00 seconds (should be max 0.05); Considering it to be crashed, server will forcibly shutdown. Once this criterion is met, it calls System.exit(1).
	// max-tick-time
	// type: integer (-1 or 0–(2^63 - 1))
	// default: 60000
	//
	// -1 - disable watchdog entirely (this disable option was added in 14w32a)
	MaxTickTime int64 `properties:"max-tick-time"`

	// This sets the maximum possible size in blocks, expressed as a radius, that the world border can obtain. Setting the world border bigger causes the commands to complete successfully but the actual border does not move past this block limit. Setting the max-world-size higher than the default doesn't appear to do anything.
	// max-world-size
	// type: integer (1-29999984)
	// default: 29999984
	//
	// Examples:
	//
	// Setting max-world-size to 1000 allows the player to have a 2000×2000 world border.
	// Setting max-world-size to 4000 gives the player an 8000×8000 world border.
	MaxWorldSize int64 `properties:"max-world-size"`

	// This is the message that is displayed in the server list of the client, below the name.
	// motd
	// type: string
	// default: A Minecraft Server
	//
	// The MOTD supports color and formatting codes.
	// The MOTD supports special characters, such as "♥". However, such characters must be converted to escaped Unicode form. An online converter can be found here.
	// If the MOTD is over 59 characters, the server list may report a communication error.
	Motd string `properties:"motd"`

	// By default it allows packets that are n-1 bytes big to go normally, but a packet of n bytes or more gets compressed down. So, a lower number means more compression but compressing small amounts of bytes might actually end up with a larger result than what went in.
	// network-compression-threshold
	// type: integer
	// default: 256
	//
	// -1 - disable compression entirely
	// 0 - compress everything
	//
	// Note: The Ethernet spec requires that packets less than 64 bytes become padded to 64 bytes. Thus, setting a value lower than 64 may not be beneficial. It is also not recommended to exceed the MTU, typically 1500 bytes.
	NetworkCompressionThreshold int `properties:"network-compression-threshold"`

	// Server checks connecting players against Minecraft account database. Set this to false only if the player's server is not connected to the Internet. Hackers with fake accounts can connect if this is set to false! If minecraft.net is down or inaccessible, no players can connect if this is set to true. Setting this variable to off purposely is called "cracking" a server, and servers that are present with online mode off are called "cracked" servers, allowing players with unlicensed copies of Minecraft to join.
	// online-mode
	// type: boolean
	// default: true
	//
	//true - Enabled. The server assumes it has an Internet connection and checks every connecting player.
	//false - Disabled. The server does not attempt to check connecting players.
	OnlineMode bool `properties:"online-mode"`

	// Sets the default permission level for ops when using /op.
	// op-permission-level
	// type: integer (0-4)
	// default: 4
	OpPermissionLevel int `properties:"op-permission-level"`

	// If non-zero, players are kicked from the server if they are idle for more than that many minutes.
	// player-idle-timeout
	// type: integer
	// default: 0
	//
	// Note: Idle time is reset when the server receives one of the following packets:
	//
	// Click Window
	// Enchant Item
	// Update Sign
	// Player Digging
	// Player Block Placement
	// Held Item Change
	// Animation (swing arm)
	// Entity Action
	// Client Status
	// Chat Message
	// Use Entity
	PlayerIdleTimeout int64 `properties:"player-idle-timeout"`

	// If the ISP/AS sent from the server is different from the one from Mojang Studios' authentication server, the player is kicked.
	// prevent-proxy-connections
	// type: boolean
	// default: false
	PreventProxyConnections bool `properties:"prevent-proxy-connections"`

	// If set to true, chat preview will be enabled.
	// previews-chat
	// type: boolean
	// default: false
	//
	// true - Enabled. When enabled, a server-controlled preview appears above the chat edit box, showing how the message will look when sent.
	// false - Disabled.
	PreviewChat bool `properties:"previews-chat"`

	// Enable PvP on the server. Players shooting themselves with arrows receive damage only if PvP is enabled.
	// pvp
	// type: boolean
	// default: true
	//
	// true - Players can kill each other.
	// false - Players cannot kill other players (also known as Player versus Environment (PvE)).
	//
	// Note: Indirect damage sources spawned by players (such as lava, fire, TNT and to some extent water, sand and gravel) still deal damage to other players.
	Pvp bool `properties:"pvp"`

	// Sets the port for the query server (see enable-query).
	// query.port
	// type: integer (1-(2^16 - 2))
	// default: 25565
	QueryPort int `properties:"query.port"`

	// Sets the maximum amount of packets a user can send before getting kicked. Setting to 0 disables this feature.
	// rate-limit
	// type: integer
	// default: 0
	RateLimit int `properties:"rate-limit"`

	// Sets the password for RCON: a remote console protocol that can allow other applications to connect and interact with a Minecraft server over the internet.
	// rcon.password
	// type: string
	// default: blank
	RconPassword string `properties:"rcon.password"`

	// Sets the RCON network port.
	// rcon.port
	// type: integer (1-(2^16 - 2))
	// default: 25575
	RconPort int `properties:"rcon.port"`

	// Optional URI to a resource pack. The player may choose to use it.
	// resource-pack
	// type: string
	// default: blank
	//
	// Note that (in some versions before 1.15.2), the ":" and "=" characters need to be escaped with a backslash (\), e.g. http\://somedomain.com/somepack.zip?someparam\=somevalue
	//
	// The resource pack may not have a larger file size than 250 MiB (Before 1.18: 100 MiB (≈ 100.8 MB)) (Before 1.15: 50 MiB (≈ 50.4 MB)). Note that download success or failure is logged by the client, and not by the server.
	ResourcePack string `properties:"resource-pack"`

	// Optional, adds a custom message to be shown on resource pack prompt when require-resource-pack is used.
	// resource-pack-prompt
	// type: string
	// default: blank
	//
	// Expects chat component syntax, can contain multiple lines.
	ResourcePackPrompt string `properties:"resource-pack-prompt"`

	// Optional SHA-1 digest of the resource pack, in lowercase hexadecimal. It is recommended to specify this, because it is used to verify the integrity of the resource pack.
	// resource-pack-sha1
	// type: string
	// default: blank
	//
	// Note: If the resource pack is any different, a yellow message "Invalid sha1 for resource-pack-sha1" appears in the console when the server starts. Due to the nature of hash functions, errors have a tiny probability of occurring, so this consequence has no effect.
	ResourcePackSHA1 string `properties:"resource-pack-sha1"`

	// When this option is enabled (set to true), players will be prompted for a response and will be disconnected if they decline the required pack.
	// require-resource-pack
	// type: boolean
	// default: false
	RequireResourcePack bool `properties:"require-resource-pack"`

	// The player should set this if they want the server to bind to a particular IP. It is strongly recommended that the player leaves server-ip blank.
	// server-ip
	// type: string
	// default: blank
	//
	// Set to blank, or the IP the player want their server to run (listen) on.
	ServerIP string `properties: "server-ip"`

	// Changes the port the server is hosting (listening) on. This port must be forwarded if the server is hosted in a network using NAT (if the player has a home router/firewall).
	// server-port
	// type: integer (1-(2^16 - 2))
	// default: 25565
	ServerPort int `properties:"server-port"`

	// Sets the maximum distance from players that living entities may be located in order to be updated by the server, measured in chunks in each direction of the player (radius, not diameter). If entities are outside of this radius, then they will not be ticked by the server nor will they be visible to players.
	// simulation-distance
	// type: integer (3-32)
	// default: 10
	//
	// 10 is the default/recommended. If the player has major lag, this value is recommended to be reduced.
	SimulationDistance int `properties: "simulation-distance"`

	// Sets whether the server sends snoop data regularly to http://snoop.minecraft.net.
	// snooper-enabled
	// type: boolean
	// default: true
	//
	// false - disable snooping.
	// true - enable snooping.
	SnooperEnabled bool `properties:"snooper-enabled"`

	// Determines if animals can spawn.
	// spawn-animals
	// type: boolean
	// default: true
	//
	// true - Animals spawn as normal.
	// false - Animals immediately vanish.
	//
	// If the player has major lag, it is recommended to turn this off/set to false.
	SpawnAnimals bool `properties:"spawn-animals"`

	// This setting has no effect if difficulty = 0 (peaceful). If difficulty is not = 0, a monster can still spawn from a monster spawner.
	// spawn-monsters
	// type: boolean
	// default: true
	// Determines if monsters can spawn.
	//
	// true - Enabled. Monsters appear at night and in the dark.
	// false - Disabled. No monsters.
	//
	// If the player has major lag, it is recommended to turn this off/set to false.
	SpawnMonsters bool `properties:"spawn-monsters"`

	// Determines whether villagers can spawn.
	// spawn-npcs
	// type: boolean
	// default: true
	//
	// true - Enabled. Villagers spawn.
	// false - Disabled. No villagers.
	SpawnNPCs bool `properties:"spawn-npcs"`

	// Determines the side length of the square spawn protection area as 2x+1. Setting this to 0 disables the spawn protection. A value of 1 protects a 3×3 square centered on the spawn point. 2 protects 5×5, 3 protects 7×7, etc. This option is not generated on the first server start and appears when the first player joins. If there are no ops set on the server, the spawn protection is disabled automatically as well.
	// spawn-protection
	// type: integer
	// default: 16
	SpawnProtection int `properties:"spawn-protection"`

	// Enables synchronous chunk writes.
	// sync-chunk-writes
	// type: boolean
	// default: true
	SyncChunkWrites bool `properties:"sync-chunk-writes"`

	// [more information needed]
	// text-filtering-config
	// type: string
	// default: blank
	TextFilteringConfig string `properties:"text-filtering-config"`

	// Linux server performance improvements: optimized packet sending/receiving on Linux
	// use-native-transport
	// type: boolean
	// default: true
	//
	// true - Enabled. Enable Linux packet sending/receiving optimization
	// false - Disabled. Disable Linux packet sending/receiving optimization
	UseNativeTransport bool `properties:"use-native-transport"`

	// Sets the amount of world data the server sends the client, measured in chunks in each direction of the player (radius, not diameter). It determines the server-side viewing distance.
	// view-distance
	// type: integer (3-32)
	// default: 10
	//
	// 10 is the default/recommended. If the player has major lag, this value is recommended to be reduced.
	ViewDistance int `properties:"view-distance"`

	// Enables a whitelist on the server.
	// white-list
	// type: boolean
	// default: false
	//
	// With a whitelist enabled, users not on the whitelist cannot connect. Intended for private servers, such as those for real-life friends or strangers carefully selected via an application process, for example.
	//
	// false - No white list is used.
	// true - The file whitelist.json is used to generate the white list.
	//
	// Note: Ops are automatically whitelisted, and there is no need to add them to the whitelist.
	WhiteList bool `properties:"white-list"`
}

func WithMotd(m string) ServerConfig {
	return func(s *ServerProperties) *ServerProperties {
		s.Motd = m
		return s
	}
}

func WithLevelName(n string) ServerConfig {
	return func(s *ServerProperties) *ServerProperties {
		s.LevelName = n
		return s
	}
}

func WithServerPort(p int) ServerConfig {
	return func(s *ServerProperties) *ServerProperties {
		s.ServerPort = p
		return s
	}
}

// WithRcon defines RCON protocol configuration
// 'port' to be used for this protocol
// 'enabled' is to enable/disable feature
// 'pass' define the RCON password
func WithRcon(port int, enabled bool, pass string) ServerConfig {
	return func(s *ServerProperties) *ServerProperties {
		s.RconPort = port
		s.EnableRcon = enabled
		s.RconPassword = pass
		return s
	}
}

// WithQuery defines Query protocol configuration
// 'port' to be used for this protocol
// 'enabled' is to enable/disable feature
func WithQuery(port int, enabled bool) ServerConfig {
	return func(s *ServerProperties) *ServerProperties {
		s.QueryPort = port
		s.EnableQuery = enabled
		return s
	}
}

// WithSeed defines level seed
func WithSeed(seed string) ServerConfig {
	return func(s *ServerProperties) *ServerProperties {
		s.LevelSeed = seed
		return s
	}
}

// GetDefaultServerProperties returns the default server.properties representation
func GetDefaultServerProperties() (*ServerProperties, error) {
	var resp ServerProperties
	in, err := defaultConfigFiles.Open("default_values/server.properties")
	if err != nil {
		err = fmt.Errorf("reading default server.properties values: %w", err)
		return nil, err
	}
	if err := properties.NewDecoder(in).Decode(&resp); err != nil {
		err = fmt.Errorf("reading default server.properties values: %w", err)
		return nil, err
	}
	return &resp, nil
}

// GetServerProperties returns a customized server.properties representation
func GetServerProperties(cfgs ...ServerConfig) (*ServerProperties, error) {
	resp, err := GetDefaultServerProperties()
	if err != nil {
		err = fmt.Errorf("loading default values: %w", err)
		return resp, err
	}
	for _, c := range cfgs {
		resp = c(resp)
	}
	return resp, nil
}
