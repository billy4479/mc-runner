# MC Server Runner

`config.yml`:

```yml
world_dir: <path to server>
db_path: ./mc-runner.db
connect_url: localhost:25565    # Where players should connect

# Restic config for automatic backups
restic_path: /some/backup/path  # Where the backups live

# Use these defaults when using the container image
java8: java8                    # Java 8 executable
java17: java17                  # Java 17 executable
java21: java21                  # Java 21 executable
java25: java25                  # Java 25 executable

# Log level
log_level: debug
```

Inside `world_dir` place a `mc-runner.yml` file
```yml
# Java version needed for this server
java_version: 21                        

# Jar file to run (might be different if you use paper/fabric)
jar_file: server.jar                   

# Flags to pass to java _before_ `-jar`, see the following sources:
# - https://minecraft.fandom.com/wiki/Tutorials/Setting_up_a_server#Java_options
# - https://docs.papermc.io/paper/aikars-flags/
java_flags:                             
 - -Xmx6G
 - -Xms6G
 - -XX:+UnlockExperimentalVMOptions
 - -XX:+UseZGC

# Flags to pass _after_ the jar file, see https://minecraft.wiki/w/Tutorial:Setting_up_a_Java_Edition_server#Server_options
minecraft_flags:                        
 - --no-gui

# The name of this server, will be diplasplayed in the web UI
server_name: My amazing new server!     

# Ignore if you don't use carpet, otherwise set the same value you used in https://carpet.tis.world/docs/rules#fakeplayernameprefix
bot_prefix: "[BOT]"                     
```
