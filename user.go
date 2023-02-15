package lastfm

import (
	"encoding/xml"
	"strconv"
)

type RecentTracks struct {
	User       string  `xml:"user,attr"`
	Total      int     `xml:"total,attr"`
	Tracks     []Track `xml:"track"`
	NowPlaying *Track  `xml:"-"` // Points to the currently playing track, if any
}

func (tracks *RecentTracks) unmarshalHelper() (err error) {
	for i, track := range tracks.Tracks {
		if track.NowPlaying {
			tracks.NowPlaying = &tracks.Tracks[i]
		}
		err = tracks.Tracks[i].unmarshalHelper()
		if err != nil {
			return
		}
	}
	return
}




// Gets a list of recent tracks from the user. The .Tracks field includes the currently playing track,
// if any, and up to the count most recent scrobbles.
// The .NowPlaying field points to any currently playing track.
//
// See http://www.last.fm/api/show/user.getRecentTracks.
func (lfm *LastFM) GetRecentTracks(user string, count int) (tracks *RecentTracks, err error) {
	method := "user.getRecentTracks"
	query := map[string]string{
		"user":     user,
		"extended": "1",
		"limit":    strconv.Itoa(count)}

	if data, err := lfm.cacheGet(method, query); data != nil {
		switch v := data.(type) {
		case RecentTracks:
			return &v, err
		case *RecentTracks:
			return v, err
		}
	} else if err != nil {
		return nil, err
	}

	body, hdr, err := lfm.doQuery(method, query)
	if err != nil {
		return
	}
	defer body.Close()

	status := lfmStatus{}
	err = xml.NewDecoder(body).Decode(&status)
	if err != nil {
		return
	}
	tracks = &status.RecentTracks
	err = tracks.unmarshalHelper() // Add this line
	if err == nil {
		go lfm.cacheSet(method, query, tracks, hdr)
	}
	return
	if status.Error.Code != 0 {
		err = &status.Error
		go lfm.cacheSet(method, query, err, hdr)
		return
	}
	return
}

type User struct {
	ID          int    `xml:"id"`
	Name        string `xml:"name"`
	RealName    string `xml:"realname"`
	URL         string `xml:"url"`
	Image       string `xml:"image"`
	Country     string `xml:"country"`
	Age         int    `xml:"age"`
	Gender      string `xml:"gender"`
	Subscriber  int    `xml:"subscriber"`
	Playcount   int    `xml:"playcount"`
	Playlists   int    `xml:"playlists"`
	Bootstrap   int    `xml:"bootstrap"`
	Registered  string `xml:"registered"`
	RegisterUNX string `xml:"registered,attr"`
}

// Returns user info
// 
//
// See http://www.last.fm/api/show/user.getInfo
func (lfm *LastFM) GetUserInfo(username string) (user *User, err error) {
	method := "user.getInfo"
	query := map[string]string{"user": username}

	if data, err := lfm.cacheGet(method, query); data != nil {
		switch v := data.(type) {
		case User:
			return &v, err
		case *User:
			return v, err
		}
	} else if err != nil {
		return nil, err
	}

	body, hdr, err := lfm.doQuery(method, query)
	if err != nil {
		return
	}
	defer body.Close()

	status := lfmStatus{}
	err = xml.NewDecoder(body).Decode(&status)
	if err != nil {
		return
	}
	if status.Error.Code != 0 {
		err = &status.Error
		go lfm.cacheSet(method, query, err, hdr)
		return
	}

	user = &status.User
	go lfm.cacheSet(method, query, user, hdr)
	return
}



type Tasteometer struct {
	Users   []string `xml:"input>user>name"`            // The compared users
	Score   float32  `xml:"result>score"`               // Varies from 0.0 to 1.0
	Artists []string `xml:"result>artists>artist>name"` // Short list of up to 5 common artists with the most affinity
}

// Compares the taste of 2 users.
//
// See http://www.last.fm/api/show/tasteometer.compare.
func (lfm *LastFM) CompareTaste(user1 string, user2 string) (taste *Tasteometer, err error) {
	method := "tasteometer.compare"
	query := map[string]string{
		"type1":  "user",
		"type2":  "user",
		"value1": user1,
		"value2": user2}

	if data, err := lfm.cacheGet(method, query); data != nil {
		switch v := data.(type) {
		case Tasteometer:
			return &v, err
		case *Tasteometer:
			return v, err
		}
	} else if err != nil {
		return nil, err
	}

	body, hdr, err := lfm.doQuery(method, query)
	if err != nil {
		return
	}
	defer body.Close()

	status := lfmStatus{}
	err = xml.NewDecoder(body).Decode(&status)
	if err != nil {
		return
	}
	if status.Error.Code != 0 {
		err = &status.Error
		go lfm.cacheSet(method, query, err, hdr)
		return
	}

	taste = &status.Tasteometer
	go lfm.cacheSet(method, query, taste, hdr)
	return
}

type Neighbour struct {
	Name  string  `xml:"name"`
	Match float32 `xml:"match"`
}
type Neighbours []Neighbour

// Gets a list of up to limit closest neighbours of a user. A neighbour is another user
// that has high tasteometer comparison scores.
//
// See http://www.last.fm/api/show/user.getNeighbours
func (lfm *LastFM) GetUserNeighbours(user string, limit int) (neighbours Neighbours, err error) {
	method := "user.getNeighbours"
	query := map[string]string{
		"user":  user,
		"limit": strconv.Itoa(limit)}

	if data, err := lfm.cacheGet(method, query); data != nil {
		return data.(Neighbours), err
	} else if err != nil {
		return nil, err
	}

	body, hdr, err := lfm.doQuery(method, query)
	if err != nil {
		return
	}
	defer body.Close()

	status := lfmStatus{}
	err = xml.NewDecoder(body).Decode(&status)
	if err != nil {
		return
	}
	if status.Error.Code != 0 {
		err = &status.Error
		go lfm.cacheSet(method, query, err, hdr)
		return
	}

	neighbours = status.Neighbours
	go lfm.cacheSet(method, query, neighbours, hdr)
	return
}

type Period int

const (
	Overall Period = 1 + iota
	OneWeek
	OneMonth
	ThreeMonths
	SixMonths
	OneYear
)

var periodStringMap = map[Period]string{
	Overall:     "overall",
	OneWeek:     "7day",
	OneMonth:    "1month",
	ThreeMonths: "3month",
	SixMonths:   "6month",
	OneYear:     "12month"}

func (p Period) String() string {
	return periodStringMap[p]
}

type TopArtists struct {
	User   string `xml:"user,attr"`
	Period Period `xml:"-"`
	Total  int    `xml:"total,attr"`

	Artists []Artist `xml:"artist"`

	// For internal use
	RawPeriod string `xml:"type,attr"`
}

func (top *TopArtists) unmarshalHelper() (err error) {
	for k, v := range periodStringMap {
		if top.RawPeriod == v {
			top.Period = k
			break
		}
	}
	return
}

// Gets a list of the (up to limit) most played artists of a user within a Period.
//
// See http://www.last.fm/api/show/user.getTopArtists.
func (lfm *LastFM) GetUserTopArtists(user string, period Period, limit int) (top *TopArtists, err error) {
	method := "user.getTopArtists"
	query := map[string]string{
		"user":   user,
		"period": periodStringMap[period],
		"limit":  strconv.Itoa(limit)}

	if data, err := lfm.cacheGet(method, query); data != nil {
		switch v := data.(type) {
		case TopArtists:
			return &v, err
		case *TopArtists:
			return v, err
		}
	} else if err != nil {
		return nil, err
	}

	body, hdr, err := lfm.doQuery(method, query)
	if err != nil {
		return
	}
	defer body.Close()

	status := lfmStatus{}
	err = xml.NewDecoder(body).Decode(&status)
	if err != nil {
		return
	}
	if status.Error.Code != 0 {
		err = &status.Error
		go lfm.cacheSet(method, query, err, hdr)
		return
	}

	top = &status.TopArtists
	err = top.unmarshalHelper()
	if err == nil {
		go lfm.cacheSet(method, query, top, hdr)
	}
	return
}

type ArtistInfo struct {
	Name string `xml:"name"`
	MBID string `xml:"mbid"`
	URL string `xml:"url"`
	Image struct {
	  SizeSmall string `xml:"small"`
	  SizeMedium string `xml:"medium"`
	  SizeLarge string `xml:"large"`
	  SizeExtraLarge string `xml:"extralarge"`
	  SizeMega string `xml:"mega"`
	  Size string `xml:""`
	} `xml:"image"`
	Streamable int `xml:"streamable"`
	Ontour int `xml:"ontour"`
	Stats struct {
	  Listeners int `xml:"listeners"`
	  Playcount int `xml:"playcount"`
	  Userplaycount int `xml:"userplaycount"`
	} `xml:"stats"`
}

func (info *ArtistInfo) unmarshalHelper() (err error) {
	info.Duration, err = time.ParseDuration(info.RawDuration + "ms")
	if err != nil {
		return
	}
	if info.Wiki != nil {
		err = info.Wiki.unmarshalHelper()
	}
	return
}

// Gets information for a Artist. The user argument can either be empty ("") or specify a last.fm username, in which
// case .UserPlaycount will be valid in the returned struct. The autocorrect parameter controls whether
// last.fm's autocorrection algorithms should be run on the artist name.
//
// The Artist struct must specify either the MBID or Artist.Name.
// Example literals that can be given as the first argument:
//   lastfm.Artist{MBID: "mbid"}
//   lastfm.Artist{Name: "Artist"}
//
// See http://www.last.fm/api/show/artist.getInfo.
func (lfm *LastFM) GetArtistInfo(artist Artist, user string, autocorrect bool) (info *ArtistInfo, err error) {
	method := "artist.getInfo"
	query := map[string]string{}
	if autocorrect {
		query["autocorrect"] = "1"
	} else {
		query["autocorrect"] = "0"
	}

	if user != "" {
		query["username"] = user
	}

	if artist.MBID != "" {
		query["mbid"] = artist.MBID
	} else {
		query["artist"] = artist.Name
	}

	if data, err := lfm.cacheGet(method, query); data != nil {
		switch v := data.(type) {
		case ArtistInfo:
			return &v, err
		case *ArtistInfo:
			return v, err
		}
	} else if err != nil {
		return nil, err
	}

	body, hdr, err := lfm.doQuery(method, query)
	if err != nil {
		return
	}
	defer body.Close()

	status := lfmStatus{}
	err = xml.NewDecoder(body).Decode(&status)
	if err != nil {
		return
	}
	if status.Error.Code != 0 {
		err = &status.Error
		go lfm.cacheSet(method, query, err, hdr)
		return
	}

	info = &status.ArtistInfo
	err = info.unmarshalHelper()
	if err == nil {
		go lfm.cacheSet(method, query, info, hdr)
	}
	return
}
