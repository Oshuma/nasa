package nasa

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestStringListUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected StringList
	}{
		{"array", `["Apollo","Apollo 11"]`, StringList{"Apollo", "Apollo 11"}},
		{"empty array", `[]`, StringList{}},
		{"single string", `"Hubble Space Telescope"`, StringList{"Hubble Space Telescope"}},
		{"empty string", `""`, nil},
		{"null", `null`, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var l StringList
			if err := json.Unmarshal([]byte(tt.input), &l); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(l, tt.expected) {
				t.Errorf("expected: %#v, got: %#v", tt.expected, l)
			}
		})
	}

	t.Run("invalid type", func(t *testing.T) {
		var l StringList
		if err := json.Unmarshal([]byte(`42`), &l); err == nil {
			t.Error("expected an error for a number")
		}
	})
}

func TestMediaMetadataUnmarshalJSON(t *testing.T) {
	// Trimmed from live images-assets.nasa.gov metadata.json documents, which use
	// either shape for AVAIL:Album and IPTC:Keywords.
	tests := []struct {
		name     string
		input    string
		album    StringList
		keywords StringList
	}{
		{
			name: "album array, keywords array",
			input: `{
				"AVAIL:NASAID": "jsc2007e034221",
				"AVAIL:Title": "Apollo 11 spacecraft pre-launch",
				"AVAIL:Album": ["KSC_50th_Anniversary"],
				"IPTC:Keywords": ["Apollo", "Apollo 11", "Launch"],
				"File:ImageWidth": 2341,
				"Composite:Megapixels": 7.2
			}`,
			album:    StringList{"KSC_50th_Anniversary"},
			keywords: StringList{"Apollo", "Apollo 11", "Launch"},
		},
		{
			name: "album empty string, keywords single string",
			input: `{
				"AVAIL:NASAID": "PIA12235",
				"AVAIL:Title": "Nearside of the Moon",
				"AVAIL:Album": "",
				"IPTC:Keywords": "Moon",
				"File:ImageWidth": 2000,
				"Composite:Megapixels": 2.8
			}`,
			album:    nil,
			keywords: StringList{"Moon"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m MediaMetadata
			if err := json.Unmarshal([]byte(tt.input), &m); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(m.AVAILAlbum, tt.album) {
				t.Errorf("AVAILAlbum expected: %#v, got: %#v", tt.album, m.AVAILAlbum)
			}
			if !reflect.DeepEqual(m.IPTCKeywords, tt.keywords) {
				t.Errorf("IPTCKeywords expected: %#v, got: %#v", tt.keywords, m.IPTCKeywords)
			}
			if m.AVAILTitle == "" || m.FileImageWidth == 0 || m.CompositeMegapixels == 0 {
				t.Errorf("other fields not decoded: %+v", m)
			}
		})
	}
}
