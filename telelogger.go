// Package telelogger provides advanced logging capabilities through Telegram.
//
// This package allows users to implement sophisticated logging
// with various output formats and destinations through Telegram Bot API.
package telelogger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Version represents the current version of the package
const Version = "0.1.0"

// ParseMode represents the available formatting modes for Telegram messages.
// Can be one of: "HTML", "Markdown", or "MarkdownV2".
type ParseMode string

const (
	// ParseModeHTML enables HTML-style formatting
	ParseModeHTML ParseMode = "HTML"
	// ParseModeMarkdown enables basic Markdown formatting
	ParseModeMarkdown ParseMode = "Markdown"
	// ParseModeMarkdownV2 enables enhanced Markdown formatting with more features
	ParseModeMarkdownV2 ParseMode = "MarkdownV2"
)

// FormatterFunc is a function type for message formatting.
// It takes a message string and returns a formatted string.
type FormatterFunc func(message string) string

// Default formatters with emojis and predefined formats
func baseInfoFormat(msg string) string    { return fmt.Sprintf("ℹ️ Info:\n%s", msg) }
func baseErrorFormat(msg string) string   { return fmt.Sprintf("❌ Error:\n%s", msg) }
func baseSuccessFormat(msg string) string { return fmt.Sprintf("✅ Success:\n%s", msg) }
func baseWarnFormat(msg string) string    { return fmt.Sprintf("🚨 Warning:\n%s", msg) }

// Config holds the configuration for the Telelogger instance.
type Config struct {
	// BotToken is the Telegram Bot Token obtained from BotFather
	BotToken string

	// ChatID is the Telegram Chat ID where messages will be sent
	ChatID int64

	// ParseMode specifies the formatting mode for messages
	// Can be HTML, Markdown, or MarkdownV2
	// If not provided, no formatting will be applied
	ParseMode ParseMode

	// InfoFormatter is a custom formatter for info messages
	// If not provided, uses default format with ℹ️ emoji
	InfoFormatter FormatterFunc

	// ErrorFormatter is a custom formatter for error messages
	// If not provided, uses default format with ❌ emoji
	ErrorFormatter FormatterFunc

	// SuccessFormatter is a custom formatter for success messages
	// If not provided, uses default format with ✅ emoji
	SuccessFormatter FormatterFunc

	// WarnFormatter is a custom formatter for warning messages
	// If not provided, uses default format with 🚨 emoji
	WarnFormatter FormatterFunc
}

// Telelogger is the main struct for sending formatted log messages to Telegram.
// It provides methods for sending different types of messages (info, error, success, warning)
// with optional message formatting and custom formatters.
type Telelogger struct {
	chatID           int64
	baseURL          string
	parseMode        ParseMode
	infoFormatter    FormatterFunc
	errorFormatter   FormatterFunc
	successFormatter FormatterFunc
	warnFormatter    FormatterFunc
	client           *http.Client
}

// message represents the structure of a Telegram message for API requests
type message struct {
	ChatID    int64     `json:"chat_id"`
	Text      string    `json:"text"`
	ParseMode ParseMode `json:"parse_mode,omitempty"`
}

// New creates a new Telelogger instance with the provided configuration.
//
// Example:
//
//	logger := telelogger.New(telelogger.Config{
//	    BotToken: "your-bot-token",
//	    ChatID:   123456789,
//	    ParseMode: telelogger.ParseModeHTML,
//	})
func New(config Config) *Telelogger {
	t := &Telelogger{
		chatID:           config.ChatID,
		baseURL:          fmt.Sprintf("https://api.telegram.org/bot%s", config.BotToken),
		parseMode:        config.ParseMode,
		infoFormatter:    config.InfoFormatter,
		errorFormatter:   config.ErrorFormatter,
		successFormatter: config.SuccessFormatter,
		warnFormatter:    config.WarnFormatter,
		client:           &http.Client{},
	}

	// Set default formatters if not provided
	if t.infoFormatter == nil {
		t.infoFormatter = baseInfoFormat
	}
	if t.errorFormatter == nil {
		t.errorFormatter = baseErrorFormat
	}
	if t.successFormatter == nil {
		t.successFormatter = baseSuccessFormat
	}
	if t.warnFormatter == nil {
		t.warnFormatter = baseWarnFormat
	}

	return t
}

// Log sends a generic message to Telegram.
//
// Example:
//
//	err := logger.Log("Generic message")
func (t *Telelogger) Log(msg string) error {
	return t.sendMessage(msg, t.parseMode)
}

// LogWithParseMode sends a generic message to Telegram with a specific parse mode.
//
// Example:
//
//	err := logger.LogWithParseMode("Message with <b>bold</b> text", telelogger.ParseModeHTML)
func (t *Telelogger) LogWithParseMode(msg string, parseMode ParseMode) error {
	return t.sendMessage(msg, parseMode)
}

// LogError sends an error message to Telegram.
// The error parameter can be either an error object or a string.
//
// Example:
//
//	err := logger.LogError("Database connection failed")
//	// or
//	err := logger.LogError(fmt.Errorf("Database connection failed"))
func (t *Telelogger) LogError(err interface{}) error {
	var msg string
	switch v := err.(type) {
	case error:
		msg = v.Error()
	case string:
		msg = v
	default:
		msg = fmt.Sprintf("%v", v)
	}
	return t.sendMessage(t.errorFormatter(msg), t.parseMode)
}

// LogInfo sends an info message to Telegram.
//
// Example:
//
//	err := logger.LogInfo("Application started successfully")
func (t *Telelogger) LogInfo(msg string) error {
	return t.sendMessage(t.infoFormatter(msg), t.parseMode)
}

// LogSuccess sends a success message to Telegram.
//
// Example:
//
//	err := logger.LogSuccess("Backup completed successfully")
func (t *Telelogger) LogSuccess(msg string) error {
	return t.sendMessage(t.successFormatter(msg), t.parseMode)
}

// LogWarn sends a warning message to Telegram.
//
// Example:
//
//	err := logger.LogWarn("Low disk space")
func (t *Telelogger) LogWarn(msg string) error {
	return t.sendMessage(t.warnFormatter(msg), t.parseMode)
}

// sendMessage handles the actual sending of messages to Telegram.
// It formats the message according to the specified parse mode and sends it via the Telegram Bot API.
func (t *Telelogger) sendMessage(text string, parseMode ParseMode) error {
	msg := message{
		ChatID:    t.chatID,
		Text:      text,
		ParseMode: parseMode,
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	resp, err := t.client.Post(
		fmt.Sprintf("%s/sendMessage", t.baseURL),
		"application/json",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned non-200 status code: %d", resp.StatusCode)
	}

	return nil
}
