package main

//Node struct
type Node struct {
	id        int64
	latitude  float64
	longitude float64
	version   int64
	timestamp int64
	changeset int64
	uid       int32
	sid       string
	tags      map[string]string
}
