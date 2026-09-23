package nasa

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

// Base URLs for NASA's raw image feeds. These are variables so tests can point them at a
// local server.
var (
	curiosityRawImagesURL    = "https://mars.nasa.gov/api/v1/raw_image_items/"
	perseveranceRawImagesURL = "https://mars.nasa.gov/rss/api/"
	pdsImagingSearchURL      = "https://pds-imaging.jpl.nasa.gov/solr/pds_archives/search"
)

const (
	// secondsPerSol is the length of a Martian solar day in Earth seconds.
	secondsPerSol = 88775.244

	curiosityPerPage    = 1000
	perseverancePerPage = 100 // The feed caps page sizes at 100.
	pdsPerPage          = 1000
)

// marsSource fetches full-frame photos from a rover's raw image feed.
type marsSource interface {
	// fetchSol returns every full-frame photo taken on sol. When instruments is not
	// empty, only photos from those source instruments are returned.
	fetchSol(sol int, instruments []string) ([]sourcePhoto, error)

	// latestSol returns the most recent sol with images of any kind.
	latestSol() (int, error)

	// totalPhotos returns the total number of images the feed reports for the mission.
	totalPhotos() (int, error)
}

// dateSource is implemented by sources that can filter by capture date themselves, which
// saves fetching whole sols.
type dateSource interface {
	// fetchDate returns every full-frame photo captured on the UTC day starting at day.
	fetchDate(day time.Time, instruments []string) ([]sourcePhoto, error)
}

// sourcePhoto is a photo as reported by a raw image feed, before it is mapped onto a
// RoverPhoto.
type sourcePhoto struct {
	id         int
	imageID    string
	sol        int
	instrument string
	url        string
	takenAt    time.Time
}

// marsMission ties a rover to its raw image feed.
type marsMission struct {
	source marsSource

	// solEpoch is the UTC start of sol 0 (local mean midnight at the landing site). It is
	// only needed when the source is not a dateSource.
	solEpoch time.Time

	// cameras maps RoverCamera slugs to the instrument names the feed uses.
	cameras map[string][]string
}

var marsMissions = map[string]marsMission{
	RoverCuriosity.Slug: {
		source:   curiositySource{},
		solEpoch: time.Date(2012, 8, 5, 13, 52, 34, 0, time.UTC),
		cameras: map[string][]string{
			RoverCameraFHAZ.Slug:    {"FHAZ_LEFT_A", "FHAZ_RIGHT_A", "FHAZ_LEFT_B", "FHAZ_RIGHT_B"},
			RoverCameraRHAZ.Slug:    {"RHAZ_LEFT_A", "RHAZ_RIGHT_A", "RHAZ_LEFT_B", "RHAZ_RIGHT_B"},
			RoverCameraMAST.Slug:    {"MAST_LEFT", "MAST_RIGHT"},
			RoverCameraCHEMCAM.Slug: {"CHEMCAM_RMI"},
			RoverCameraMAHLI.Slug:   {"MAHLI"},
			RoverCameraMARDI.Slug:   {"MARDI"},
			RoverCameraNAVCAM.Slug:  {"NAV_LEFT_A", "NAV_RIGHT_A", "NAV_LEFT_B", "NAV_RIGHT_B"},
		},
	},
	RoverPerseverance.Slug: {
		source:   perseveranceSource{},
		solEpoch: time.Date(2021, 2, 18, 4, 24, 16, 0, time.UTC),
		cameras:  identityCameraMap(RoverPerseverance),
	},
	RoverOpportunity.Slug: {
		source:  pdsSource{spacecraft: "opportunity"},
		cameras: merCameras,
	},
	RoverSpirit.Slug: {
		source:  pdsSource{spacecraft: "spirit"},
		cameras: merCameras,
	},
}

// merCameras maps the Mars Exploration Rover cameras to the instrument names pdsSource
// reports. The archive calls both hazard cameras HAZCAM, so front and rear are told apart
// by the camera letter in the product file name.
var merCameras = map[string][]string{
	RoverCameraFHAZ.Slug:   {"FRONT_HAZCAM"},
	RoverCameraRHAZ.Slug:   {"REAR_HAZCAM"},
	RoverCameraNAVCAM.Slug: {"NAVCAM"},
	RoverCameraPANCAM.Slug: {"PANCAM"},
	RoverCameraMI.Slug:     {"MI"},
}

// merCameraLetters maps pdsSource instrument names to the camera letter that is the second
// character of MER product file names.
var merCameraLetters = map[string]string{
	"FRONT_HAZCAM": "f",
	"REAR_HAZCAM":  "r",
	"NAVCAM":       "n",
	"PANCAM":       "p",
	"MI":           "m",
}

// identityCameraMap maps each of the rover's cameras to an instrument of the same name.
func identityCameraMap(rover Rover) map[string][]string {
	m := make(map[string][]string, len(rover.Cameras))
	for _, c := range rover.Cameras {
		m[c.Slug] = []string{c.Name}
	}
	return m
}

func missionFor(rover Rover) (marsMission, error) {
	m, ok := marsMissions[rover.Slug]
	if !ok {
		return marsMission{}, ErrorRoverUnsupported
	}
	return m, nil
}

// instrumentsFor returns the feed instrument names for camera. An empty camera means no
// camera filter and returns nil.
func (m marsMission) instrumentsFor(rover Rover, camera RoverCamera) ([]string, error) {
	if camera.Slug == "" {
		return nil, nil
	}
	if !hasCamera(rover, camera) {
		return nil, &ErrorRoverCameraMissing{rover, camera}
	}
	return m.cameras[camera.Slug], nil
}

// cameraFor returns the RoverCamera that owns the feed instrument name.
func (m marsMission) cameraFor(rover Rover, instrument string) (RoverCamera, bool) {
	for _, c := range rover.Cameras {
		for _, name := range m.cameras[c.Slug] {
			if name == instrument {
				return c, true
			}
		}
	}
	return RoverCamera{}, false
}

// solsOn returns the sols that overlap the UTC calendar day of date.
func (m marsMission) solsOn(date time.Time) []int {
	day := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	first := int(day.Sub(m.solEpoch).Seconds() / secondsPerSol)
	last := int(day.Add(24*time.Hour-time.Nanosecond).Sub(m.solEpoch).Seconds() / secondsPerSol)

	sols := []int{}
	for sol := max(first, 0); sol <= last; sol++ {
		sols = append(sols, sol)
	}
	return sols
}

// curiositySource reads https://mars.nasa.gov/api/v1/raw_image_items/, which filters by sol
// and instrument but not by sample type, so full frames are picked out here.
type curiositySource struct{}

type curiosityResponse struct {
	Items []struct {
		ID         int       `json:"id"`
		ImageID    string    `json:"imageid"`
		Sol        int       `json:"sol"`
		Instrument string    `json:"instrument"`
		URL        string    `json:"https_url"`
		DateTaken  time.Time `json:"date_taken"`
		Extended   struct {
			SampleType string `json:"sample_type"`
		} `json:"extended"`
	} `json:"items"`
	More  bool `json:"more"`
	Total int  `json:"total"`
}

func (curiositySource) query(page, perPage int, conditions ...string) (curiosityResponse, error) {
	v := url.Values{}
	v.Set("order", "sol desc")
	v.Set("per_page", strconv.Itoa(perPage))
	v.Set("page", strconv.Itoa(page))
	v.Set("condition_1", "msl:mission")
	for i, c := range conditions {
		v.Set(fmt.Sprintf("condition_%d", i+2), c)
	}

	content, err := getContent(curiosityRawImagesURL+"?"+v.Encode(), nil)
	if err != nil {
		return curiosityResponse{}, err
	}

	r := curiosityResponse{}
	err = json.Unmarshal(content, &r)
	return r, err
}

func (s curiositySource) fetchSol(sol int, instruments []string) ([]sourcePhoto, error) {
	conditions := []string{fmt.Sprintf("%d:sol:in", sol)}
	if len(instruments) > 0 {
		conditions = append(conditions, strings.Join(instruments, ",")+":instrument:in")
	}

	photos := []sourcePhoto{}
	for page := 0; ; page++ {
		r, err := s.query(page, curiosityPerPage, conditions...)
		if err != nil {
			return nil, err
		}

		for _, i := range r.Items {
			if i.Extended.SampleType != "full" {
				continue
			}
			photos = append(photos, sourcePhoto{
				id:         i.ID,
				imageID:    i.ImageID,
				sol:        i.Sol,
				instrument: i.Instrument,
				url:        i.URL,
				takenAt:    i.DateTaken.UTC(),
			})
		}

		if !r.More || len(r.Items) == 0 {
			return photos, nil
		}
	}
}

func (s curiositySource) latestSol() (int, error) {
	r, err := s.query(0, 1)
	if err != nil {
		return 0, err
	}
	if len(r.Items) == 0 {
		return 0, nil
	}
	return r.Items[0].Sol, nil
}

func (s curiositySource) totalPhotos() (int, error) {
	r, err := s.query(0, 1)
	return r.Total, err
}

// perseveranceSource reads https://mars.nasa.gov/rss/api/?feed=raw_images&category=mars2020,
// which only returns full frames and filters by sol and instrument.
type perseveranceSource struct{}

type perseveranceResponse struct {
	Images []struct {
		ImageID    string `json:"imageid"`
		Sol        int    `json:"sol"`
		SampleType string `json:"sample_type"`
		Camera     struct {
			Instrument string `json:"instrument"`
		} `json:"camera"`
		ImageFiles struct {
			Large string `json:"large"`
		} `json:"image_files"`
		DateTakenUTC utcTime `json:"date_taken_utc"`
	} `json:"images"`
	TotalResults int `json:"total_results"`
	LatestSol    int `json:"latest_sol"`
}

func (perseveranceSource) query(v url.Values) (perseveranceResponse, error) {
	// The feed times out or returns 403 for some queries that omit an order.
	v.Set("order", "sol desc")
	v.Set("feed", "raw_images")
	v.Set("category", "mars2020")
	v.Set("feedtype", "json")
	v.Set("ver", "1.2")

	content, err := getContent(perseveranceRawImagesURL+"?"+v.Encode(), nil)
	if err != nil {
		return perseveranceResponse{}, err
	}

	r := perseveranceResponse{}
	err = json.Unmarshal(content, &r)
	return r, err
}

func (s perseveranceSource) fetchSol(sol int, instruments []string) ([]sourcePhoto, error) {
	photos := []sourcePhoto{}
	for page := 0; ; page++ {
		v := url.Values{}
		v.Set("num", strconv.Itoa(perseverancePerPage))
		v.Set("page", strconv.Itoa(page))
		v.Set("condition_2", fmt.Sprintf("%d:sol:gte", sol))
		v.Set("condition_3", fmt.Sprintf("%d:sol:lte", sol))
		if len(instruments) > 0 {
			v.Set("search", "|"+strings.Join(instruments, "|"))
		}

		r, err := s.query(v)
		if err != nil {
			return nil, err
		}

		for _, i := range r.Images {
			if !strings.EqualFold(i.SampleType, "full") {
				continue
			}
			photos = append(photos, sourcePhoto{
				imageID:    i.ImageID,
				sol:        i.Sol,
				instrument: i.Camera.Instrument,
				url:        i.ImageFiles.Large,
				takenAt:    i.DateTakenUTC.Time,
			})
		}

		if len(r.Images) < perseverancePerPage {
			return photos, nil
		}
	}
}

func (s perseveranceSource) latestSol() (int, error) {
	r, err := s.query(url.Values{"latest": {"true"}})
	return r.LatestSol, err
}

func (s perseveranceSource) totalPhotos() (int, error) {
	r, err := s.query(url.Values{"num": {"1"}, "page": {"0"}})
	return r.TotalResults, err
}

// utcTime parses timestamps that may or may not carry a zone; zoneless values are UTC.
type utcTime struct {
	time.Time
}

// UnmarshalJSON unmarshals an RFC 3339 timestamp, with or without a zone.
func (t *utcTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" || s == "null" {
		return nil
	}

	parsed, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		parsed, err = time.Parse("2006-01-02T15:04:05.999999999", s)
		if err != nil {
			return err
		}
	}
	*t = utcTime{Time: parsed.UTC()}
	return nil
}

// pdsSource reads the PDS Imaging Node Atlas search index for the Mars Exploration Rovers.
// It returns ILF products: full-frame images with only the onboard compression lookup table
// reversed. ILF is the one full-frame product archived for every camera, Pancam included.
type pdsSource struct {
	spacecraft string
}

type pdsResponse struct {
	Response struct {
		NumFound int `json:"numFound"`
		Docs     []struct {
			FileName  string  `json:"FILE_NAME"`
			ProductID string  `json:"PRODUCT_ID"`
			Sol       int     `json:"PLANET_DAY_NUMBER"`
			StartTime utcTime `json:"START_TIME"`
			BrowseURL string  `json:"ATLAS_BROWSE_URL"`
		} `json:"docs"`
	} `json:"response"`
}

// pdsILFFilter matches ILF products by their position in the file name, e.g.
// "1n137056144ilf2002p1825r0m1.img". Matching "ilf" anywhere would also catch mosaics. The
// camera letter comes second, so "?" is replaced by it to filter by camera.
const pdsILFFilter = "FILE_NAME:?%s?????????ilf*"

func (s pdsSource) query(filters []string, v url.Values) (pdsResponse, error) {
	v.Set("q", "*:*")
	v.Set("wt", "json")
	v.Set("fl", "FILE_NAME,PRODUCT_ID,PLANET_DAY_NUMBER,START_TIME,ATLAS_BROWSE_URL")
	v.Add("fq", `ATLAS_MISSION_NAME:"mars exploration rover"`)
	v.Add("fq", "ATLAS_SPACECRAFT_NAME:"+s.spacecraft)
	for _, f := range filters {
		v.Add("fq", f)
	}

	content, err := getContent(pdsImagingSearchURL+"?"+v.Encode(), nil)
	if err != nil {
		return pdsResponse{}, err
	}

	r := pdsResponse{}
	err = json.Unmarshal(content, &r)
	return r, err
}

// productFilter returns the ILF filter, limited to the given instruments when any are set.
func (pdsSource) productFilter(instruments []string) string {
	if len(instruments) == 0 {
		return fmt.Sprintf(pdsILFFilter, "?")
	}

	filters := make([]string, len(instruments))
	for i, instrument := range instruments {
		filters[i] = fmt.Sprintf(pdsILFFilter, merCameraLetters[instrument])
	}
	return strings.Join(filters, " OR ")
}

// fetch returns every ILF photo matching filters, following the index's pagination.
func (s pdsSource) fetch(instruments []string, filters ...string) ([]sourcePhoto, error) {
	filters = append(filters, s.productFilter(instruments))

	photos := []sourcePhoto{}
	for start := 0; ; start += pdsPerPage {
		v := url.Values{}
		v.Set("rows", strconv.Itoa(pdsPerPage))
		v.Set("start", strconv.Itoa(start))
		v.Set("sort", "START_TIME asc, PRODUCT_ID asc")

		r, err := s.query(filters, v)
		if err != nil {
			return nil, err
		}

		for _, d := range r.Response.Docs {
			photos = append(photos, sourcePhoto{
				imageID:    d.ProductID,
				sol:        d.Sol,
				instrument: merInstrument(d.FileName),
				url:        pdsFullSizeURL(d.BrowseURL),
				takenAt:    d.StartTime.Time,
			})
		}

		if len(r.Response.Docs) == 0 || start+pdsPerPage >= r.Response.NumFound {
			return photos, nil
		}
	}
}

func (s pdsSource) fetchSol(sol int, instruments []string) ([]sourcePhoto, error) {
	return s.fetch(instruments, fmt.Sprintf("PLANET_DAY_NUMBER:%d", sol))
}

func (s pdsSource) fetchDate(day time.Time, instruments []string) ([]sourcePhoto, error) {
	return s.fetch(instruments, fmt.Sprintf("START_TIME:[%s TO %s}",
		day.Format(time.RFC3339), day.AddDate(0, 0, 1).Format(time.RFC3339)))
}

func (s pdsSource) latestSol() (int, error) {
	v := url.Values{}
	v.Set("rows", "1")
	v.Set("sort", "PLANET_DAY_NUMBER desc")

	r, err := s.query([]string{s.productFilter(nil)}, v)
	if err != nil || len(r.Response.Docs) == 0 {
		return 0, err
	}
	return r.Response.Docs[0].Sol, nil
}

func (s pdsSource) totalPhotos() (int, error) {
	r, err := s.query([]string{s.productFilter(nil)}, url.Values{"rows": {"0"}})
	return r.Response.NumFound, err
}

// merInstrument returns the instrument name for a MER product file name, whose second
// character is the camera letter.
func merInstrument(fileName string) string {
	if len(fileName) > 1 {
		for instrument, letter := range merCameraLetters {
			if strings.EqualFold(fileName[1:2], letter) {
				return instrument
			}
		}
	}
	return ""
}

// pdsFullSizeURL turns an Atlas browse URL, which points at a 256x256 catalog JPEG, into the
// matching 1024x1024 JPEG in the volume's browse directory. The Atlas only indexes the
// catalog copy; both share the same path below the volume root.
func pdsFullSizeURL(browseURL string) string {
	u, err := url.Parse(browseURL)
	if err != nil || !strings.Contains(u.Path, "/extras/catalog/") {
		return browseURL
	}
	u.Path = strings.Replace(path.Clean(u.Path), "/extras/catalog/", "/browse/", 1)
	return u.String()
}
