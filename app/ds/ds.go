package ds

import (
	"log"
	"upper.io/db.v2"
	"upper.io/db.v2/postgresql"
)

var sess postgresql.Database

var (
	userSource     					db.Collection
	leadersSource 					db.Collection
	query										db.Result
)

func init() {
	settings := postgresql.ConnectionURL{
		Database: "cc",
		Host:     "localhost",
		User:     "danyel",
		Password: "passsword",
	}

	// Conexion a la DB y comunicarse con las tables
	var err error
	sess, err = postgresql.Open(settings)
	if err != nil {
		log.Fatal("Session Open Error: ", err)
	}

	userSource = sess.Collection("users")
	leadersSource = sess.Collection("leaders_public_footprint")
}
