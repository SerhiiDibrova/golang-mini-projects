package helpers

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func LogError(message string) {
	logErrorToFile(message)
	logErrorToDatabase(message)
	logErrorToConsole(message)
}

func logErrorToFile(message string) {
	file, err := os.OpenFile("error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println(message)
}

func logErrorToDatabase(message string) {
	db, err := sql.Open("sqlite3", "./error.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS errors (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		message TEXT
	);`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO errors (message) VALUES (?)", message)
	if err != nil {
		log.Fatal(err)
	}
}

func logErrorToConsole(message string) {
	logger := log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println(message)
}