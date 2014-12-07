package main

import "net/http"
import "log"
import "os"
import "fmt"
import "path"
import "html/template"
import "time"
import "sort"

type LayoutData struct {
	Title string
	Data  interface{}
}

type Index struct {
	Title    string
	Content  string
	Subtitle string
}

type Event struct {
	Id   int
	Name string
	Date time.Time
}

func dateFromString(stamp string) time.Time {
	t, _ := time.Parse("2006, 1, 2", stamp)
	return t
}

var events = []Event{
	Event{Name: "Pier", Date: dateFromString("2012, 2, 24")},
	Event{Name: "Frits", Date: dateFromString("2012, 3, 30")},
	Event{Name: "Karin", Date: dateFromString("2012, 4, 27")},
	Event{Name: "Simon", Date: dateFromString("2012, 5, 25")},
	Event{Name: "Stijn", Date: dateFromString("2012, 6, 22")},
	Event{Name: "Dennis", Date: dateFromString("2012, 7, 27")},
	Event{Name: "Erwin", Date: dateFromString("2012, 8, 31")},
	Event{Name: "Kurt", Date: dateFromString("2012, 9, 28")},
	Event{Name: "Peter", Date: dateFromString("2012, 10,	26")},
	Event{Name: "Katrien", Date: dateFromString("2012, 11,	30")},
	Event{Name: "Gijs", Date: dateFromString("2012, 12,	28")},
	Event{Name: "Jan", Date: dateFromString("2013, 1, 25")},
	Event{Name: "Piet", Date: dateFromString("2013, 2, 22")},
	Event{Name: "Maarten", Date: dateFromString("2013, 3, 29")},
	Event{Name: "Koen", Date: dateFromString("2013, 4, 26")},
	Event{Name: "Frits", Date: dateFromString("2013, 5, 31")},
	Event{Name: "Dennis", Date: dateFromString("2013, 6, 28")},
	Event{Name: "Stijn", Date: dateFromString("2013, 7, 26")},
	Event{Name: "Dieter", Date: dateFromString("2013, 8, 30")},
	Event{Name: "Kurt", Date: dateFromString("2013, 9, 27")},
	Event{Name: "Peter", Date: dateFromString("2013, 10, 25")},
	Event{Name: "Katrien", Date: dateFromString("2013, 11, 29")},
	Event{Name: "Gijs", Date: dateFromString("2013, 12, 27")},
	Event{Name: "Jan", Date: dateFromString("2014, 1, 31")},
	Event{Name: "Piet", Date: dateFromString("2014, 2, 28")},
	Event{Name: "Maarten", Date: dateFromString("2014, 3, 28")},
	Event{Name: "Koen", Date: dateFromString("2014, 4, 26")},
	Event{Name: "Pier", Date: dateFromString("2014, 5, 30")},
	Event{Name: "Karin", Date: dateFromString("2014, 6, 27")},
	Event{Name: "Simon", Date: dateFromString("2014, 7, 25")},
	Event{Name: "Erwin", Date: dateFromString("2014, 8, 29")},
}

type ByDate []Event

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }

func nextCardMoment(sortedEvents []Event, date time.Time) (event Event, isToday bool) {
	for _, event := range sortedEvents {
		if event.Date.Equal(date) || event.Date.After(date) {
			return event, (event.Date.YearDay() == time.Now().YearDay())
		}
	}
	return Event{Name: "Niemand"}, false
}

func main() {
	// setupDatabase()
	for _, event := range events {
		createEvent(event.Name, event.Date)
	}
	log.Printf("Event %s", nextEvent(dateFromString("2014, 12, 23")))
	log.Printf("All events: %s", allEvents())
	sort.Sort(ByDate(events))

	// Serve our static assets
	// (Make sure to strip out the public, before `FileServer` sees it)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		layout := path.Join("templates", "layout.html")
		parse := func(name string) (*template.Template, error) {
			t := template.New("")
			return t.ParseFiles(layout, path.Join("templates", name))
		}
		index, err := parse("index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var layoutData = LayoutData{Title: "Dinges"}
		event, isToday := nextCardMoment(events, time.Now())
		if isToday {
			layoutData.Data = Index{"Is het kaarten vandaag?", "Ja", fmt.Sprintf("Bij %s", event.Name)}
		} else {
			layoutData.Data = Index{"Is het kaarten vandaag?", "Neeje", fmt.Sprintf("Volgende is bij %s op %s", event.Name, event.Date.Format("Mon Jan 2"))}
		}

		if err := index.ExecuteTemplate(w, "layout.html", layoutData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}
	log.Printf("Whizzing and wurring at http://0.0.0.0:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
