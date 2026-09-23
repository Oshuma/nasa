package nasa

import (
	"encoding/json"
	"errors"
	"fmt"
	neturl "net/url"
	"time"
)

const (
	mediaAPIURL    = "https://images-api.nasa.gov/search"
	assetAPIURL    = "https://images-api.nasa.gov/asset/%s"
	metadataAPIURL = "https://images-api.nasa.gov/metadata/%s"
	captionsAPIURL = "https://images-api.nasa.gov/captions/%s"
	albumAPIURL    = "https://images-api.nasa.gov/album/%s"
)

// Media represents a NASA Image and Video Library search response. The structure maps
// directly to the Collection+JSON payload returned by the /search endpoint.
type Media struct {
	Metadata struct {
		TotalHits int `json:"total_hits"`
	} `json:"metadata"`

	Items []struct {
		Data []struct {
			Center           string    `json:"center"`
			Album            []string  `json:"album"`
			SecondaryCreator string    `json:"secondary_creator"`
			Keywords         []string  `json:"keywords"`
			Description      string    `json:"description"`
			Description508   string    `json:"description_508"`
			MediaType        string    `json:"media_type"`
			NasaID           string    `json:"nasa_id"`
			DateCreated      time.Time `json:"date_created"`
			Location         string    `json:"location"`
			Photographer     string    `json:"photographer"`
			Title            string    `json:"title"`
		} `json:"data"`

		Links []struct {
			Render string `json:"render"`
			Href   string `json:"href"`
			Rel    string `json:"rel"`
			Prompt string `json:"prompt"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
			Size   int64  `json:"size"`
		} `json:"links"`

		Href string `json:"href"`
	} `json:"items"`

	Links []struct {
		Prompt string `json:"prompt"`
		Href   string `json:"href"`
		Rel    string `json:"rel"`
	} `json:"links"`

	Version string `json:"version"`
	Href    string `json:"href"`
}

type mediaResponse struct {
	Collection Media `json:"collection"`
}

// MediaSearch queries the NASA Image and Video Library /search endpoint using the
// provided MediaParams. At least one search parameter must be supplied; the helper
// unmarshals the Collection+JSON response into a Media value for easier consumption.
func MediaSearch(p ParamEncoder) (Media, error) {
	params, ok := p.(*MediaParams)
	if !ok {
		return Media{}, ErrorParamsMismatch
	}

	content, err := getContent(mediaAPIURL, params)
	if err != nil {
		return Media{}, err
	}

	mr := mediaResponse{}
	err = json.Unmarshal(content, &mr)
	if err != nil {
		return Media{}, err
	}

	return mr.Collection, nil
}

// MediaAssets represents the results from the media asset manifest endpoint,
// describing downloadable resources for a media item identified by NASA ID.
type MediaAssets struct {
	Items []struct {
		Href string `json:"href"`
	} `json:"items"`
	Version string `json:"version"`
	Href    string `json:"href"`
}

type mediaAssetResponse struct {
	Collection MediaAssets `json:"collection"`
}

// GetMediaAssets retrieves the media asset manifest for the given NASA ID, returning
// download links for each available rendition (image, video, metadata, etc.).
func GetMediaAssets(nasaID string) (MediaAssets, error) {
	url := fmt.Sprintf(assetAPIURL, neturl.PathEscape(nasaID))
	content, err := getContent(url, nil)
	if err != nil {
		return MediaAssets{}, err
	}

	mar := mediaAssetResponse{}
	err = json.Unmarshal(content, &mar)
	if err != nil {
		return MediaAssets{}, err
	}

	return mar.Collection, nil
}

type mediaMetadataResponse struct {
	Location string `json:"location"`
}

// StringList is a []string that also accepts a single JSON string, which some metadata
// fields use when they hold only one value.
type StringList []string

// UnmarshalJSON unmarshals either a JSON string or an array of strings.
func (l *StringList) UnmarshalJSON(b []byte) error {
	var one string
	if err := json.Unmarshal(b, &one); err == nil {
		if one == "" {
			*l = nil
		} else {
			*l = StringList{one}
		}
		return nil
	}

	var many []string
	if err := json.Unmarshal(b, &many); err != nil {
		return err
	}
	*l = many
	return nil
}

// MediaMetadata holds metadata info for a media resource. Use GetMediaMetadata to fetch
// the metadata JSON and unmarshal it into this structure.
type MediaMetadata struct {
	// TODO: Split these off into structs based on their namespace.
	// TODO: Parse the weird timestamps: "2006:01:02 15:04:05" and "2006:01:02 15:04:05-0700"
	AVAILAlbum                   StringList `json:"AVAIL:Album"`
	AVAILCenter                  string     `json:"AVAIL:Center"`
	AVAILDateCreated             string     `json:"AVAIL:DateCreated"`
	AVAILDescription             string     `json:"AVAIL:Description"`
	AVAILDescription508          string     `json:"AVAIL:Description508"`
	AVAILKeywords                []string   `json:"AVAIL:Keywords"`
	AVAILLocation                string     `json:"AVAIL:Location"`
	AVAILMediaType               string     `json:"AVAIL:MediaType"`
	AVAILNASAID                  string     `json:"AVAIL:NASAID"`
	AVAILOwner                   string     `json:"AVAIL:Owner"`
	AVAILPhotographer            string     `json:"AVAIL:Photographer"`
	AVAILSecondaryCreator        string     `json:"AVAIL:SecondaryCreator"`
	AVAILTitle                   string     `json:"AVAIL:Title"`
	CompositeImageSize           string     `json:"Composite:ImageSize"`
	CompositeMegapixels          float64    `json:"Composite:Megapixels"`
	EXIFColorSpace               string     `json:"EXIF:ColorSpace"`
	EXIFComponentsConfiguration  string     `json:"EXIF:ComponentsConfiguration"`
	EXIFCreateDate               string     `json:"EXIF:CreateDate"`
	EXIFExifVersion              string     `json:"EXIF:ExifVersion"`
	EXIFFlashpixVersion          string     `json:"EXIF:FlashpixVersion"`
	EXIFImageDescription         string     `json:"EXIF:ImageDescription"`
	EXIFResolutionUnit           string     `json:"EXIF:ResolutionUnit"`
	EXIFXResolution              int        `json:"EXIF:XResolution"`
	EXIFYCbCrPositioning         string     `json:"EXIF:YCbCrPositioning"`
	EXIFYResolution              int        `json:"EXIF:YResolution"`
	ExifToolExifToolVersion      float64    `json:"ExifTool:ExifToolVersion"`
	FileBitsPerSample            int        `json:"File:BitsPerSample"`
	FileColorComponents          int        `json:"File:ColorComponents"`
	FileCurrentIPTCDigest        string     `json:"File:CurrentIPTCDigest"`
	FileDirectory                string     `json:"File:Directory"`
	FileEncodingProcess          string     `json:"File:EncodingProcess"`
	FileExifByteOrder            string     `json:"File:ExifByteOrder"`
	FileFileAccessDate           string     `json:"File:FileAccessDate"`
	FileFileInodeChangeDate      string     `json:"File:FileInodeChangeDate"`
	FileFileModifyDate           string     `json:"File:FileModifyDate"`
	FileFileName                 string     `json:"File:FileName"`
	FileFilePermissions          string     `json:"File:FilePermissions"`
	FileFileSize                 string     `json:"File:FileSize"`
	FileFileType                 string     `json:"File:FileType"`
	FileFileTypeExtension        string     `json:"File:FileTypeExtension"`
	FileImageHeight              int        `json:"File:ImageHeight"`
	FileImageWidth               int        `json:"File:ImageWidth"`
	FileMIMEType                 string     `json:"File:MIMEType"`
	FileYCbCrSubSampling         string     `json:"File:YCbCrSubSampling"`
	IPTCApplicationRecordVersion int        `json:"IPTC:ApplicationRecordVersion"`
	IPTCKeywords                 StringList `json:"IPTC:Keywords"`
	JFIFJFIFVersion              float64    `json:"JFIF:JFIFVersion"`
	JFIFResolutionUnit           string     `json:"JFIF:ResolutionUnit"`
	JFIFXResolution              int        `json:"JFIF:XResolution"`
	JFIFYResolution              int        `json:"JFIF:YResolution"`
	SourceFile                   string     `json:"SourceFile"`
	XMPCreateDate                string     `json:"XMP:CreateDate"`
	XMPCreatedate                string     `json:"XMP:Createdate"`
	XMPCredit                    string     `json:"XMP:Credit"`
	XMPDateCreated               string     `json:"XMP:DateCreated"`
	XMPDescription               string     `json:"XMP:Description"`
	XMPImageDescription          string     `json:"XMP:ImageDescription"`
	XMPNasaID                    string     `json:"XMP:Nasa_id"`
	XMPSource                    string     `json:"XMP:Source"`
	XMPTitle                     string     `json:"XMP:Title"`
	XMPXMPToolkit                string     `json:"XMP:XMPToolkit"`
}

// GetMediaMetadata gets the metadata for media with the provided NASA ID. The helper
// first resolves the metadata location and then retrieves the JSON payload.
func GetMediaMetadata(nasaID string) (MediaMetadata, error) {
	url := fmt.Sprintf(metadataAPIURL, neturl.PathEscape(nasaID))
	content, err := getContent(url, nil)
	if err != nil {
		return MediaMetadata{}, err
	}

	resp := mediaMetadataResponse{}
	err = json.Unmarshal(content, &resp)
	if err != nil {
		return MediaMetadata{}, err
	}

	// Make sure we get a valid, absolute URL.
	u, err := neturl.Parse(resp.Location)
	if err != nil || !u.IsAbs() || u.Host == "" {
		return MediaMetadata{}, ErrorNoMetadata
	}

	content, err = getContent(resp.Location, nil)
	if err != nil {
		return MediaMetadata{}, err
	}

	metadata := MediaMetadata{}
	err = json.Unmarshal(content, &metadata)
	if err != nil {
		return MediaMetadata{}, err
	}

	return metadata, nil
}

// GetMediaCaptions returns the captions for the given nasaID.
// TODO: Maybe parse the captions in later versions.
func GetMediaCaptions(nasaID string) (string, error) {
	url := fmt.Sprintf(captionsAPIURL, neturl.PathEscape(nasaID))
	content, err := getContent(url, nil)
	if err != nil {
		return "", err
	}

	type errorResponse struct {
		Reason string `json:"reason"`
	}
	r := errorResponse{}
	err = json.Unmarshal(content, &r)
	if err != nil {
		return "", err
	}
	if r.Reason != "" {
		return "", errors.New(r.Reason)
	}

	type captionLocation struct {
		Location string `json:"location"`
	}
	c := captionLocation{}
	err = json.Unmarshal(content, &c)
	if err != nil {
		return "", err
	}
	if c.Location == "" {
		return "", ErrorNoCaptions
	}

	captions, err := getContent(c.Location, nil)
	if err != nil {
		return "", err
	}

	return string(captions), nil
}
