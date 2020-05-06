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

// Media represents a media search response.
type Media struct {
	Metadata struct {
		TotalHits int `json:"total_hits"`
	} `json:"metadata"`

	Items []struct {
		Data []struct {
			Center           string    `json:"center"`
			SecondaryCreator string    `json:"secondary_creator"`
			Keywords         []string  `json:"keywords"`
			Description      string    `json:"description"`
			Description508   string    `json:"description_508"`
			MediaType        string    `json:"media_type"`
			NasaID           string    `json:"nasa_id"`
			DateCreated      time.Time `json:"date_created"`
		} `json:"data"`

		Links []struct {
			Render string `json:"render"`
			Href   string `json:"href"`
			Rel    string `json:"rel"`
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

// MediaSearch searches the NASA Image and Video Library.
func MediaSearch(p ParamEncoder) (Media, error) {
	content, err := getContent(mediaAPIURL, p)
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

// MediaAssets represents a media search query.
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

// GetMediaAssets gets the media assets for the given nasaID.
func GetMediaAssets(nasaID string) (MediaAssets, error) {
	url := fmt.Sprintf(assetAPIURL, nasaID)
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

// MediaMetadata holds metadata info for a media resource.
type MediaMetadata struct {
	// TODO: Split these off into structs based on their namespace.
	// TODO: Parse the weird timestamps: "2006:01:02 15:04:05" and "2006:01:02 15:04:05-0700"
	AVAILAlbum                   string   `json:"AVAIL:Album"`
	AVAILCenter                  string   `json:"AVAIL:Center"`
	AVAILDateCreated             string   `json:"AVAIL:DateCreated"`
	AVAILDescription             string   `json:"AVAIL:Description"`
	AVAILDescription508          string   `json:"AVAIL:Description508"`
	AVAILKeywords                []string `json:"AVAIL:Keywords"`
	AVAILLocation                string   `json:"AVAIL:Location"`
	AVAILMediaType               string   `json:"AVAIL:MediaType"`
	AVAILNASAID                  string   `json:"AVAIL:NASAID"`
	AVAILOwner                   string   `json:"AVAIL:Owner"`
	AVAILPhotographer            string   `json:"AVAIL:Photographer"`
	AVAILSecondaryCreator        string   `json:"AVAIL:SecondaryCreator"`
	AVAILTitle                   string   `json:"AVAIL:Title"`
	CompositeImageSize           string   `json:"Composite:ImageSize"`
	CompositeMegapixels          float64  `json:"Composite:Megapixels"`
	EXIFColorSpace               string   `json:"EXIF:ColorSpace"`
	EXIFComponentsConfiguration  string   `json:"EXIF:ComponentsConfiguration"`
	EXIFCreateDate               string   `json:"EXIF:CreateDate"`
	EXIFExifVersion              string   `json:"EXIF:ExifVersion"`
	EXIFFlashpixVersion          string   `json:"EXIF:FlashpixVersion"`
	EXIFImageDescription         string   `json:"EXIF:ImageDescription"`
	EXIFResolutionUnit           string   `json:"EXIF:ResolutionUnit"`
	EXIFXResolution              int      `json:"EXIF:XResolution"`
	EXIFYCbCrPositioning         string   `json:"EXIF:YCbCrPositioning"`
	EXIFYResolution              int      `json:"EXIF:YResolution"`
	ExifToolExifToolVersion      float64  `json:"ExifTool:ExifToolVersion"`
	FileBitsPerSample            int      `json:"File:BitsPerSample"`
	FileColorComponents          int      `json:"File:ColorComponents"`
	FileCurrentIPTCDigest        string   `json:"File:CurrentIPTCDigest"`
	FileDirectory                string   `json:"File:Directory"`
	FileEncodingProcess          string   `json:"File:EncodingProcess"`
	FileExifByteOrder            string   `json:"File:ExifByteOrder"`
	FileFileAccessDate           string   `json:"File:FileAccessDate"`
	FileFileInodeChangeDate      string   `json:"File:FileInodeChangeDate"`
	FileFileModifyDate           string   `json:"File:FileModifyDate"`
	FileFileName                 string   `json:"File:FileName"`
	FileFilePermissions          string   `json:"File:FilePermissions"`
	FileFileSize                 string   `json:"File:FileSize"`
	FileFileType                 string   `json:"File:FileType"`
	FileFileTypeExtension        string   `json:"File:FileTypeExtension"`
	FileImageHeight              int      `json:"File:ImageHeight"`
	FileImageWidth               int      `json:"File:ImageWidth"`
	FileMIMEType                 string   `json:"File:MIMEType"`
	FileYCbCrSubSampling         string   `json:"File:YCbCrSubSampling"`
	IPTCApplicationRecordVersion int      `json:"IPTC:ApplicationRecordVersion"`
	IPTCKeywords                 []string `json:"IPTC:Keywords"`
	JFIFJFIFVersion              float64  `json:"JFIF:JFIFVersion"`
	JFIFResolutionUnit           string   `json:"JFIF:ResolutionUnit"`
	JFIFXResolution              int      `json:"JFIF:XResolution"`
	JFIFYResolution              int      `json:"JFIF:YResolution"`
	SourceFile                   string   `json:"SourceFile"`
	XMPCreateDate                string   `json:"XMP:CreateDate"`
	XMPCreatedate                string   `json:"XMP:Createdate"`
	XMPCredit                    string   `json:"XMP:Credit"`
	XMPDateCreated               string   `json:"XMP:DateCreated"`
	XMPDescription               string   `json:"XMP:Description"`
	XMPImageDescription          string   `json:"XMP:ImageDescription"`
	XMPNasaID                    string   `json:"XMP:Nasa_id"`
	XMPSource                    string   `json:"XMP:Source"`
	XMPTitle                     string   `json:"XMP:Title"`
	XMPXMPToolkit                string   `json:"XMP:XMPToolkit"`
}

// GetMediaMetadata gets the metadata for media with nasaID.
func GetMediaMetadata(nasaID string) (MediaMetadata, error) {
	url := fmt.Sprintf(metadataAPIURL, nasaID)
	content, err := getContent(url, nil)
	if err != nil {
		return MediaMetadata{}, err
	}

	resp := mediaMetadataResponse{}
	err = json.Unmarshal(content, &resp)
	if err != nil {
		return MediaMetadata{}, err
	}

	// Make sure we get a valid URL.
	_, err = neturl.Parse(resp.Location)
	if err != nil {
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
	url := fmt.Sprintf(captionsAPIURL, nasaID)
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

	captions, err := getContent(c.Location, nil)
	if err != nil {
		return "", err
	}

	return string(captions), nil
}
