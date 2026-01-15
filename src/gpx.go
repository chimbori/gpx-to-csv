package main

import (
	"bytes"
	"encoding/xml"
)

type GPX struct {
	Version   string      `xml:"version,attr"`
	Tracks    []*GPXTrack `xml:"trk"`
	Waypoints []*GPXPoint `xml:"wpt"`
}

type GPXTrack struct {
	Segments []*GPXTrackSegment `xml:"trkseg,omitempty"`
}

type GPXTrackSegment struct {
	Points []*GPXPoint `xml:"trkpt"`
}

type GPXPoint struct {
	Latitude  float64 `xml:"lat,attr"`
	Longitude float64 `xml:"lon,attr"`
	Timestamp string  `xml:"time,omitempty"`
}

func parseGpx(buf []byte) (*GPX, error) {
	decoder := xml.NewDecoder(bytes.NewReader(buf))
	gpx := &GPX{}
	err := decoder.Decode(&gpx)
	if err != nil {
		return nil, err
	}
	return gpx, nil
}
