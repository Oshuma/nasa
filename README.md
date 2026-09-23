# NASA Open APIs

Go package for [NASA Open APIs](https://api.nasa.gov/).

### APIs Implemented

- [x] **APOD**: Astronomy Picture of the Day
- [ ] **Asteroids NeoWs**: Near Earth Object Web Service
- [ ] **DONKI**: Space Weather Database of Notifications, Knowledge, Information
- [ ] **Earth**: Unlock the significant public investment in earth observation data
- [ ] **EONET**: The Earth Observatory Natural Event Tracker
- [x] **EPIC**: Earth Polychromatic Imaging Camera
- [ ] **Exoplanet**: Programmatic access to NASA's Exoplanet Archive database
- [ ] **GeneLab**: Programmatic interface for GeneLab's public data repository website
- [ ] **Insight**: Mars Weather Service API
- [x] **Mars Rover Photos**: Image data gathered by NASA's Curiosity, Opportunity, Spirit, and Perseverance rovers on Mars, fetched live from NASA's raw image feeds and the PDS Imaging Atlas
- [x] **NASA Image and Video Library**: API to access the NASA Image and Video Library site at images.nasa.gov
- [ ] **TechTransfer**: Patents, Software, and Tech Transfer Reports
- [ ] **Satallite Situation Center**: System to cast geocentric spacecraft location information into a framework of (empirical) geophysical regions
- [ ] **SSD/CNEOS**: Solar System Dynamics and Center for Near-Earth Object Studies
- [ ] **Techport**: API to make NASA technology project data available in a machine-readable format
- [ ] **TLE API**: Two line element data for earth-orbiting objects at a given point in time
- [ ] **Vesta/Moon/Mars Trek WMTS**: A Web Map Tile Service for the Vesta, Moon, and Mars Trek imagery projects

## Using From Another Project

Install the module with Go modules.

```
go get github.com/Oshuma/nasa
```

Import the package and call the client helpers from your code.

## API Examples

All snippets assume a shared set of imports:

```go
import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Oshuma/nasa"
)
```

### APOD: Astronomy Picture of the Day

```go
apod, err := nasa.APOD(&nasa.APODParams{
	APIKey: os.Getenv("NASA_API_KEY"),
	Date:   time.Now(),
})
if err != nil {
	log.Fatal(err)
}

fmt.Printf("%s: %s\n", apod.Date.Format("2006-01-02"), apod.Title)
fmt.Printf("Media: %s\n", apod.URL)
```

If you need more than one image (for example, a date range) use `APODList`:

```go
items, err := nasa.APODList(&nasa.APODParams{
	APIKey:    os.Getenv("NASA_API_KEY"),
	StartDate: time.Now().AddDate(0, 0, -3),
	EndDate:   time.Now(),
	Thumbs:    true,
})
if err != nil {
	log.Fatal(err)
}

for _, item := range items {
	fmt.Printf("%s: %s (%s)\n", item.Date.Format("2006-01-02"), item.Title, item.MediaType)
}
```

### EPIC: Earth Polychromatic Imaging Camera

```go
images, err := nasa.EPIC(&nasa.EPICParams{
	APIKey:     os.Getenv("NASA_API_KEY"),
	Collection: "natural",
})
if err != nil {
	log.Fatal(err)
}

for _, img := range images {
	fmt.Printf("%s - %s\n", img.Date.Format("2006-01-02"), img.Caption)
	fmt.Println("Natural image:", img.URL.Natural)
	fmt.Println("Thumbnail:", img.URL.Thumb.Natural)
}
```

Passing a `Date` retrieves imagery for a specific day; otherwise the helper returns the most recent image for the requested collection (natural or enhanced).

### Mars Rover Photos

The original [Mars Rover Photos API](https://github.com/corincerami/mars-photo-api) was retired on October 8, 2025. These helpers now fetch full-frame images on the fly, with no API key and no caching:

- **Curiosity and Perseverance** come from NASA's raw image feeds, which update as new images arrive.
- **Opportunity and Spirit** come from the [PDS Imaging Atlas](https://pds-imaging.jpl.nasa.gov/search/), NASA's archive of mission data. These are `ILF` products: full-frame images from the Front and Rear Hazcams, Navcam, Pancam and Microscopic Imager (`nasa.RoverCameraMI`), returned as 1024×1024 JPEGs. Pancam captured each scene through several color filters, so a scene appears as several grayscale images.

> **Stability:** none of these sources is a documented public API. The raw image feeds (`mars.nasa.gov/api/v1/raw_image_items` for Curiosity and `mars.nasa.gov/rss/api` for Perseverance) power NASA's own raw image galleries, and the PDS search index (`pds-imaging.jpl.nasa.gov/solr/pds_archives/search`) powers the Imaging Atlas website. The full-size Opportunity and Spirit image URLs are derived from the archive's folder layout rather than returned by the index. Any of these may change or disappear without notice, which would break these helpers until the library is updated.

```go
photos, err := nasa.MarsRoverPhotos(&nasa.MarsPhotosParams{
	Sol:    1000,
	Camera: nasa.RoverCameraNAVCAM,
	Page:   1,
}, nasa.RoverCuriosity)
if err != nil {
	log.Fatal(err)
}

for _, photo := range photos.Photos {
	fmt.Printf("%s - %s\n", photo.EarthDate.Format("2006-01-02"), photo.Image)
}
```

Helpers such as `nasa.RoverCuriosity`, `nasa.RoverPerseverance`, and the predefined `RoverCamera*` constants ensure only valid camera selections are sent for a rover.

- Set `EarthDate` instead of `Sol` to get the photos captured on that UTC date.
- `Page` returns 25 photos per page; leave it at 0 to get every matching photo. `photos.Total` holds the number of matches across all pages.
- `MarsRoverPhotosLatest` returns the photos from the most recent sol that has full-frame images.
- `MarsMissionManifest` returns a partial manifest: mission dates and status, `MaxSol`, `MaxDate` and `TotalPhotos`, but no per-sol `Photos` breakdown.

### NASA Image and Video Library

```go
results, err := nasa.MediaSearch(&nasa.MediaParams{
	Query:     "apollo 11",
	MediaType: "image,video",
	Page:      1,
})
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Total hits: %d\n", results.Metadata.TotalHits)

for _, item := range results.Items {
	if len(item.Data) == 0 {
		continue
	}
	data := item.Data[0]
	fmt.Printf("%s - %s\n", data.NasaID, data.Title)
	for _, link := range item.Links {
		if link.Rel == "preview" {
			fmt.Println("Preview:", link.Href)
		}
	}
}
```

Any combination of documented search parameters can be supplied through `MediaParams`. Use pagination helpers like `Page` and `PageSize`, or refine searches with fields such as `Center`, `Keywords`, `Photographer`, and `YearStart` / `YearEnd`.
