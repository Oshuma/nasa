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
- [x] **Mars Rover Photos**: Image data gathered by NASA's Curiosity, Opportunity, Spirit, and Perseverance rovers on Mars
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

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Oshuma/nasa"
)

func main() {
	params := &nasa.APODParams{
		APIKey:    "YOUR_API_KEY",
		StartDate: time.Now().AddDate(0, 0, -3),
		EndDate:   time.Now(),
		Thumbs:    true,
	}

	items, err := nasa.APODList(params)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range items {
		fmt.Printf("%s: %s\n", item.Date.Format("2006-01-02"), item.Title)
		fmt.Printf("Media: %s\n\n", item.URL)
	}
}
```
