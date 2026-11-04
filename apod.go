package nasa

import (
	"bytes"
	"encoding/json"
	"errors"
	"html"
	"math/rand/v2"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// apodAPIURL is the WordPress APOD feed that replaced api.nasa.gov/planetary/apod.
const apodAPIURL = "https://science.nasa.gov/wp-json/wp/v2/apod-basic"

const (
	apodPerPage = 25

	// apodCountWorkers bounds concurrent requests when emulating Count.
	apodCountWorkers = 5

	// apodCountRounds is how many times Count retries dates that have no APOD.
	apodCountRounds = 3
)

// apodFirstDate is the date of the first Astronomy Picture of the Day.
var apodFirstDate = time.Date(1995, 6, 16, 0, 0, 0, 0, time.UTC)

// APODImage represents an Astronomy Picture Of the Day. Most callers obtain an APODImage
// via APOD or APODList rather than instantiating the struct directly.
//
// URL is the image (or video embed for videos) and HDURL the full-size image, matching
// the retired api.nasa.gov service. Explanation and Copyright may contain HTML.
type APODImage struct {
	Date        Date   `json:"date"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	HDURL       string `json:"hdurl"`
	Explanation string `json:"explanation"`
	MediaType   string `json:"media_type"`
	Copyright   string `json:"copyright"`

	// ThumbnailURL is only set for YouTube videos when APODParams.Thumbs is true.
	ThumbnailURL string `json:"thumbnail_url"`

	// Deprecated: The APOD feed has no service version; it is always empty.
	ServiceVersion string `json:"service_version"`
}

// apodItem is a single entry from the APOD WordPress feed.
type apodItem struct {
	Date        Date   `json:"date"`
	Title       string `json:"title"`
	Permalink   string `json:"permalink"`
	MediaType   string `json:"media_type"`
	Explanation string `json:"explanation"`
	Copyright   string `json:"copyright"`
	HDURL       string `json:"hdurl"`
	BasicHTML   string `json:"basic_html"`
}

// APOD returns the Astronomy Picture Of the Day for APODParams.Date, or the latest one when
// Date is not set. Other APODParams filters are ignored; see APODList.
func APOD(p ParamEncoder) (APODImage, error) {
	params, err := apodParams(p)
	if err != nil {
		return APODImage{}, err
	}

	images, err := apodSingle(params)
	if err != nil || len(images) == 0 {
		return APODImage{}, err
	}

	return images[0], nil
}

// APODList returns APOD images for a StartDate/EndDate range (oldest first) or Count random
// images. Count returns ErrorNoResults if no images are found. With neither set it behaves
// like APOD.
func APODList(p ParamEncoder) ([]APODImage, error) {
	params, err := apodParams(p)
	if err != nil {
		return nil, err
	}

	switch {
	case params.Count > 0:
		return apodRandom(params)
	case !params.StartDate.IsZero():
		return apodRange(params)
	default:
		return apodSingle(params)
	}
}

func apodParams(p ParamEncoder) (*APODParams, error) {
	switch params := p.(type) {
	case *APODParams:
		return params, nil
	case *APIParam, nil:
		return &APODParams{}, nil
	default:
		return nil, ErrorParamsMismatch
	}
}

// apodDateCode formats a date as the feed's legacy YYMMDD code.
func apodDateCode(t time.Time) string {
	return t.Format("060102")
}

func apodSingle(p *APODParams) ([]APODImage, error) {
	u := apodAPIURL + "?per_page=1"
	if !p.Date.IsZero() {
		u = apodAPIURL + "/" + apodDateCode(p.Date)
	}

	content, err := getContent(u, nil)
	if err != nil {
		return nil, err
	}

	return parseAPODContent(content, p.Thumbs)
}

func apodRange(p *APODParams) ([]APODImage, error) {
	end := p.EndDate
	if end.IsZero() {
		end = time.Now()
	}

	v := url.Values{}
	v.Set("date_from", apodDateCode(p.StartDate))
	v.Set("date_to", apodDateCode(end))
	v.Set("per_page", strconv.Itoa(apodPerPage))

	images := []APODImage{}
	for page := 1; ; page++ {
		v.Set("page", strconv.Itoa(page))
		content, header, err := getResponse(apodAPIURL + "?" + v.Encode())
		if err != nil {
			return nil, err
		}

		items, err := parseAPODContent(content, p.Thumbs)
		if err != nil {
			return nil, err
		}
		images = append(images, items...)

		pages, _ := strconv.Atoi(header.Get("X-WP-TotalPages"))
		if page >= pages || len(items) == 0 {
			break
		}
	}

	sort.SliceStable(images, func(i, j int) bool { return images[i].Date.Before(images[j].Date.Time) })
	return images, nil
}

// apodRandom emulates the retired count parameter by fetching random dates. Dates with no
// APOD are retried with new dates for a few rounds; any other error is returned.
func apodRandom(p *APODParams) ([]APODImage, error) {
	now := time.Now().UTC()
	days := int(now.Sub(apodFirstDate).Hours() / 24)
	tried := map[string]bool{}
	found := map[string]APODImage{}

	for round := 0; round < apodCountRounds && len(found) < p.Count; round++ {
		var dates []time.Time
		for attempts := 0; len(dates) < p.Count-len(found) && attempts < p.Count*10; attempts++ {
			d := apodFirstDate.AddDate(0, 0, rand.IntN(days+1))
			if code := apodDateCode(d); !tried[code] {
				tried[code] = true
				dates = append(dates, d)
			}
		}

		var (
			mu       sync.Mutex
			wg       sync.WaitGroup
			firstErr error
			sem      = make(chan struct{}, apodCountWorkers)
		)
		for _, d := range dates {
			wg.Add(1)
			sem <- struct{}{}
			go func(d time.Time) {
				defer func() { <-sem; wg.Done() }()
				images, err := apodSingle(&APODParams{Date: d, Thumbs: p.Thumbs})

				mu.Lock()
				defer mu.Unlock()
				var statusErr *ErrorHTTPStatus
				switch {
				case errors.As(err, &statusErr) && statusErr.StatusCode == http.StatusNotFound:
				case err != nil:
					if firstErr == nil {
						firstErr = err
					}
				case len(images) > 0:
					found[apodDateCode(d)] = images[0]
				}
			}(d)
		}
		wg.Wait()

		if firstErr != nil {
			return nil, firstErr
		}
	}

	if len(found) == 0 {
		return nil, ErrorNoResults
	}

	images := make([]APODImage, 0, len(found))
	for _, img := range found {
		images = append(images, img)
	}
	return images, nil
}

func parseAPODContent(content []byte, thumbs bool) ([]APODImage, error) {
	trimmed := bytes.TrimSpace(content)
	if len(trimmed) == 0 {
		return nil, nil
	}

	var items []apodItem
	if trimmed[0] == '{' {
		item := apodItem{}
		if err := json.Unmarshal(trimmed, &item); err != nil {
			return nil, err
		}
		items = []apodItem{item}
	} else if err := json.Unmarshal(trimmed, &items); err != nil {
		return nil, err
	}

	images := make([]APODImage, 0, len(items))
	for _, item := range items {
		images = append(images, item.image(thumbs))
	}
	return images, nil
}

var (
	apodLinkedImageRe = regexp.MustCompile(`(?is)<a\s[^>]*href="([^"]+)"[^>]*>\s*<img\s[^>]*src="([^"]+)"`)
	apodImageRe       = regexp.MustCompile(`(?is)<img\s[^>]*src="([^"]+)"`)
	apodIframeRe      = regexp.MustCompile(`(?is)<iframe\s[^>]*src="([^"]+)"`)
	apodYouTubeRe     = regexp.MustCompile(`youtube(?:-nocookie)?\.com/embed/([\w-]+)`)
)

func (item apodItem) image(thumbs bool) APODImage {
	img := APODImage{
		Date:        item.Date,
		Title:       item.Title,
		Explanation: item.Explanation,
		MediaType:   item.MediaType,
		Copyright:   item.Copyright,
	}

	hd := item.HDURL
	if apodPlaceholder(hd) {
		hd = ""
	}

	// basic_html mirrors the legacy APOD page; its first <center> block holds the media.
	media := item.BasicHTML
	if i := strings.Index(strings.ToLower(media), "</center>"); i >= 0 {
		media = media[:i]
	}

	if item.MediaType == "image" {
		var link, src string
		if m := apodLinkedImageRe.FindStringSubmatch(media); m != nil {
			link, src = m[1], m[2]
		} else if m := apodImageRe.FindStringSubmatch(media); m != nil {
			src = m[1]
		}
		img.URL = firstNonEmpty(hd, apodURL(src), apodURL(link))
		img.HDURL = firstNonEmpty(hd, apodURL(link), apodURL(src))
		return img
	}

	if m := apodIframeRe.FindStringSubmatch(media); m != nil {
		img.URL = apodURL(m[1])
		if yt := apodYouTubeRe.FindStringSubmatch(img.URL); thumbs && yt != nil {
			img.ThumbnailURL = "https://img.youtube.com/vi/" + yt[1] + "/0.jpg"
		}
	} else {
		img.URL = hd
	}
	return img
}

// apodPlaceholder reports whether u is the generic image the feed uses when an APOD has
// no featured image.
func apodPlaceholder(u string) bool {
	return strings.Contains(u, "/images/misc/news-thumbnail.")
}

func apodURL(u string) string {
	u = strings.TrimSpace(html.UnescapeString(u))
	if strings.HasPrefix(u, "//") {
		u = "https:" + u
	}
	return u
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
