package nasa

import "testing"

func TestParseAPODContent(t *testing.T) {
	t.Run("single item", func(t *testing.T) {
		in := []byte("{\"date\":\"2023-01-02\",\"title\":\"Sample\",\"url\":\"https://example.com\",\"thumbnail_url\":\"https://example.com/thumb.jpg\"}")
		items, err := parseAPODContent(in)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(items) != 1 {
			t.Fatalf("expected 1 item, got %d", len(items))
		}
		item := items[0]
		if item.Title != "Sample" {
			t.Fatalf("unexpected title: %s", item.Title)
		}
		if item.ThumbnailURL != "https://example.com/thumb.jpg" {
			t.Fatalf("unexpected thumbnail url: %s", item.ThumbnailURL)
		}
		if item.Date.Format("2006-01-02") != "2023-01-02" {
			t.Fatalf("unexpected date: %s", item.Date.Format("2006-01-02"))
		}
	})

	t.Run("multiple items", func(t *testing.T) {
		in := []byte("[{\"date\":\"2023-01-02\",\"title\":\"First\",\"url\":\"https://example.com/1\"},{\"date\":\"2023-01-03\",\"title\":\"Second\",\"url\":\"https://example.com/2\"}]")
		items, err := parseAPODContent(in)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(items) != 2 {
			t.Fatalf("expected 2 items, got %d", len(items))
		}
	})

	t.Run("empty response", func(t *testing.T) {
		in := []byte("\n")
		items, err := parseAPODContent(in)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if items != nil {
			t.Fatalf("expected nil slice, got %v", items)
		}
	})

	t.Run("invalid payload", func(t *testing.T) {
		in := []byte("{\"date\":")
		if _, err := parseAPODContent(in); err == nil {
			t.Fatalf("expected an error, got nil")
		}
	})
}
