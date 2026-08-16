package handler

import (
	"slices"
	"strings"
	"testing"
)

func TestPortfolioJSONIsValid(t *testing.T) {
	items := LoadAllPortfolioItems()
	if len(items) == 0 {
		t.Fatal("portfolio.json produced no entries")
	}

	for _, item := range items {
		if item.Description == "" {
			t.Errorf("%q: missing description", item.Title)
		}
		if len(item.Tags) == 0 {
			t.Errorf("%q: missing tags", item.Title)
		}
		// The site previously shipped six links to example.com.
		if strings.Contains(item.URL, "example.com") {
			t.Errorf("%q: placeholder URL %q", item.Title, item.URL)
		}
	}
}

func TestSearchPortfolio(t *testing.T) {
	tests := []struct {
		name      string
		query     string
		tag       string
		wantTitle string
		wantGone  string
	}{
		{name: "empty query returns everything", wantTitle: "Homelab"},
		{name: "tag only", tag: "certifications", wantTitle: "AWS Cloud Practitioner", wantGone: "Dotfiles"},
		{name: "query only", query: "dotfiles", wantTitle: "Dotfiles", wantGone: "Homelab"},
		{name: "query and tag combined", query: "aws", tag: "certifications", wantTitle: "AWS Cloud Practitioner", wantGone: "3 Tier Web App"},
		{name: "matches description", query: "terraform", wantTitle: "3 Tier Web App"},
		{name: "unknown tag falls back to all", tag: "nonsense", wantTitle: "Homelab"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SearchPortfolio(tt.query, tt.tag)
			titles := make([]string, len(got))
			for i, item := range got {
				titles[i] = item.Title
			}

			if tt.wantTitle != "" && !slices.Contains(titles, tt.wantTitle) {
				t.Errorf("want %q in results, got %v", tt.wantTitle, titles)
			}
			if tt.wantGone != "" && slices.Contains(titles, tt.wantGone) {
				t.Errorf("did not want %q in results, got %v", tt.wantGone, titles)
			}
		})
	}
}

func TestSearchPortfolioNoMatch(t *testing.T) {
	if got := SearchPortfolio("zzzzz-no-such-thing", "all"); len(got) != 0 {
		t.Errorf("want no results, got %d", len(got))
	}
}

func TestNewContactMsg(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		subject string
		message string
		wantErr bool
	}{
		{name: "valid", email: "a@b.com", subject: "hi", message: "hello", wantErr: false},
		{name: "trims whitespace", email: "  a@b.com ", subject: "hi", message: " hello ", wantErr: false},
		{name: "bad email", email: "nope", subject: "hi", message: "hello", wantErr: true},
		{name: "empty email", email: "", subject: "hi", message: "hello", wantErr: true},
		{name: "empty message", email: "a@b.com", subject: "hi", message: "   ", wantErr: true},
		{name: "subject too long", email: "a@b.com", subject: strings.Repeat("x", maxSubjectLen+1), message: "hello", wantErr: true},
		{name: "message too long", email: "a@b.com", subject: "hi", message: strings.Repeat("x", maxMessageLen+1), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewContactMsg(tt.email, tt.subject, tt.message)
			if tt.wantErr {
				if err == nil {
					t.Fatal("want error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if strings.TrimSpace(got.Email) != got.Email || strings.TrimSpace(got.Message) != got.Message {
				t.Error("fields were not trimmed")
			}
		})
	}
}

func TestSendDiscordMsgWithoutConfig(t *testing.T) {
	t.Setenv("DISCORD_BOT_TOKEN", "")
	t.Setenv("DISCORD_CHANNEL_ID", "")

	// Must return an error, not exit the process.
	if err := SendDiscordMsg(ContactMsg{Email: "a@b.com", Message: "hi"}); err == nil {
		t.Error("want an error when Discord is unconfigured, got nil")
	}
}

func TestImagePath(t *testing.T) {
	if got := (Portfolio{Image: ""}).ImagePath(); got != "" {
		t.Errorf("empty image should yield empty path, got %q", got)
	}
	if got := (Portfolio{Image: "crc.png"}).ImagePath(); got != imgPath+"crc.png" {
		t.Errorf("got %q", got)
	}
}
