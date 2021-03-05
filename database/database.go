package database

import (
	"log"

	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

func getTableDefinitions() map[string]interface{} {
	return map[string]interface{}{
		"level_metadata": "",
		"users":          "",
		"levels":         "",
	}
}

func Init(db string, url string, reseed bool) *r.Session {
	var err error

	session, err := r.Connect(r.ConnectOpts{
		Address: url,
	})

	if err != nil {
		log.Fatalln(err)
	}

	SetUpDatabase(db, session, reseed)

	return session
}

func SetUpDatabase(DBName string, session *r.Session, reseed bool) {

	if reseed {
		r.DBDrop(DBName).Exec(session)
	}

	createDatabaseIfNotExist(DBName, session)
	createTablesIfNotExist(DBName, session)
	setupDatabaseUsers(DBName, session)
}

func createDatabaseIfNotExist(DBName string, session *r.Session) {
	// Create database if it doesn't exist
	err := r.DB(DBName).TableList().Exec(session)

	if err != nil {
		log.Printf("Error retrieving database: %v", err)

		// There was an error retrieving the database, attempt to create it
		_, err = r.DBCreate(DBName).Run(session)

		if err != nil {
			log.Fatalf("Failed to create database %s: %v", DBName, err)
		}

		log.Printf("Successfully created database %s; creating tables", DBName)
	}
}

func createTablesIfNotExist(DBName string, session *r.Session) {
	for k := range getTableDefinitions() {
		if err := r.DB(DBName).Table(k).IndexList().Exec(session); err != nil {
			log.Printf("%s table doesn't exist, attempting to create: %v", k, err)

			err = r.DB(DBName).TableCreate(k).Exec(session)

			if err != nil {
				log.Fatalf("Failed to create table %s: %v", k, err)
			}
			log.Printf("Table %s successfully created!", k)
		} else {
			log.Printf("Table %s exists.", k)
		}

	}
}

func setupDatabaseUsers(DBName string, session *r.Session) {
	err := r.DB("rethinkdb").Table("users").Insert(map[string]string{
		"id":       "system",
		"password": "p455w0rd", // TODO: MAKE THIS EXTERNAL
	}).Exec(session)

	if err != nil {
		log.Fatalf("Could not set up database users: %v", err)
	}
}
