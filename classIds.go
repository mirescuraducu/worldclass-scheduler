package main

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/gocolly/colly"
)

type Credentials struct {
	email    string
	password string
}

type Club struct {
	id   string
	name string
}
type class struct {
	day      string
	hour     string
	title    string
	trainer  string
	room     string
	classId  string
	clubName string
}

func main() {
	classes := []class{}
	clubName := "Not defined"
	radu := Credentials{"mirescu.raducu@gmail.com", os.Getenv("RADU_PASSKEY")}
	baseUrl, err := url.Parse("https://members.worldclass.ro")
	if err != nil {
		panic(err)
	}

	// create a new collector
	c := colly.NewCollector()

	// authenticate
	c.Post(baseUrl.JoinPath("_process_login.php").String(), map[string]string{"email": radu.email, "member_password": radu.password, "remember_me": "false"})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnHTML(".daily-schedule", func(e *colly.HTMLElement) {
		temp := class{}
		temp.day = e.ChildText("div.schedule-day>strong")

		e.ForEach(".schedule-class", func(_ int, el *colly.HTMLElement) {
			if len(el.DOM.Find(".btn-book-class").Nodes) > 0 {
				temp.hour = el.ChildText("div.col-xs-7.col-sm-12>span.class-hours")
				temp.room = el.ChildText("div.col-xs-7.col-sm-12>span.room")
				temp.title = el.ChildText("div.col-xs-7.col-sm-12>strong.class-title")
				temp.trainer = el.ChildText("div.col-xs-7.col-sm-12>span.trainers")
				temp.classId = el.ChildAttr("div.col-xs-5.col-sm-12.text-right>a", "data-target")
				temp.clubName = clubName
				classes = append(classes, temp)
			}
		})

	})

	for _, clubObj := range []Club{Club{"446", "Caro"}, Club{"433", "Promenada"}, Club{"410", "Upground"}, Club{"464", "Oregon"}, Club{"466", "Planet"}} {
		clubName = clubObj.name
		c.Post(baseUrl.JoinPath("member-schedule.php").String(), map[string]string{"clubid": clubObj.id, "group": "-1"})
	}

	c.Wait()

	for _, classFound := range classes {
		_, currentWeekNumber := time.Now().ISOWeek()
		fmt.Println(currentWeekNumber, classFound.clubName, classFound.title, classFound.classId, classFound.trainer, classFound.hour)
	}
}
