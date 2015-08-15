package main

//Node struct
type Node struct {
	ID        int64             `gorethink:"osm_id,omitempty"`
	Latitude  float64           `gorethink:"latitude,omitempty"`
	Longitude float64           `gorethink:"longitude,omitempty"`
	Version   int64             `gorethink:"version,omitempty"`
	Timestamp string            `gorethink:"timestamp,omitempty"`
	Changeset int64             `gorethink:"changeset,omitempty"`
	UID       int32             `gorethink:"uid,omitempty"`
	User      string            `gorethink:"user,omitempty"`
	Tags      map[string]string `gorethink:"tags,omitempty"`
}
