# Go-LastFM

This package was forked and updated to support (most) of the current Last.FM routes. Feel free to open a pull requests for any improvments. 
This package should be rewritten soon, if I feel like it.


# How To Install
```
go get github.com/lxi1400/go-lastfm
```
# Example
## Print Current/Last Track Name
```golang
package main

import (
	"fmt"
	"github.com/lxi1400/go-lastfm"
)

func main() {
	// Get Tracks
	tracks, err := lfm.GetRecentTracks("username", 1) 
	if err != nil {
		fmt.Println(err.Error())
	}

	// Get current track
	trackSlice := tracks.Tracks
	currentTrack := trackSlice[0]

	// Print the current tracks name and if it is currently playing
	fmt.Println(currentTrack["nam"])
	fmt.Println(tracks.NowPlaying)
}

var (
	lfm = lastfm.New("API_KEY")
)

```
