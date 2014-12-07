package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"time"
)

func setupDatabase() {
	os.Remove("./kkez.db")

	db, err := sql.Open("sqlite3", "./kkez.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table events (id integer not null primary key, name varchar(255), date date);
	delete from events;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
	log.Printf("Set up a database")
	return
}

func allEvents() []Event {
	db, err := sql.Open("sqlite3", "./kkez.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("select id, name, date from events")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var events []Event
	for rows.Next() {
		var id int
		var name string
		var date time.Time
		rows.Scan(&id, &name, &date)
		events = append(events, Event{Id: id, Name: name, Date: date})
	}
	return events
}
func nextEvent(on time.Time) (event Event) {
	db, err := sql.Open("sqlite3", "./kkez.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var name string
	var date time.Time
	db.QueryRow("SELECT name, date from events where date >= ? LIMIT 1", on).Scan(&name, &date)
	return Event{Name: name, Date: date}
}

func createEvent(name string, date time.Time) {
	db, err := sql.Open("sqlite3", "./kkez.db")
	defer db.Close()
	var count int
	db.QueryRow("SELECT COUNT(*) from events where name = ? and date = ?", name, date).Scan(&count)
	log.Printf("Count: %d", count)
	if count < 1 {
		_, err = db.Exec("insert into events(name, date) values(?, ?)", name, date)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Saved the record")
	} else {
		log.Printf("Found the record")
	}
	return
}
