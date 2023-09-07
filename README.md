# MfuseBot
This is a simple Discord bot to test your MOHAA Morphuse scripts on.

It depends on: [Morfuse](https://github.com/morfuse/morfuse)

![Demo GIF](https://i.ibb.co/MN1PZ9w/demo.gif)
## 🛠️ Installation Steps:

This project is [hosted on github](https://github.com/DraGo0oN/mfuse-bot). You can clone this project directly using this command:

```
git clone https://github.com/DraGo0oN/mfuse-bot
```

Use go mod tidy to install/update all the dependencies:
```
go mod tidy
```

Build from source:
```
go build
```

## Usage
Start:
```
./mfusebot
```

## 🛠️ Slash Command

```
/mfuse pops up a modal to input your Morpheus code on.
```

## 🛠️ Configuration

Configuration is done via `.env`:

```json
# discord bot token
TOKEN=<your discord bot token>
# bot client id
CLIENT_ID=<your discord bot client it>
# guild id (server id)
GUILD_ID=<your server id>
```

# Contributing
Public contributions are welcome!  
You can create a [new issue](https://github.com/DraGo0oN/mfuse-bot/issues/new) for bugs, or feel free to open a [pull request](https://github.com/DraGo0oN/mfuse-bot/pulls) for any and all your changes or work-in-progress features.

## License

This project is licensed under the MIT License. See the LICENSE file for details.