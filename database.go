package main

import r "github.com/dancannon/gorethink"

var (
	session       *r.Session
	databaseName  string
	nodeTable     = "node"
	wayTable      = "way"
	relationTable = "relation"
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

func checkTable(tname string) {
	tables, err := r.DB(databaseName).TableList().Run(session)
	var result string
	for tables.Next(&result) {
		if result == tname {
			return
		}
	}
	_, err = r.DB(databaseName).TableCreate(tname).Run(session)
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
	checkTable(nodeTable)
	checkTable(wayTable)
	checkTable(relationTable)
}

//KillSession disconnects from the database
func KillSession() {
	session.Close()
}

//SaveNodes saves nodes to the database
func SaveNodes(nodes []Node) {
	_, err := r.DB(databaseName).Table(nodeTable).Insert(nodes).RunWrite(session)
	LogError(err)
}

//SaveWays saves a node to the database
func SaveWays(ways []Way) {
	_, err := r.DB(databaseName).Table(wayTable).Insert(ways).RunWrite(session)
	LogError(err)
}

//SaveRelations saves a node to the database
func SaveRelations(relations []Relation) {
	_, err := r.DB(databaseName).Table(relationTable).Insert(relations).RunWrite(session)
	LogError(err)
}
