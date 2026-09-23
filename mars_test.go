package nasa

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestHasCamera(t *testing.T) {
	t.Run("has camera", func(t *testing.T) {
		r := RoverCuriosity
		c := RoverCameraFHAZ
		if !hasCamera(r, c) {
			t.Errorf("rover %s should have camera %s", r.Name, c.Name)
		}
	})

	t.Run("perseverance has camera", func(t *testing.T) {
		r := RoverPerseverance
		c := RoverCameraMCZLEFT
		if !hasCamera(r, c) {
			t.Errorf("rover %s should have camera %s", r.Name, c.Name)
		}
	})

	t.Run("does not have camera", func(t *testing.T) {
		r := RoverCuriosity
		c := RoverCameraPANCAM
		if hasCamera(r, c) {
			t.Errorf("rover %s should not have camera %s", r.Name, c.Name)
		}
	})

	t.Run("perseverance missing camera", func(t *testing.T) {
		r := RoverPerseverance
		c := RoverCameraMAST
		if hasCamera(r, c) {
			t.Errorf("rover %s should not have camera %s", r.Name, c.Name)
		}
	})
}

func TestMarsParamsMismatch(t *testing.T) {
	p := &APIParam{APIKey: "key"}

	if _, err := MarsRoverPhotos(p, RoverCuriosity); err != ErrorParamsMismatch {
		t.Errorf("MarsRoverPhotos: expected ErrorParamsMismatch, got: %v", err)
	}

	if _, err := MarsRoverPhotosLatest(p, RoverCuriosity); err != ErrorParamsMismatch {
		t.Errorf("MarsRoverPhotosLatest: expected ErrorParamsMismatch, got: %v", err)
	}

	if _, err := MarsMissionManifest(p, RoverCuriosity); err != ErrorParamsMismatch {
		t.Errorf("MarsMissionManifest: expected ErrorParamsMismatch, got: %v", err)
	}
}

// fakeCuriosity serves raw_image_items responses from items grouped by sol, honouring the
// sol and instrument conditions and paginating with perPage.
func fakeCuriosity(t *testing.T, latestSol int, bySol map[int][]map[string]any, perPage int) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("condition_1") != "msl:mission" {
			t.Errorf("missing mission condition: %s", r.URL.RawQuery)
		}

		if q.Get("condition_2") == "" {
			json.NewEncoder(w).Encode(map[string]any{
				"items": []map[string]any{{"sol": latestSol}},
				"total": 12345,
			})
			return
		}

		var sol int
		fmt.Sscanf(q.Get("condition_2"), "%d:sol:in", &sol)
		items := bySol[sol]

		if c := q.Get("condition_3"); c != "" {
			wanted := strings.Split(strings.TrimSuffix(c, ":instrument:in"), ",")
			filtered := []map[string]any{}
			for _, i := range items {
				if slices.Contains(wanted, i["instrument"].(string)) {
					filtered = append(filtered, i)
				}
			}
			items = filtered
		}

		page, _ := strconv.Atoi(q.Get("page"))
		start := min(page*perPage, len(items))
		end := min(start+perPage, len(items))
		json.NewEncoder(w).Encode(map[string]any{
			"items": items[start:end],
			"more":  end < len(items),
			"total": len(items),
		})
	}))
	t.Cleanup(srv.Close)

	orig := curiosityRawImagesURL
	curiosityRawImagesURL = srv.URL + "/"
	t.Cleanup(func() { curiosityRawImagesURL = orig })
}

func curiosityItem(id, sol int, instrument, sampleType, taken string) map[string]any {
	return map[string]any{
		"id":         id,
		"imageid":    fmt.Sprintf("IMG%d", id),
		"sol":        sol,
		"instrument": instrument,
		"https_url":  fmt.Sprintf("https://example.com/%d.jpg", id),
		"date_taken": taken,
		"extended":   map[string]any{"sample_type": sampleType},
	}
}

func TestMarsRoverPhotosCuriosity(t *testing.T) {
	fakeCuriosity(t, 1000, map[int][]map[string]any{
		1000: {
			curiosityItem(1, 1000, "NAV_LEFT_B", "full", "2015-05-30T15:00:00.000Z"),
			curiosityItem(2, 1000, "FHAZ_RIGHT_B", "full", "2015-05-30T16:00:00.000Z"),
			curiosityItem(3, 1000, "MAST_LEFT", "subframe", "2015-05-30T14:00:00.000Z"),
			curiosityItem(4, 1000, "FHAZ_LEFT_B", "full", "2015-05-30T13:00:00.000Z"),
			curiosityItem(5, 1000, "MAST_RIGHT", "thumbnail", "2015-05-30T14:00:00.000Z"),
			curiosityItem(6, 1000, "NEW_CAM", "full", "2015-05-30T12:00:00.000Z"),
		},
	}, 2)

	t.Run("sol, full frames only, sorted by camera then time", func(t *testing.T) {
		r, err := MarsRoverPhotos(&MarsPhotosParams{Sol: 1000}, RoverCuriosity)
		if err != nil {
			t.Fatal(err)
		}

		got := []string{}
		for _, p := range r.Photos {
			got = append(got, p.ImageID+"/"+p.Camera.Name)
		}
		expected := []string{"IMG4/FHAZ", "IMG2/FHAZ", "IMG1/NAVCAM", "IMG6/NEW_CAM"}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("\nexpected: %v\ngot: %v", expected, got)
		}
		if r.Total != 4 || r.Page != 1 {
			t.Errorf("expected total 4 page 1, got total %d page %d", r.Total, r.Page)
		}

		p := r.Photos[0]
		if p.ID != 4 || p.Sol != 1000 || p.Image != "https://example.com/4.jpg" || p.Camera.Instrument != "FHAZ_LEFT_B" ||
			p.Camera.FullName != RoverCameraFHAZ.FullName || p.EarthDate.Format("2006-01-02") != "2015-05-30" ||
			!p.TakenAt.Equal(time.Date(2015, 5, 30, 13, 0, 0, 0, time.UTC)) || p.Rover.Name != "Curiosity" ||
			p.Rover.Status != "active" || len(p.Rover.Cameras) != len(RoverCuriosity.Cameras) {
			t.Errorf("photo fields not mapped: %+v", p)
		}
	})

	t.Run("camera", func(t *testing.T) {
		r, err := MarsRoverPhotos(&MarsPhotosParams{Sol: 1000, Camera: RoverCameraFHAZ}, RoverCuriosity)
		if err != nil {
			t.Fatal(err)
		}
		if r.Total != 2 || r.Photos[0].ImageID != "IMG4" || r.Photos[1].ImageID != "IMG2" {
			t.Errorf("unexpected photos: %+v", r.Photos)
		}
	})

	t.Run("earth date uses capture time", func(t *testing.T) {
		r, err := MarsRoverPhotos(&MarsPhotosParams{EarthDate: time.Date(2015, 5, 30, 0, 0, 0, 0, time.UTC)}, RoverCuriosity)
		if err != nil {
			t.Fatal(err)
		}
		if r.Total != 4 {
			t.Errorf("expected 4 photos, got %d", r.Total)
		}

		r, err = MarsRoverPhotos(&MarsPhotosParams{EarthDate: time.Date(2015, 5, 31, 0, 0, 0, 0, time.UTC)}, RoverCuriosity)
		if err != nil {
			t.Fatal(err)
		}
		if r.Total != 0 {
			t.Errorf("expected no photos on the following day, got %d", r.Total)
		}
	})

	t.Run("unknown camera", func(t *testing.T) {
		_, err := MarsRoverPhotos(&MarsPhotosParams{Sol: 1000, Camera: RoverCameraPANCAM}, RoverCuriosity)
		var missing *ErrorRoverCameraMissing
		if !errors.As(err, &missing) {
			t.Errorf("expected ErrorRoverCameraMissing, got: %v", err)
		}
	})
}

func TestMarsRoverPhotosPagination(t *testing.T) {
	items := []map[string]any{}
	for i := 1; i <= 60; i++ {
		items = append(items, curiosityItem(i, 5, "MAHLI", "full", fmt.Sprintf("2020-01-01T00:%02d:00.000Z", i%60)))
	}
	fakeCuriosity(t, 5, map[int][]map[string]any{5: items}, 1000)

	tests := []struct {
		page, expected int
	}{{0, 60}, {1, 25}, {2, 25}, {3, 10}, {4, 0}}

	for _, tt := range tests {
		r, err := MarsRoverPhotos(&MarsPhotosParams{Sol: 5, Page: tt.page}, RoverCuriosity)
		if err != nil {
			t.Fatal(err)
		}
		if len(r.Photos) != tt.expected || r.Total != 60 {
			t.Errorf("page %d: expected %d of 60, got %d of %d", tt.page, tt.expected, len(r.Photos), r.Total)
		}
	}
}

func TestMarsRoverPhotosLatest(t *testing.T) {
	fakeCuriosity(t, 12, map[int][]map[string]any{
		12: {curiosityItem(1, 12, "MAHLI", "thumbnail", "2020-01-03T00:00:00.000Z")},
		11: {
			curiosityItem(2, 11, "MAHLI", "full", "2020-01-02T00:00:00.000Z"),
			curiosityItem(3, 11, "NAV_LEFT_B", "full", "2020-01-02T01:00:00.000Z"),
		},
	}, 1000)

	photos, err := MarsRoverPhotosLatest(&MarsPhotosParams{}, RoverCuriosity)
	if err != nil {
		t.Fatal(err)
	}
	if len(photos) != 2 || photos[0].Sol != 11 {
		t.Errorf("expected the 2 photos from sol 11, got: %+v", photos)
	}

	photos, err = MarsRoverPhotosLatest(&MarsPhotosParams{Camera: RoverCameraNAVCAM}, RoverCuriosity)
	if err != nil {
		t.Fatal(err)
	}
	if len(photos) != 1 || photos[0].ImageID != "IMG3" {
		t.Errorf("expected the NAVCAM photo, got: %+v", photos)
	}

	m, err := MarsMissionManifest(&MarsPhotosParams{}, RoverCuriosity)
	if err != nil {
		t.Fatal(err)
	}
	if m.Name != "Curiosity" || m.MaxSol != 11 || m.MaxDate.Format("2006-01-02") != "2020-01-02" ||
		m.TotalPhotos != 12345 || m.LandingDate.Format("2006-01-02") != "2012-08-06" || m.Photos != nil {
		t.Errorf("unexpected manifest: %+v", m)
	}
}

func TestMarsRoverPhotosPerseverance(t *testing.T) {
	var searches []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("category") != "mars2020" || q.Get("order") == "" {
			t.Errorf("unexpected query: %s", r.URL.RawQuery)
		}
		if q.Get("condition_2") != "100:sol:gte" || q.Get("condition_3") != "100:sol:lte" {
			t.Errorf("unexpected sol conditions: %s", r.URL.RawQuery)
		}
		searches = append(searches, q.Get("search"))

		// Serve one full page followed by a short one to exercise pagination.
		images := []map[string]any{}
		n := perseverancePerPage
		if q.Get("page") != "0" {
			n = 3
		}
		for i := 0; i < n; i++ {
			images = append(images, map[string]any{
				"imageid":        fmt.Sprintf("P%s-%03d", q.Get("page"), i),
				"sol":            100,
				"sample_type":    "Full",
				"camera":         map[string]any{"instrument": "MCZ_LEFT"},
				"image_files":    map[string]any{"large": "https://example.com/p.jpg"},
				"date_taken_utc": "2021-06-01T16:00:54.287",
			})
		}
		json.NewEncoder(w).Encode(map[string]any{"images": images})
	}))
	defer srv.Close()

	orig := perseveranceRawImagesURL
	perseveranceRawImagesURL = srv.URL + "/"
	defer func() { perseveranceRawImagesURL = orig }()

	r, err := MarsRoverPhotos(&MarsPhotosParams{Sol: 100, Camera: RoverCameraMCZLEFT, Page: 5}, RoverPerseverance)
	if err != nil {
		t.Fatal(err)
	}

	if r.Total != perseverancePerPage+3 || len(r.Photos) != 3 {
		t.Errorf("expected page 5 of %d photos to hold 3, got %d of %d", perseverancePerPage+3, len(r.Photos), r.Total)
	}
	if !reflect.DeepEqual(searches, []string{"|MCZ_LEFT", "|MCZ_LEFT"}) {
		t.Errorf("unexpected search params: %q", searches)
	}

	p := r.Photos[0]
	if p.ID != 0 || p.Camera.Name != "MCZ_LEFT" || !p.TakenAt.Equal(time.Date(2021, 6, 1, 16, 0, 54, 287000000, time.UTC)) {
		t.Errorf("photo fields not mapped: %+v", p)
	}
}

func TestMarsRoverUnsupported(t *testing.T) {
	// A caller-built Rover with no image source, e.g. a mistyped slug or an unlisted rover.
	rover := Rover{Name: "Sojourner", Slug: "sojourner"}

	if _, err := MarsRoverPhotos(&MarsPhotosParams{}, rover); err != ErrorRoverUnsupported {
		t.Errorf("MarsRoverPhotos: expected ErrorRoverUnsupported, got: %v", err)
	}
	if _, err := MarsRoverPhotosLatest(&MarsPhotosParams{}, rover); err != ErrorRoverUnsupported {
		t.Errorf("MarsRoverPhotosLatest: expected ErrorRoverUnsupported, got: %v", err)
	}
	if _, err := MarsMissionManifest(&MarsPhotosParams{}, rover); err != ErrorRoverUnsupported {
		t.Errorf("MarsMissionManifest: expected ErrorRoverUnsupported, got: %v", err)
	}
}

func TestSolsOn(t *testing.T) {
	m := marsMissions[RoverPerseverance.Slug]

	tests := []struct {
		date     time.Time
		expected []int
	}{
		// Sol 100 runs from 2021-05-31 22:22 to 2021-06-01 23:01 UTC.
		{time.Date(2021, 6, 1, 0, 0, 0, 0, time.UTC), []int{100, 101}},
		{time.Date(2021, 5, 31, 12, 0, 0, 0, time.UTC), []int{99, 100}},
		{time.Date(2021, 2, 18, 0, 0, 0, 0, time.UTC), []int{0}},
		{time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), []int{}},
	}

	for _, tt := range tests {
		if got := m.solsOn(tt.date); !reflect.DeepEqual(got, tt.expected) {
			t.Errorf("%s: expected %v, got %v", tt.date.Format("2006-01-02"), tt.expected, got)
		}
	}
}

func TestUTCTimeUnmarshalJSON(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Time
	}{
		{`"2021-06-01T16:00:54.287"`, time.Date(2021, 6, 1, 16, 0, 54, 287000000, time.UTC)},
		{`"2021-02-22T09:57:38Z"`, time.Date(2021, 2, 22, 9, 57, 38, 0, time.UTC)},
		{`"2021-02-22T11:57:38+02:00"`, time.Date(2021, 2, 22, 9, 57, 38, 0, time.UTC)},
		{`null`, time.Time{}},
	}

	for _, tt := range tests {
		var got utcTime
		if err := json.Unmarshal([]byte(tt.input), &got); err != nil {
			t.Errorf("%s: %v", tt.input, err)
			continue
		}
		if !got.Equal(tt.expected) {
			t.Errorf("%s: expected %s, got %s", tt.input, tt.expected, got.Time)
		}
	}
}

// fakePDS serves PDS Atlas search responses from docs, applying the sol, date and camera
// filters and paginating with rows/start.
func fakePDS(t *testing.T, docs []map[string]any) *[]url.Values {
	t.Helper()
	var queries []url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		queries = append(queries, q)

		fqs := q["fq"]
		if !slices.Contains(fqs, `ATLAS_MISSION_NAME:"mars exploration rover"`) || !slices.Contains(fqs, "ATLAS_SPACECRAFT_NAME:opportunity") {
			t.Errorf("missing mission filters: %v", fqs)
		}

		matched := []map[string]any{}
		for _, d := range docs {
			ok := true
			for _, fq := range fqs {
				switch {
				case strings.HasPrefix(fq, "PLANET_DAY_NUMBER:"):
					ok = ok && fq == fmt.Sprintf("PLANET_DAY_NUMBER:%d", d["PLANET_DAY_NUMBER"])
				case strings.HasPrefix(fq, "START_TIME:["):
					var from, to string
					fmt.Sscanf(strings.NewReplacer("[", "", "}", "", " TO ", " ").Replace(strings.TrimPrefix(fq, "START_TIME:")), "%s %s", &from, &to)
					st := d["START_TIME"].(string)
					ok = ok && st >= from && st < to
				case strings.HasPrefix(fq, "FILE_NAME:"):
					letter := d["FILE_NAME"].(string)[1:2]
					ok = ok && (strings.Contains(fq, "FILE_NAME:??") || strings.Contains(fq, "FILE_NAME:?"+letter+"?"))
				}
			}
			if ok {
				matched = append(matched, d)
			}
		}

		if q.Get("sort") == "PLANET_DAY_NUMBER desc" {
			sort.Slice(matched, func(i, j int) bool {
				return matched[i]["PLANET_DAY_NUMBER"].(int) > matched[j]["PLANET_DAY_NUMBER"].(int)
			})
		}

		rows, _ := strconv.Atoi(q.Get("rows"))
		start, _ := strconv.Atoi(q.Get("start"))
		page := matched[min(start, len(matched)):min(start+rows, len(matched))]

		json.NewEncoder(w).Encode(map[string]any{
			"response": map[string]any{"numFound": len(matched), "docs": page},
		})
	}))
	t.Cleanup(srv.Close)

	orig := pdsImagingSearchURL
	pdsImagingSearchURL = srv.URL
	t.Cleanup(func() { pdsImagingSearchURL = orig })
	return &queries
}

func pdsDoc(fileName string, sol int, start string) map[string]any {
	return map[string]any{
		"FILE_NAME":         fileName,
		"PRODUCT_ID":        strings.ToUpper(strings.TrimSuffix(fileName, ".img")),
		"PLANET_DAY_NUMBER": sol,
		"START_TIME":        start,
		"ATLAS_BROWSE_URL":  "https://pds.example.com///data/mer/mer1no_0xxx/extras/catalog/sol0100/rdr//" + fileName + ".jpg",
	}
}

func TestMarsRoverPhotosOpportunity(t *testing.T) {
	queries := fakePDS(t, []map[string]any{
		pdsDoc("1p137064657ilf2019p2356l7m1.img", 100, "2004-05-05T21:30:23Z"),
		pdsDoc("1n137056144ilf2002p1825r0m1.img", 100, "2004-05-05T19:08:30Z"),
		pdsDoc("1f137064531ilf2019p1214r0m1.img", 100, "2004-05-05T21:28:17Z"),
		pdsDoc("1r137064581ilf2019p1312l0m1.img", 100, "2004-05-05T21:29:08Z"),
		pdsDoc("1m137060845ilf2002p2956m2m1.img", 100, "2004-05-05T20:26:51Z"),
		pdsDoc("1n137150000ilf2100p1825r0m1.img", 101, "2004-05-06T22:00:00Z"),
	})

	t.Run("sol, sorted by camera", func(t *testing.T) {
		r, err := MarsRoverPhotos(&MarsPhotosParams{Sol: 100}, RoverOpportunity)
		if err != nil {
			t.Fatal(err)
		}

		got := []string{}
		for _, p := range r.Photos {
			got = append(got, p.Camera.Name+"/"+p.Camera.Instrument)
		}
		expected := []string{"FHAZ/FRONT_HAZCAM", "RHAZ/REAR_HAZCAM", "NAVCAM/NAVCAM", "PANCAM/PANCAM", "MI/MI"}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("\nexpected: %v\ngot: %v", expected, got)
		}

		p := r.Photos[0]
		if p.ID != 0 || p.ImageID != "1F137064531ILF2019P1214R0M1" || p.Sol != 100 ||
			p.Image != "https://pds.example.com/data/mer/mer1no_0xxx/browse/sol0100/rdr/1f137064531ilf2019p1214r0m1.img.jpg" ||
			!p.TakenAt.Equal(time.Date(2004, 5, 5, 21, 28, 17, 0, time.UTC)) || p.Rover.Status != "complete" {
			t.Errorf("photo fields not mapped: %+v", p)
		}

		last := (*queries)[len(*queries)-1]
		if !slices.Contains(last["fq"], "PLANET_DAY_NUMBER:100") || !slices.Contains(last["fq"], "FILE_NAME:???????????ilf*") {
			t.Errorf("unexpected filters: %v", last["fq"])
		}
	})

	t.Run("camera", func(t *testing.T) {
		r, err := MarsRoverPhotos(&MarsPhotosParams{Sol: 100, Camera: RoverCameraRHAZ}, RoverOpportunity)
		if err != nil {
			t.Fatal(err)
		}
		if r.Total != 1 || r.Photos[0].Camera.Name != "RHAZ" {
			t.Errorf("unexpected photos: %+v", r.Photos)
		}

		last := (*queries)[len(*queries)-1]
		if !slices.Contains(last["fq"], "FILE_NAME:?r?????????ilf*") {
			t.Errorf("unexpected filters: %v", last["fq"])
		}
	})

	t.Run("earth date is filtered by the index", func(t *testing.T) {
		r, err := MarsRoverPhotos(&MarsPhotosParams{EarthDate: time.Date(2004, 5, 6, 0, 0, 0, 0, time.UTC)}, RoverOpportunity)
		if err != nil {
			t.Fatal(err)
		}
		if r.Total != 1 || r.Photos[0].Sol != 101 {
			t.Errorf("unexpected photos: %+v", r.Photos)
		}

		last := (*queries)[len(*queries)-1]
		if !slices.Contains(last["fq"], "START_TIME:[2004-05-06T00:00:00Z TO 2004-05-07T00:00:00Z}") {
			t.Errorf("unexpected filters: %v", last["fq"])
		}
	})

	t.Run("latest and manifest", func(t *testing.T) {
		photos, err := MarsRoverPhotosLatest(&MarsPhotosParams{}, RoverOpportunity)
		if err != nil {
			t.Fatal(err)
		}
		if len(photos) != 1 || photos[0].Sol != 101 {
			t.Errorf("expected the sol 101 photo, got: %+v", photos)
		}

		m, err := MarsMissionManifest(&MarsPhotosParams{}, RoverOpportunity)
		if err != nil {
			t.Fatal(err)
		}
		if m.MaxSol != 101 || m.MaxDate.Format("2006-01-02") != "2004-05-06" || m.TotalPhotos != 6 || m.Status != "complete" {
			t.Errorf("unexpected manifest: %+v", m)
		}
	})
}

func TestPDSPagination(t *testing.T) {
	docs := []map[string]any{}
	for i := 0; i < pdsPerPage+5; i++ {
		docs = append(docs, pdsDoc(fmt.Sprintf("1n%09dilf2002p1825r0m1.img", i), 7, "2004-02-01T00:00:00Z"))
	}
	queries := fakePDS(t, docs)

	r, err := MarsRoverPhotos(&MarsPhotosParams{Sol: 7}, RoverOpportunity)
	if err != nil {
		t.Fatal(err)
	}
	if r.Total != pdsPerPage+5 || len(*queries) != 2 {
		t.Errorf("expected %d photos over 2 requests, got %d over %d", pdsPerPage+5, r.Total, len(*queries))
	}
}

func TestPDSFullSizeURL(t *testing.T) {
	tests := []struct{ input, expected string }{
		{
			"https://pds-imaging.jpl.nasa.gov///data/mer/mer1no_0xxx/extras/catalog/sol0100/rdr//1n137056144ilf2002p1825r0m1.img.jpg",
			"https://pds-imaging.jpl.nasa.gov/data/mer/mer1no_0xxx/browse/sol0100/rdr/1n137056144ilf2002p1825r0m1.img.jpg",
		},
		{"https://example.com/other/path.jpg", "https://example.com/other/path.jpg"},
	}
	for _, tt := range tests {
		if got := pdsFullSizeURL(tt.input); got != tt.expected {
			t.Errorf("\nexpected: %s\ngot: %s", tt.expected, got)
		}
	}
}
