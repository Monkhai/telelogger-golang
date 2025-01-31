# telelogger

A lightweight, easy-to-use Telegram logging utility for Go applications.

## Features

- Simple integration with Telegram Bot API
- Multiple log levels (info, error, success, warn)
- Customizable message formatters
- Strong Go type system
- Zero external dependencies

## Installation

```bash
go get github.com/monkhai/telelogger-golang
```

## Quick Start

```go
package main

import "github.com/monkhai/telelogger-golang"

func main() {
    logger := telelogger.New(telelogger.Config{
        BotToken: "YOUR_BOT_TOKEN",
        ChatID:   YOUR_CHAT_ID,
    })

    // Basic logging
    logger.LogInfo("Hello, world!")
    logger.LogError("Something went wrong!")
    logger.LogSuccess("Operation completed successfully!")
    logger.LogWarn("Warning: Resource running low")
}
```

## Configuration

The `New` function accepts a `Config` struct with the following options:

```go
type Config struct {
    // Your Telegram Bot Token
    BotToken string

    // Target Chat ID where messages will be sent
    ChatID int64

    // The formatting of the message
    // Can be ParseModeHTML, ParseModeMarkdown, or ParseModeMarkdownV2
    ParseMode ParseMode

    // Custom formatter for info messages
    InfoFormatter FormatterFunc

    // Custom formatter for error messages
    ErrorFormatter FormatterFunc

    // Custom formatter for success messages
    SuccessFormatter FormatterFunc

    // Custom formatter for warning messages
    WarnFormatter FormatterFunc
}

// FormatterFunc is a function type for message formatting
type FormatterFunc func(message string) string
```

### Custom Formatters Example

You can customize how messages are formatted before they're sent to Telegram.
Notice the `<b>` tags in the formatters, adding bold text to the titles.
This allows you to add more information to the messages, such as links, bold text, etc.
For more information on the different parse modes, see the [Telegram API documentation](https://core.telegram.org/bots/api#formatting-options).

```go
logger := telelogger.New(telelogger.Config{
    BotToken:  "YOUR_BOT_TOKEN",
    ChatID:    YOUR_CHAT_ID,
    ParseMode: telelogger.ParseModeHTML, // allows us to add <b> tags and more
    InfoFormatter: func(msg string) string {
        return fmt.Sprintf("ℹ️ <b>INFO:</b>\n%s", msg)
    },
    ErrorFormatter: func(msg string) string {
        return fmt.Sprintf("❌ <b>ERROR:</b>\n%s", msg)
    },
    SuccessFormatter: func(msg string) string {
        return fmt.Sprintf("✅ <b>SUCCESS:</b>\n%s", msg)
    },
    WarnFormatter: func(msg string) string {
        return fmt.Sprintf("🚨️ <b>WARNING:</b>\n%s", msg)
    },
})
```

### Error Handling

Unlike the TypeScript version, this package follows Go's error handling patterns:

```go
// All logging methods return an error that you can handle
if err := logger.LogInfo("Hello, world!"); err != nil {
    // Handle error
}

// You can pass either a string or an error to LogError
err := someFunction()
if err != nil {
    logger.LogError(err) // Accepts error interface
}
logger.LogError("Something went wrong") // Also accepts string
```

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
