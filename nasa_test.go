package nasa

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetContent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"ok":true}`))
		}))
		defer srv.Close()

		content, err := getContent(srv.URL, nil)
		if err != nil {
			t.Fatal(err)
		}
		if string(content) != `{"ok":true}` {
			t.Errorf("unexpected content: %s", content)
		}
	})

	t.Run("non-2xx status", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error":"bad key"}`))
		}))
		defer srv.Close()

		_, err := getContent(srv.URL, nil)
		var statusErr *ErrorHTTPStatus
		if !errors.As(err, &statusErr) {
			t.Fatalf("expected *ErrorHTTPStatus, got: %v", err)
		}
		if statusErr.StatusCode != http.StatusForbidden {
			t.Errorf("expected status %d, got: %d", http.StatusForbidden, statusErr.StatusCode)
		}
		if statusErr.Body != `{"error":"bad key"}` {
			t.Errorf("unexpected body: %s", statusErr.Body)
		}
	})
}
