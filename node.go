package main

//Node struct
type Node struct {
	ID        int64             `gorethink:"osm_id"`
	Latitude  float64           `gorethink:"latitude"`
	Longitude float64           `gorethink:"longitude"`
	Version   int64             `gorethink:"version"`
	Timestamp int64             `gorethink:"timestamp"`
	Changeset int64             `gorethink:"changeset"`
	UID       int32             `gorethink:"uid"`
	SID       string            `gorethink:"sid"`
	Tags      map[string]string `gorethink:"tags"`
}
