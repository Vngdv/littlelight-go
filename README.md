# LittleLight

This is Discordbot is used to create Channels for Users

This Bot is written in a way that allows it to work without any config files.

## Configuration

Runtime Flags
```
-t <Discord Token> -- Sets the Discord Token
-j <Channel Name> -- Sets the default Join Channel name. Defaul: "📢 Join to own"
-c <Category Identifier> -- Sets the Default Emoji, Sentence etc to Check for. Default: "🎤"
-n <Channel Names> -- Sets the List of Channel Names to choose from. Seperated by ";". Default "Voice Channel; 🎈 Party Room"
```
Example usage: `-t <TOKEN> -j "Join here" -c "🚩" -n "Voice 1;Voice2"`

**Some options are also available from the environment.**

Environment variables
```
TOKEN=<Discord Token>-- Sets the Discord Token
```