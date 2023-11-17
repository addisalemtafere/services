package psql

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/lib/pq"
)

func New(log *log.Logger) (*sql.DB, error) {

	// Errors
	const (
		ErrFailedToConnect = "FAILED_TO_CONNECT"
		ErrFailedToPing    = "FAILED_TO_PING"
		ErrFailedToCreate  = "FAILED_TO_CREATE"
	)

	log.SetPrefix(log.Prefix() + " [AUTH] [INFRA] [STORAGE] [PSQL] ")
	log.Println("Initiating pg db")

	// Configure Database
	// Get host, username, password from env
	var (
		host, _     = os.LookupEnv("DB_HOST")
		user, _     = os.LookupEnv("DB_USER")
		password, _ = os.LookupEnv("DB_PASS")
		dbName, _   = os.LookupEnv("DB_NAME")
		sslMode, _  = os.LookupEnv("SSL_MODE")
	)

	// pg connection
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", user, password, host, dbName, sslMode))
	if err != nil {
		log.Println("Failed to connect to db server : " + err.Error())
		return nil, err
	}

	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(0)

	if err = db.Ping(); err != nil {
		log.Println("Ping is not responding")
		if err, ok := err.(*pq.Error); ok {
			switch err.Code.Name() {
			case "invalid_catalog_name":
				{
					log.Println("The specified database does not exist, trying to create one")
					_, err := db.Exec("create database " + dbName)
					if err != nil {
						//handle the error
						log.Println(err.Error())
					}
					if err = db.Ping(); err != nil {
						log.Println("Pinging the new database failed : " + err.Error())
					}
					return db, nil
				}
			}
			// log.Println(err.Code.Name())
			return nil, err
		}
		return nil, err

	}

	log.Println("Successfully connected to database")
	return db, nil
}
