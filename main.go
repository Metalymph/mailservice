package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	app, err := NewApp(Sqlite)
	app.db.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err.Error())
		return
	}
	log.Fatal(http.ListenAndServe(":8080", app.mux))
}
