package database

import (
	"errors"
	"os"

	"github.com/jinzhu/gorm"

	// PosgregSQL communication
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Gorm stores the connection to the database
type Gorm struct {
	Conn *gorm.DB
}

// ConnectToDatabase connects to the PostgreSQL dabase
func ConnectToDatabase() (*Gorm, error) {
	db, err := gorm.Open("postgres", "host="+os.Getenv("DB_HOST")+" port="+os.Getenv("DB_PORT")+" user="+os.Getenv("DB_USER")+" dbname="+os.Getenv("DB_DB_NAME")+" password="+os.Getenv("DB_PASSWORD")+" sslmode=disable")
	if err != nil {
		return nil, err
	}
	if !db.HasTable(&RepoTag{}) {
		db.CreateTable(&RepoTag{})
	}
	if !db.HasTable(&LanguageTag{}) {
		db.CreateTable(&LanguageTag{})
	}
	return &Gorm{Conn: db}, nil
}

// Ping checks database connection
func (db *Gorm) Ping() error {
	if err := db.Conn.DB().Ping(); err != nil {
		return errors.New("There is no connection to the database")
	}
	return nil
}
