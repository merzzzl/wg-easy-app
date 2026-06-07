package notification

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/go-telegram/bot"
	_ "modernc.org/sqlite"

	"wg-easy-app/backend/internal/config"
	"wg-easy-app/backend/internal/repository/postgres"
	telegramrepo "wg-easy-app/backend/internal/repository/telegram"
)

func TestNotifyApprovedUsersSendsOnlyToApprovedUsers(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := openNotificationTestDB(t)
	insertNotificationTestUser(t, db, 1001, "approved_one", 7001, "approved")
	insertNotificationTestUser(t, db, 1002, "waiting_one", 7002, "waiting_approve")
	insertNotificationTestUser(t, db, 1003, "approved_two", 7003, "approved")

	client := &telegramTestClient{t: t}

	botClient, err := bot.New("test", bot.WithHTTPClient(time.Second, client), bot.WithSkipGetMe())
	if err != nil {
		t.Fatalf("new bot: %v", err)
	}

	service := New(&config.Config{AdminUsername: "admin"}, postgres.NewRepository(db), telegramrepo.New(botClient))

	sent, err := service.NotifyApprovedUsers(ctx, "maintenance tonight")
	if err != nil {
		t.Fatalf("notify approved users: %v", err)
	}

	if sent != 2 {
		t.Fatalf("sent = %d, want 2", sent)
	}

	if len(client.requests) != 2 {
		t.Fatalf("telegram requests = %d, want 2", len(client.requests))
	}

	gotChatIDs := map[string]bool{}
	for _, request := range client.requests {
		if request.Get("text") != "maintenance tonight" {
			t.Fatalf("text = %q, want %q", request.Get("text"), "maintenance tonight")
		}

		gotChatIDs[request.Get("chat_id")] = true
	}

	if !gotChatIDs["7001"] || !gotChatIDs["7003"] || gotChatIDs["7002"] {
		encoded, _ := json.Marshal(gotChatIDs)
		t.Fatalf("chat ids = %s, want only 7001 and 7003", encoded)
	}
}

type telegramTestClient struct {
	t        *testing.T
	requests []url.Values
}

func (c *telegramTestClient) Do(r *http.Request) (*http.Response, error) {
	c.t.Helper()

	if r.URL.Path != "/bottest/sendMessage" {
		c.t.Fatalf("unexpected telegram path: %s", r.URL.Path)
	}

	if err := r.ParseMultipartForm(1 << 20); err != nil {
		c.t.Fatalf("parse telegram request: %v", err)
	}

	c.requests = append(c.requests, r.Form)

	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(`{"ok":true,"result":{"message_id":1,"date":1,"chat":{"id":1,"type":"private"},"text":"ok"}}`)),
	}, nil
}

func openNotificationTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	_, err = db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			telegram_id INTEGER NOT NULL UNIQUE,
			username TEXT NOT NULL,
			language_code TEXT NOT NULL DEFAULT '',
			chat_id INTEGER NOT NULL DEFAULT 0,
			status TEXT NOT NULL DEFAULT 'pending',
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		)`)
	if err != nil {
		t.Fatalf("create users table: %v", err)
	}

	return db
}

func insertNotificationTestUser(t *testing.T, db *sql.DB, telegramID int64, username string, chatID int64, status string) {
	t.Helper()

	now := time.Now().UTC()
	_, err := db.Exec(`
		INSERT INTO users (telegram_id, username, language_code, chat_id, status, created_at, updated_at)
		VALUES (?, ?, 'en', ?, ?, ?, ?)`,
		telegramID,
		username,
		chatID,
		status,
		now,
		now,
	)
	if err != nil {
		t.Fatalf("insert user %s: %v", username, err)
	}
}
