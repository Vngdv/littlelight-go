# littlelight-go

This is a discord bot that creates voice channels on demand.

LittleLight is a bot that was made with three core features in mind:
- simplicity
- no configuration files or storage
- fast channel creation

Unlike other Bots, LittleLight uses your private chat with the bot as its store for channel names. We also just rename the empty channel you join instead of moving you away.

## Table of Contents

- [Install](#install)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contributing](#contributing)
- [License](#license)

## Install

Currently, this package is only runnable inside a docker container.
If you want to run this as a standalone binary, you have to compile it yourself.

Run with docker:
```
docker run ghcr.io/vngdv/littlelight:latest -b 96 -t <Your token>
```

### Runtime flags
```
-t <Discord Token> -- Sets the Discord Token
-j <Channel Name> -- Sets the default Join Channel name. Default: "📢 Join to own"
-c <Category Identifier> -- Sets the Default Emoji, Sentence etc to Check for. Default: "🎤"
-n <Channel Names> -- Sets the List of Channel Names to choose from. Seperated by ";". Default "Voice Channel; 🎈 Party Room"
-b <Default Bitrate> -- Sets the default bitrate for the Channels in kbps. Default: 64
--allow-channelnames=<True/false> -- Allows the use of custom channel names via dm chat. Default: true
```

## Usage

Create a category inside discord that contains your category identifier. The default is: 🎤
Now create a single voice channel inside and the bot will take over the rest.

If you want to modify your channel name, just write the bot a DM with your preferred channel name inside. The next channel you create will get that name.

## Maintainers

[@Vngdv](https://github.com/Vngdv)

## Contributing

PRs accepted.

Small note: If editing the README, please conform to the [standard-readme](https://github.com/RichardLitt/standard-readme) specification.

## License

MIT © 2022 Vngdv
