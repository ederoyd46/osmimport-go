package main

import "time"

//Way struct
type Way struct {
	ID        int64             `gorethink:"osm_id,omitempty"`
	Version   int32             `gorethink:"version,omitempty"`
	Timestamp time.Time         `gorethink:"timestamp,omitempty"`
	Changeset int64             `gorethink:"changeset,omitempty"`
	UID       int32             `gorethink:"uid,omitempty"`
	User      string            `gorethink:"user,omitempty"`
	Tags      map[string]string `gorethink:"tags,omitempty"`
	Refs      []int64           `gorethink:"refs,omitempty"`
}
