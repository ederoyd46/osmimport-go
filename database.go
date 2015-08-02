package main

import (
	"fmt"

	r "github.com/dancannon/gorethink"
)

var (
	session *r.Session
)

func init() {
	session = connect()
	checkDatabase()
}

func checkDatabase() {
	databases, err := r.DBList().Run(session)
	LogError(err)
	fmt.Println(databases)
}

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
