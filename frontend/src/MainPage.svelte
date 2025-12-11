<script lang="ts">
  import { getMeOrError, setMeOrError } from "$lib/state.svelte";
  import { logout } from "$lib/api";
  import { ServerSocket } from "$lib/ws.svelte";

  import { Button } from "flowbite-svelte";

  let consoleOutput: HTMLDivElement;

  let serverSocket = new ServerSocket();

  $effect(() => {
    consoleOutput.scrollTo(0, consoleOutput.scrollHeight);
  });

  $effect(() => {
    serverSocket.connect();

    return () => {
      serverSocket.close();
    };
  });

  const dummyConsoleOutput = `[20:38:18] [main/INFO]: Loading Minecraft 1.21.1 with Fabric Loader 0.16.2
[20:38:18] [ForkJoinPool-1-worker-1/WARN]: Mod novaultlimit uses the version r1.1.0-1.21.1 which isn't compatible with Loader's extended semantic version format (Could not parse version number component 'r1'!), SemVer is recommended for reliably evaluating dependencies and prioritizing newer version
[20:38:18] [main/WARN]: Warnings were found!
 - Mod 'Forge Config API Port' (forgeconfigapiport) 21.1.0 recommends any version of modmenu, which is missing!
	 - You should install any version of modmenu for the optimal experience.
[20:38:18] [main/INFO]: Loading 69 mods:
	- accelerateddecay 21.0.0
	- architectury 13.0.8
	- badpackets 0.8.1
	- balm-fabric 21.0.23
	   \-- kuma_api 21.0.5-SNAPSHOT
	- carpet 1.4.147+v240613
	- carpet-extra 1.4.148
	- carpet-tis-addition 1.63.1
	   \-- conditional-mixin 0.6.2
	- cloth-config 15.0.130
	   \-- cloth-basic-math 0.6.1
	- clumps 19.0.0.1
	- fabric-api 0.102.1+1.21.1
	   |-- fabric-api-base 0.4.42+6573ed8c6a
	   |-- fabric-api-lookup-api-v1 1.6.68+b55973446a
	   |-- fabric-biome-api-v1 13.0.29+5bd9f1bc6a
	   |-- fabric-block-api-v1 1.0.22+0af3f5a76a
	   |-- fabric-block-view-api-v2 1.0.10+6573ed8c6a
	   |-- fabric-blockrenderlayer-v1 1.1.52+0af3f5a76a
	   |-- fabric-client-tags-api-v1 1.1.15+6573ed8c6a
	   |-- fabric-command-api-v1 1.2.49+f71b366f6a
	   |-- fabric-command-api-v2 2.2.28+6ced4dd96a
	   |-- fabric-commands-v0 0.2.66+df3654b36a
	   |-- fabric-content-registries-v0 8.0.16+b55973446a
	   |-- fabric-convention-tags-v1 2.0.20+7f945d5b6a
	   |-- fabric-convention-tags-v2 2.6.0+605f22ad6a
	   |-- fabric-crash-report-info-v1 0.2.29+0af3f5a76a
	   |-- fabric-data-attachment-api-v1 1.1.27+1daea2156a
	   |-- fabric-data-generation-api-v1 20.2.17+16c4ae256a
	   |-- fabric-dimensions-v1 4.0.0+6fc22b996a
	   |-- fabric-entity-events-v1 1.6.12+6fc22b996a
	   |-- fabric-events-interaction-v0 0.7.12+ba9dae066a
	   |-- fabric-game-rule-api-v1 1.0.53+6ced4dd96a
	   |-- fabric-item-api-v1 11.0.0+afdfc9216a
	   |-- fabric-item-group-api-v1 4.1.4+780172706a
	   |-- fabric-key-binding-api-v1 1.0.47+0af3f5a76a
	   |-- fabric-keybindings-v0 0.2.45+df3654b36a
	   |-- fabric-lifecycle-events-v1 2.3.12+6c1df3606a
	   |-- fabric-loot-api-v2 3.0.14+3f89f5a56a
	   |-- fabric-loot-api-v3 1.0.2+3f89f5a56a
	   |-- fabric-message-api-v1 6.0.13+6573ed8c6a
	   |-- fabric-model-loading-api-v1 2.0.0+fe474d6b6a
	   |-- fabric-networking-api-v1 4.2.2+60c3209b6a
	   |-- fabric-object-builder-api-v1 15.2.0+40875a936a
	   |-- fabric-particles-v1 4.0.2+6573ed8c6a
	   |-- fabric-recipe-api-v1 5.0.12+650897126a
	   |-- fabric-registry-sync-v0 5.1.2+60c3209b6a
	   |-- fabric-renderer-api-v1 3.4.0+c705a49c6a
	   |-- fabric-renderer-indigo 1.7.0+c705a49c6a
	   |-- fabric-renderer-registries-v1 3.2.68+df3654b36a
	   |-- fabric-rendering-data-attachment-v1 0.3.48+73761d2e6a
	   |-- fabric-rendering-fluids-v1 3.1.6+1daea2156a
	   |-- fabric-rendering-v0 1.1.71+df3654b36a
	   |-- fabric-rendering-v1 5.0.5+df16efd06a
	   |-- fabric-resource-conditions-api-v1 4.3.0+8dc279b16a
	   |-- fabric-resource-loader-v0 1.3.0+565991296a
	   |-- fabric-screen-api-v1 2.0.25+8b68f1c76a
	   |-- fabric-screen-handler-api-v1 1.3.86+b55973446a
	   |-- fabric-sound-api-v1 1.0.23+6573ed8c6a
	   |-- fabric-transfer-api-v1 5.1.17+1db1cc1f6a
	   \-- fabric-transitive-access-wideners-v1 6.1.0+1daea2156a
	- fabricloader 0.16.2
	   \-- mixinextras 0.4.1
	- ferritecore 7.0.0
	- forgeconfigapiport 21.1.0
	   |-- com_electronwill_night-config_core 3.8.0
	   \-- com_electronwill_night-config_toml 3.8.0
	- inventorysorter 1.9.0-1.21
	   \-- kyrptconfig 1.6.0-1.20.4
	- java 21
	- krypton 0.2.8
	   \-- com_velocitypowered_velocity-native 3.3.0-SNAPSHOT
	- minecraft 1.21.1
	- mr_no_endermangrief 1-v.2.1.2
	- nochatreports 1.21-v2.8.0
	- notenoughcrashes 4.4.8+1.21
	- novaultlimit r1.1.0-1.21.1
	- servux 0.3.7
	   \-- fabric-permissions-api-v0 0.3.1
	- shulkerboxtooltip 5.1.0+1.21.1
	- wthit 12.4.0
[20:38:19] [main/INFO]: SpongePowered MIXIN Subsystem Version=0.8.7 Source=file://net/fabricmc/sponge-mixin/0.15.2+mixin.0.8.7/sponge-mixin-0.15.2+mixin.0.8.7.jar Service=Knot/Fabric Env=SERVER
[20:38:22] [main/INFO]: Compatibility level set to JAVA_17
[20:38:22] [main/INFO]: Compatibility level set to JAVA_21
[20:38:22] [main/WARN]: Reference map 'accelerated-decay-common-refmap.json' for accelerateddecay-common.mixins.json could not be read. If this is a development environment you can ignore this message
[20:38:22] [main/WARN]: Reference map 'forgeconfigapiport.common.refmap.json' for forgeconfigapiport.common.mixins.json could not be read. If this is a development environment you can ignore this message
[20:38:23] [main/WARN]: Error loading class: net/coderbot/iris/Iris (java.lang.ClassNotFoundException: net/coderbot/iris/Iris)
[20:38:23] [main/WARN]: @Mixin target net.coderbot.iris.Iris was not found notenoughcrashes.fabric.mixins.json:iris.SilentNEC from mod notenoughcrashes
[20:38:23] [main/INFO]: Initializing MixinExtras via com.llamalad7.mixinextras.service.MixinExtrasServiceImpl(version=0.4.1).
[20:38:27] [main/WARN]: Method overwrite conflict for getClimateSettings in architectury.mixins.json:BiomeAccessor from mod architectury, previously written by carpet.mixins.Biome_scarpetMixin. Skipping method.
[20:38:34] [TISCM Mapping/INFO]: Yarn mapping file yarn-1.21.1+build.3-v2.tiny loaded
[20:38:34] [main/WARN]: @Redirect conflict. Skipping carpet-extra.mixins.json:ServerPlayNetworkHandlerMixin from mod carpet-extra->@Redirect::carpetextra_removeHitPosCheck(Lnet/minecraft/class_243;Lnet/minecraft/class_243;)Lnet/minecraft/class_243; with priority 1000, already redirected by mixins.servux.json:MixinServerPlayNetworkHandler from mod servux->@Redirect::servux$removeHitPosCheck(Lnet/minecraft/class_243;Lnet/minecraft/class_243;)Lnet/minecraft/class_243; with priority 1005
[20:38:34] [main/INFO]: Compression will use Java, encryption will use Java
[20:38:34] [main/INFO]: Initialized NoVaultLimit although there really isn't anything to initialize.
[20:38:35] [main/INFO]: [WTHIT] Initializing common plugin waila:core at mcp.mobius.waila.plugin.core.CoreCommonPlugin
[20:38:35] [main/INFO]: [WTHIT] Initializing common plugin waila:vanilla at mcp.mobius.waila.plugin.vanilla.VanillaCommonPlugin
[20:38:35] [main/INFO]: [WTHIT] Initializing common plugin waila:harvest at mcp.mobius.waila.plugin.harvest.HarvestCommonPlugin
[20:38:35] [main/INFO]: [WTHIT] Initializing common plugin waila:fabric at mcp.mobius.waila.plugin.fabric.FabricCommonPlugin
[20:38:35] [main/INFO]: [WTHIT] Initializing common plugin waila:extra at mcp.mobius.waila.plugin.extra.ExtraPlugin
[20:38:35] [main/INFO]: [WTHIT] Plugin config reloaded
[20:38:35] [main/INFO]: Krypton is now accelerating your Minecraft server's networking stack 🚀
[20:38:36] [main/INFO]: Environment: Environment[sessionHost=https://sessionserver.mojang.com, servicesHost=https://api.minecraftservices.com, name=PROD]
[20:38:39] [main/INFO]: Loaded 1290 recipes
[20:38:39] [main/INFO]: Loaded 1402 advancements
[20:38:39] [main/INFO]: Applied 0 biome modifications to 0 of 64 new biomes in 1.362 ms
[20:38:40] [Server thread/INFO]: [WTHIT] Plugin config reloaded
[20:38:40] [Server thread/INFO]: Starting minecraft server version 1.21.1
[20:38:40] [Server thread/INFO]: Loading properties
[20:38:40] [Server thread/INFO]: Default game type: SURVIVAL
[20:38:40] [Server thread/INFO]: Generating keypair
[20:38:40] [Server thread/INFO]: Starting Minecraft server on *:25565
[20:38:41] [Server thread/INFO]: Using default channel type
[20:38:41] [Server thread/INFO]: Preparing level "Friends"
[20:38:41] [Server thread/INFO]: [CM] Loaded 10 settings from carpet.conf
[20:38:42] [Server thread/INFO]: Preparing start region for dimension minecraft:overworld
[20:38:42] [Worker-Main-5/INFO]: Preparing spawn area: 0%
[20:38:42] [Worker-Main-8/INFO]: Preparing spawn area: 18%
[20:38:43] [Server thread/INFO]: Time elapsed: 790 ms
[20:38:43] [Server thread/INFO]: Done (1.958s)! For help, type "help"
[20:38:51] [User Authenticator #1/INFO]: UUID of player is 
[20:38:51] [Server thread/INFO]: [/127.0.0.1:58892] logged in with entity id 337 at (193.42515339880384, 71.5, 104.04090481213456)
[20:38:51] [Server thread/INFO]: joined the game
[20:38:52] [Server thread/INFO]: [ShulkerBoxTooltip] QwertyRacing: protocol version: 2.0
[20:38:52] [Server thread/INFO]: [ShulkerBoxTooltip Plugins] Loading 1 plugin: shulkerboxtooltip
[20:38:52] [Server thread/INFO]: [ShulkerBoxTooltip] overriding preview provider shulkerboxtooltip:trapped_chest with shulkerboxtooltip:chest for item minecraft:chest
[20:38:52] [Server thread/INFO]: [ShulkerBoxTooltip Plugins] Registered 31 providers for mod shulkerboxtooltip
[20:38:54] [Server thread/ERROR]: Error loading saved data: map_76
java.lang.ArrayIndexOutOfBoundsException: arraycopy: length -1 is negative
	at java.base/java.lang.System.arraycopy(Native Method) ~[?:?]
	at java.base/java.io.PushbackInputStream.unread(Unknown Source) ~[?:?]
	at net.minecraft.class_26.method_17921(class_26.java:119) ~[server-intermediary.jar:?]
	at net.minecraft.class_26.method_17923(class_26.java:95) ~[server-intermediary.jar:?]
	at net.minecraft.class_26.method_120(class_26.java:75) ~[server-intermediary.jar:?]
	at net.minecraft.class_26.method_20786(class_26.java:63) ~[server-intermediary.jar:?]
	at net.minecraft.class_3218.method_17891(class_3218.java:1248) ~[server-intermediary.jar:?]
	at net.minecraft.class_1806.method_7997(class_1806.java:56) ~[server-intermediary.jar:?]
	at net.minecraft.class_3231.method_18756(class_3231.java:128) ~[server-intermediary.jar:?]
	at net.minecraft.class_3898.method_18727(class_3898.java:1157) ~[server-intermediary.jar:?]
	at net.minecraft.class_3215.method_12127(class_3215.java:311) ~[server-intermediary.jar:?]
	at net.minecraft.class_3218.method_18765(class_3218.java:350) ~[server-intermediary.jar:?]
	at net.minecraft.server.MinecraftServer.mixinextras$bridge$method_18765$315(MinecraftServer.java) ~[server-intermediary.jar:?]
	at net.minecraft.server.MinecraftServer.wrapOperation$cej000$carpet-tis-addition$yeetUpdateSuppressionCrash_implOnTickWorlds(MinecraftServer.java:7159) ~[server-intermediary.jar:?]
	at net.minecraft.server.MinecraftServer.method_3813(MinecraftServer.java:1021) ~[server-intermediary.jar:?]
	at net.minecraft.class_3176.method_3813(class_3176.java:299) ~[server-intermediary.jar:?]
	at net.minecraft.server.MinecraftServer.method_3748(MinecraftServer.java:912) ~[server-intermediary.jar:?]
	at net.minecraft.server.MinecraftServer.method_29741(MinecraftServer.java:697) ~[server-intermediary.jar:?]
	at net.minecraft.server.MinecraftServer.method_29739(MinecraftServer.java:281) ~[server-intermediary.jar:?]
	at java.base/java.lang.Thread.run(Unknown Source) [?:?]
[20:39:40] [Server thread/INFO]: Villager class_1646['Villager'/497, l='ServerLevel[Friends]', x=291.70, y=54.00, z=88.70] died, message: 'Villager was squished too much'
[20:46:57] [Server thread/INFO]: Villager class_1646['Villager'/494, l='ServerLevel[Friends]', x=291.70, y=54.00, z=88.30] died, message: 'Villager was squished too much'
[20:59:07] [Server thread/INFO]: Villager class_1646['Villager'/493, l='ServerLevel[Friends]', x=291.70, y=54.00, z=88.70] died, message: 'Villager was squished too much'
[21:07:16] [Server thread/INFO]: lost connection: Disconnected
[21:07:16] [Server thread/INFO]: left the game
[21:07:20] [Server thread/INFO]: Stopping server
[21:07:20] [Server thread/INFO]: Saving players
[21:07:20] [Server thread/INFO]: Saving worlds
[21:07:20] [Server thread/INFO]: Saving chunks for level 'ServerLevel[Friends]'/minecraft:overworld
[21:07:20] [Server thread/INFO]: Saving chunks for level 'ServerLevel[Friends]'/minecraft:the_nether
[21:07:20] [Server thread/INFO]: Saving chunks for level 'ServerLevel[Friends]'/minecraft:the_end
[21:07:20] [Server thread/INFO]: ThreadedAnvilChunkStorage (Friends): All chunks are saved
[21:07:20] [Server thread/INFO]: ThreadedAnvilChunkStorage (DIM-1): All chunks are saved
[21:07:20] [Server thread/INFO]: ThreadedAnvilChunkStorage (DIM1): All chunks are saved
[21:07:20] [Server thread/INFO]: ThreadedAnvilChunkStorage: All dimensions are saved
`;
</script>

<div class="w-4/5 rounded border border-gray-300">
  <div
    class="flex flex-row items-center justify-between border-b border-gray-300 px-5 py-3"
  >
    <div>
      Welcome <b> {getMeOrError().Name.String} </b>

      <span class="text-gray-600">
        - {#if serverSocket.isConnected}
          <b>Connected </b>
          to
          <span class="inline font-mono">
            mc-runner@{serverSocket.serverState?.version}
          </span>
          for server "<b> {serverSocket.serverState?.server_name} </b>"
        {:else if serverSocket.isConnecting}
          <b> Connecting... </b>
        {:else}
          <b> Disconnected </b>
        {/if}
      </span>
    </div>

    <div>
      {#if getMeOrError().ID == 0}
        <Button size="xs" outline color="alternative">Invite</Button>
      {/if}
      <Button
        size="xs"
        outline
        color="alternative"
        onclick={async () => {
          await logout();
          await setMeOrError();
        }}
      >
        Logout
      </Button>
    </div>
  </div>
  <div class="flex flex-row">
    <div
      bind:this={consoleOutput}
      class="mx-5 my-3 h-64 w-full overflow-scroll rounded border-gray-300 bg-gray-100 px-3 py-2 font-mono text-sm whitespace-pre"
    >
      {#if !serverSocket.serverState?.is_running}
        <div class="relative h-full">
          <div class="blur-sm select-none">
            This is dummy text
            <br />
            if you're reading this
            <br />
            you should reconsider
            <br />
            your life choices.
            <br />
            join the club!
            <br />
            Since you're here though,
            <br />
            tell me something about you're life
            <br />
            how's it going lately?
          </div>
          <div
            class="absolute top-0 left-0 flex h-full w-full items-center justify-center font-sans text-base"
          >
            The server is not running.
          </div>
        </div>
      {:else}
        {dummyConsoleOutput}
      {/if}
    </div>
    <div class="border-l border-gray-300"></div>
    <div class="flex w-2/5 flex-col items-center py-3 whitespace-nowrap">
      <Button
        color="green"
        disabled={!serverSocket.serverState ||
          serverSocket.serverState.is_running}
        class="mx-5 mb-3"
        onclick={() => {
          serverSocket.startServer();
        }}
      >
        Start server
      </Button>
      <div class="w-full border-t border-gray-300 px-5 py-3">
        {#if serverSocket.serverState?.is_running}
          {serverSocket.serverState?.online_players.length} players online:
          <ul class="list-inside list-disc break-words whitespace-break-spaces">
            {#each serverSocket.serverState?.online_players as player}
              <li>{player}</li>
            {/each}
          </ul>
        {:else}
          The server is not running.
        {/if}
      </div>
    </div>
  </div>
</div>
