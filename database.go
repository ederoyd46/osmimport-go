package main

import r "github.com/dancannon/gorethink"

func connect() *r.Session {
	var session *r.Session
	session, err := r.Connect(r.ConnectOpts{
		Address: "localhost:28015",
	})
	LogError(err)
	return session
}

func disconnect(session *r.Session) {
	session.Close()
}

//SaveNodes saves a node to the database
func SaveNodes(node []Node) {
	session := connect()
	_, err := r.DB("test").Table("node").Insert(node).RunWrite(session)
	disconnect(session)
	LogError(err)
}
