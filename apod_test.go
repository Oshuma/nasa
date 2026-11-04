package nasa

import (
	"encoding/json"
	"testing"
)

func apodJSON(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

const apodPlaceholderURL = "https://assets.science.nasa.gov/dynamicimage/assets/science/astro/programs/cosmic-origins/images/misc/news-thumbnail.png?w=594&h=516"

func apodBasicHTML(media string) string {
	return "<center>\n<h1> Astronomy Picture of the Day </h1>\n<p>\n<a href=\"https://science.nasa.gov/apod/archive/\">Discover the cosmos!</a>\n<p>\n2024 January 1<br>\n" +
		media + "\n</center>\n<center>\n<b> Title </b><a href=\"https://example.com/other.jpg\"><img src=\"https://example.com/other.jpg\"></a></center>"
}

func TestParseAPODContent(t *testing.T) {
	t.Run("single image", func(t *testing.T) {
		in := apodJSON(t, map[string]any{
			"date":        "2024-01-01",
			"title":       "NGC 1232",
			"permalink":   "https://science.nasa.gov/image-article/apod-2024-january-1/",
			"url":         "https://science.nasa.gov/image-article/apod-2024-january-1/",
			"media_type":  "image",
			"explanation": "<strong>Explanation:</strong> A galaxy.",
			"copyright":   "Someone",
			"hdurl":       "https://assets.science.nasa.gov/a/ngc1232.jpg?w=3969",
			"basic_html":  apodBasicHTML(`<a href="https://assets.science.nasa.gov/a/ngc1232.jpg"><IMG SRC="https://assets.science.nasa.gov/a/ngc1232.jpg" alt="x"></a>`),
		})
		items, err := parseAPODContent(in, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(items) != 1 {
			t.Fatalf("expected 1 item, got %d", len(items))
		}
		item := items[0]
		if item.Title != "NGC 1232" || item.MediaType != "image" || item.Copyright != "Someone" {
			t.Fatalf("unexpected item: %+v", item)
		}
		if item.URL != "https://assets.science.nasa.gov/a/ngc1232.jpg?w=3969" || item.HDURL != item.URL {
			t.Fatalf("expected hdurl for URL and HDURL, got %q / %q", item.URL, item.HDURL)
		}
		if item.Date.Format("2006-01-02") != "2024-01-01" {
			t.Fatalf("unexpected date: %s", item.Date.Format("2006-01-02"))
		}
	})

	t.Run("placeholder hdurl falls back to page image", func(t *testing.T) {
		in := apodJSON(t, map[string]any{
			"date":       "1995-06-16",
			"media_type": "image",
			"hdurl":      apodPlaceholderURL,
			"basic_html": apodBasicHTML(`<a href="https://assets.science.nasa.gov/a/big.gif?a=1&amp;b=2"><IMG SRC="https://assets.science.nasa.gov/a/small.gif" alt="x"></a>`),
		})
		items, err := parseAPODContent(in, false)
		if err != nil {
			t.Fatal(err)
		}
		if items[0].URL != "https://assets.science.nasa.gov/a/small.gif" {
			t.Errorf("URL = %q", items[0].URL)
		}
		if items[0].HDURL != "https://assets.science.nasa.gov/a/big.gif?a=1&b=2" {
			t.Errorf("HDURL = %q", items[0].HDURL)
		}
	})

	t.Run("youtube video", func(t *testing.T) {
		in := apodJSON(t, []map[string]any{{
			"date":       "2024-03-25",
			"media_type": "video",
			"hdurl":      apodPlaceholderURL,
			"basic_html": apodBasicHTML(`<iframe width="960" src="//www.youtube.com/embed/NqBfQeJqkfU?rel=0" allowfullscreen></iframe>`),
		}})

		items, err := parseAPODContent(in, false)
		if err != nil {
			t.Fatal(err)
		}
		if items[0].URL != "https://www.youtube.com/embed/NqBfQeJqkfU?rel=0" || items[0].HDURL != "" || items[0].ThumbnailURL != "" {
			t.Errorf("unexpected video: %+v", items[0])
		}

		items, err = parseAPODContent(in, true)
		if err != nil {
			t.Fatal(err)
		}
		if items[0].ThumbnailURL != "https://img.youtube.com/vi/NqBfQeJqkfU/0.jpg" {
			t.Errorf("ThumbnailURL = %q", items[0].ThumbnailURL)
		}
	})

	t.Run("multiple items", func(t *testing.T) {
		in := []byte(`[{"date":"2023-01-03","title":"Second","media_type":"image","hdurl":"https://example.com/2.jpg"},{"date":"2023-01-02","title":"First","media_type":"image","hdurl":"https://example.com/1.jpg"}]`)
		items, err := parseAPODContent(in, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(items) != 2 || items[1].URL != "https://example.com/1.jpg" {
			t.Fatalf("unexpected items: %+v", items)
		}
	})

	t.Run("empty response", func(t *testing.T) {
		items, err := parseAPODContent([]byte("\n"), false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if items != nil {
			t.Fatalf("expected nil slice, got %v", items)
		}
	})

	t.Run("invalid payload", func(t *testing.T) {
		if _, err := parseAPODContent([]byte(`{"date":`), false); err == nil {
			t.Fatalf("expected an error, got nil")
		}
	})
}

func TestAPODParamsType(t *testing.T) {
	if _, err := APOD(&MediaParams{}); err != ErrorParamsMismatch {
		t.Errorf("expected ErrorParamsMismatch, got %v", err)
	}
}
