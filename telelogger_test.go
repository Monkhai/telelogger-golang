package telelogger_test

import (
	"errors"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
	"github.com/monkhai/telelogger-golang"
)

var (
	testLogger *telelogger.Telelogger
)

func TestMain(m *testing.M) {
	// Load the .env file before running tests
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found or error loading it: %v", err)
	}

	// Initialize the test logger
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatIDStr := os.Getenv("TELEGRAM_CHAT_ID")

	if botToken != "" && chatIDStr != "" {
		chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
		if err == nil {
			testLogger = telelogger.New(telelogger.Config{
				BotToken: botToken,
				ChatID:   chatID,
			})
		}
	}

	// Run tests
	os.Exit(m.Run())
}

func skipIfNoLogger(t *testing.T) {
	if testLogger == nil {
		t.Skip("Skipping test: TELEGRAM_BOT_TOKEN and TELEGRAM_CHAT_ID environment variables are required")
	}
}

func TestNew(t *testing.T) {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatIDStr := os.Getenv("TELEGRAM_CHAT_ID")

	if botToken == "" || chatIDStr == "" {
		t.Skip("Skipping test: TELEGRAM_BOT_TOKEN and TELEGRAM_CHAT_ID environment variables are required")
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		t.Fatalf("Failed to parse TELEGRAM_CHAT_ID: %v", err)
	}

	logger := telelogger.New(telelogger.Config{
		BotToken: botToken,
		ChatID:   chatID,
	})

	if logger == nil {
		t.Error("New() should return a non-nil logger")
	}
}

func TestVersion(t *testing.T) {
	if telelogger.Version == "" {
		t.Error("Version should not be empty")
	}
}

func TestCustomFormatters(t *testing.T) {
	customInfo := func(msg string) string { return "Custom:" + msg }

	logger := telelogger.New(telelogger.Config{
		BotToken:      "test-token",
		ChatID:        123456789,
		InfoFormatter: customInfo,
	})

	if logger == nil {
		t.Error("New() with custom formatter should return a non-nil logger")
	}
}

func TestLogInfo(t *testing.T) {
	skipIfNoLogger(t)
	err := testLogger.LogInfo("Test info message")
	if err != nil {
		t.Errorf("LogInfo failed: %v", err)
	}
}

func TestLogError(t *testing.T) {
	skipIfNoLogger(t)

	// Test with string
	err := testLogger.LogError("Test error message")
	if err != nil {
		t.Errorf("LogError with string failed: %v", err)
	}

	// Test with error
	testErr := errors.New("test error")
	err = testLogger.LogError(testErr)
	if err != nil {
		t.Errorf("LogError with error failed: %v", err)
	}
}

func TestLogSuccess(t *testing.T) {
	skipIfNoLogger(t)
	err := testLogger.LogSuccess("Test success message")
	if err != nil {
		t.Errorf("LogSuccess failed: %v", err)
	}
}

func TestLogWarn(t *testing.T) {
	skipIfNoLogger(t)
	err := testLogger.LogWarn("Test warning message")
	if err != nil {
		t.Errorf("LogWarn failed: %v", err)
	}
}

func TestLogWithParseMode(t *testing.T) {
	skipIfNoLogger(t)
	err := testLogger.LogWithParseMode("Test <b>bold</b> message", telelogger.ParseModeHTML)
	if err != nil {
		t.Errorf("LogWithParseMode failed: %v", err)
	}
}

func TestLog(t *testing.T) {
	skipIfNoLogger(t)
	err := testLogger.Log("Test generic message")
	if err != nil {
		t.Errorf("Log failed: %v", err)
	}
}
