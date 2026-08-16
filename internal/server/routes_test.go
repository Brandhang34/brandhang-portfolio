package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func newTestHandler() http.Handler {
	return (&Server{port: defaultPort}).RegisterRoutes()
}

func TestPagesRender(t *testing.T) {
	handler := newTestHandler()

	tests := []struct {
		name   string
		path   string
		expect string
	}{
		{"home", "/", "Brandon Hang"},
		{"about", "/about", "About Me"},
		{"portfolio", "/portfolio", "Serverless CI/CD Resume"},
		{"contact", "/contact", "Send Message"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, req)

			if resp.Code != http.StatusOK {
				t.Fatalf("GET %s: status = %d, want %d", tt.path, resp.Code, http.StatusOK)
			}
			if !strings.Contains(resp.Body.String(), tt.expect) {
				t.Errorf("GET %s: body does not contain %q", tt.path, tt.expect)
			}
		})
	}
}

func TestSearchPortfolioRoute(t *testing.T) {
	handler := newTestHandler()

	form := url.Values{}
	form.Set("search-portfolio", "terraform")
	form.Set("filter-tags", "certifications")

	req := httptest.NewRequest(http.MethodPost, "/search-portfolio", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusOK)
	}

	body := resp.Body.String()
	if !strings.Contains(body, "Hashicorp Terraform Associate") {
		t.Error("expected the Terraform certification in the results")
	}
	// Query and tag filter must apply together; previously each control posted
	// on its own and clobbered the other.
	if strings.Contains(body, "3 Tier Web App") {
		t.Error("a project leaked into a certifications-only search")
	}
}

func TestContactRejectsInvalidInput(t *testing.T) {
	handler := newTestHandler()

	form := url.Values{}
	form.Set("email", "not-an-email")
	form.Set("subject", "hi")
	form.Set("message", "hello")

	req := httptest.NewRequest(http.MethodPost, "/submit-contact-msg", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	// The endpoint must answer rather than terminate the process: it previously
	// called log.Fatal, so any failure here took the whole server down.
	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusOK)
	}
	if !strings.Contains(resp.Body.String(), "try again") {
		t.Error("expected a validation message in the response")
	}
}

func TestHoneypotIsAccepted(t *testing.T) {
	handler := newTestHandler()

	form := url.Values{}
	form.Set("email", "bot@example.com")
	form.Set("message", "spam")
	form.Set("website", "http://spam.example")

	req := httptest.NewRequest(http.MethodPost, "/submit-contact-msg", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusOK)
	}
	// A filled honeypot is silently dropped, so the bot cannot tell it failed
	// and Discord is never contacted.
	if !strings.Contains(resp.Body.String(), "on its way") {
		t.Error("expected the honeypot path to return the success fragment")
	}
}

func TestStaticAssetsAreServed(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/assets/css/output.css", nil)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusOK)
	}
}
