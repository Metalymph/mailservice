package main

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const createSqlite string = `
	CREATE TABLE IF NOT EXISTS mails (
	id INTEGER NOT NULL PRIMARY KEY,
	name TEXT,
	address TEXT,
	message TEXT
	);`

const createPostgresql string = `
	CREATE TABLE Mails (
		Id SERIAL INT PRIMARY KEY,
		Name TEXT,
		Address TEXT,
		Message TEXT,
	);
`

//go:embed config.json
var postgresqlConfigFile string

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}

// DBType is a useful alias to define an enumeration for database choice
type DBType int64

const (
	Sqlite DBType = iota
	Postgresql
)

// App represents the web application model
type App struct {
	db  *sql.DB
	mux *http.ServeMux
	dbtype DBType
}

// NewApp generates a pointer to a new App fully initialized
func NewApp(dbtype DBType) (*App, error) {
	app := App{}
	if err := app.initDB(dbtype); err != nil {
		return nil, err
	}
	app.initMux()
	return &app, nil
}

// openDB configures DB settings and initialize it
func openDB(dbtype DBType) (*sql.DB, error) {
	var db *sql.DB

	switch dbtype {
	case Sqlite:
		db, err := sql.Open("sqlite3", "mailservice.db")
		if err != nil {
			return nil, err
		}

		_, err = db.Exec(createSqlite)
		if err != nil {
			return nil, err
		}
	case Postgresql:
		//parsing configuration file for postgresql
		var pgConfig PostgresConfig
		if err := json.NewDecoder(strings.NewReader(postgresqlConfigFile)).Decode(&pgConfig); err != nil {
			return nil, err
		}

		psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", pgConfig.Host, pgConfig.Port, pgConfig.User, pgConfig.Password, pgConfig.Dbname)
		db, err := sql.Open("postgres", psqlconn)
		if err != nil {
			return nil, err
		}

		if err := db.Ping(); err != nil {
			return nil, err
		}

		_, err = db.Exec(createPostgresql)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

// initDB initialize the App internal DB
func (a *App) initDB(dbtype DBType) error {
	db, err := openDB(dbtype)
	if err != nil {
		return err
	}

	a.db = db
	a.dbtype = dbtype
	return nil
}

func (a *App) initMux() {
	a.mux = http.NewServeMux()
	a.mux.HandleFunc("/savemail", a.saveMail)
	a.mux.HandleFunc("/mails", a.getMails)
}

func (a *App) saveMail(w http.ResponseWriter, r *http.Request) {
	var mail MailItem
	if err := json.NewDecoder(r.Body).Decode(&mail); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer r.Body.Close()

	if err := saveMail(a, &mail); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(mail); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (a *App) getMails(w http.ResponseWriter, r *http.Request) {
	mails, err := getMails(a)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if json.NewEncoder(w).Encode(mails); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
