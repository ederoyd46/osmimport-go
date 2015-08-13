package main

import r "github.com/dancannon/gorethink"

var (
	session      *r.Session
	databaseName string
	nodeTable    = "node"
)

func checkDatabase() {
	databases, err := r.DBList().Run(session)
	var result string
	for databases.Next(&result) {
		if result == databaseName {
			return
		}
	}
	_, err = r.DBCreate(databaseName).Run(session)
	LogError(err)
}

func checkTables() {
	tables, err := r.DB(databaseName).TableList().Run(session)
	var result string
	for tables.Next(&result) {
		if result == nodeTable {
			return
		}
	}

	_, err = r.DB(databaseName).TableCreate(nodeTable).Run(session)
	LogError(err)
}

func connect(host string) *r.Session {
	var session *r.Session
	session, err := r.Connect(r.ConnectOpts{
		Address: host,
	})
	LogError(err)
	return session
}

//InitDB Sets up the database connection pool
func InitDB(host, dbname string) {
	session = connect(host)
	databaseName = dbname
	checkDatabase()
	checkTables()
}

//KillSession disconnects from the database
func KillSession() {
	session.Close()
}

//SaveNodes saves a node to the database
func SaveNodes(node []Node) {
	_, err := r.DB(databaseName).Table(nodeTable).Insert(node).RunWrite(session)
	LogError(err)
}
