package main

//Relation struct
type Relation struct {
	ID        int64             `gorethink:"osm_id,omitempty"`
	Version   int32             `gorethink:"version,omitempty"`
	Timestamp string            `gorethink:"timestamp,omitempty"`
	Changeset int64             `gorethink:"changeset,omitempty"`
	UID       int32             `gorethink:"uid,omitempty"`
	User      string            `gorethink:"user,omitempty"`
	Tags      map[string]string `gorethink:"tags,omitempty"`
	MemIds    []int64           `gorethink:"memids,omitempty"`
	Roles     []string          `gorethink:"roles,omitempty"`
	Types     []string          `gorethink:"types,omitempty"`
}
