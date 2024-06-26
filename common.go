package lastfm

import (
	"strings"
	"time"
)

// Some structs need extra processing after XML unmarshalling.
type unmarshalHelper interface {
	unmarshalHelper() error
}

type LastFMError struct {
	error
	Code    int    `xml:"code,attr"`
	Message string `xml:",chardata"`
}

func (e *LastFMError) Error() string {
	return strings.Trim(e.Message, "\n ")
}

// Known image sizes
const (
 	SmallImageSize      = "small"
 	MediumImageSize     = "medium"
 	LargeImageSize      = "large"
 	ExtraLargeImageSize = "extralarge"
)


type Image struct {
	Size  string `xml:"size,attr"`
	URL   string `xml:",chardata"`
}

type Artist struct {
 	Name      string   `xml:"name"`
 	PlayCount int      `xml:"playcount"` // Currently is always 0, except when part of the result of GetUserTopArtists.
 	MBID      string   `xml:"mbid"`
 	URL       string   `xml:"url"`
 	Images    []*Image `xml:"image"`
}

// Less detailed struct returned in GetRecentTracks.

type Album struct {
	Artist    string `xml:"artist"`
	MBID      string `xml:"mbid"`
	Name      string `xml:"name"`
	Streamable int    `xml:"streamable"`
	Album     string `xml:"album"`
	URL       string `xml:"url"`
	Images    []*Image `xml:"image"`
	Date      string `xml:"date"`
}

// More detailed struct returned in GetTrackInfo.
type AlbumInfo struct {
 	TrackNo int      `xml:"position,attr"`
 	Name    string   `xml:"title"`
 	Artist  string   `xml:"artist"`
 	MBID    string   `xml:"mbid"`
 	URL     string   `xml:"url"`
 	Images  []*Image `xml:"image"`
 }

 
type Track struct {
	NowPlaying bool      `xml:"nowplaying,attr"`
	Images     []*Image  `xml:"image"`
	Artist     Artist    `xml:"artist"`
	Album      Album     `xml:"album"`
	Loved      bool      `xml:"loved"`
	Name       string    `xml:"name"`
	MBID       string    `xml:"mbid"`
	URL        string    `xml:"url"`
	Date       time.Time `xml:"-"`

	// For internal use
	RawDate lfmDate `xml:"date"`
}



func (track *Track) unmarshalHelper() (err error) {
	if track.RawDate.Date != "" {
		track.Date = time.Unix(track.RawDate.UTS, 0)
	}
	return
}

type Wiki struct {
	Published time.Time `xml:"-"`
	Summary   string    `xml:"summary"`
	Content   string    `xml:"content"`

	// For internal use
	RawPublished string `xml:"published"`
}

func (wiki *Wiki) unmarshalHelper() (err error) {
	if wiki.RawPublished != "" {
		wiki.Published, err = time.Parse("02 Jan 2006, 15:04", wiki.RawPublished)
	}
	return
}
