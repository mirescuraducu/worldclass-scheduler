package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"
)

type Credentials struct {
	email    string
	password string
}

type ClassSchedule struct {
	clubId  string
	classId string
}

type Bookings struct {
	name             string
	bookingWeekDay   string
	bookingStartTime string
	class            ClassSchedule
	account          []Credentials
}

func main() {
	now := time.Now()                 // current local time
	weekday := now.Weekday().String() // get the day of the week name as a string
	hour, minutes, _ := now.Clock()
	timeString := fmt.Sprintf("%02d:%02d", hour, minutes)

	time.Sleep(time.Duration(rand.Intn(10)))

	radu := Credentials{"mirescu.raducu@gmail.com", os.Getenv("RADU_PASSKEY")}
	alina := Credentials{"alina.tucunete@gmail.com", os.Getenv("ALINA_PASSKEY")}
	bookings := []Bookings{
		Bookings{ // TRX
			name:             "TRX - Luni",
			bookingWeekDay:   "Sunday",
			bookingStartTime: "16:00",
			class:            ClassSchedule{"464", "730835"},
			account:          []Credentials{radu, alina},
		},
		Bookings{ // TRX
			name:             "TRX - Miercuri",
			bookingWeekDay:   "Tuesday",
			bookingStartTime: "17:10",
			class:            ClassSchedule{"464", "730913"},
			account:          []Credentials{radu, alina},
		},
		Bookings{ // Pilates
			name:             "Pilates - Marti",
			bookingWeekDay:   "Monday",
			bookingStartTime: "17:40",
			class:            ClassSchedule{"410", "733153"},
			account:          []Credentials{alina},
		},
		Bookings{ // Pilates
			name:             "Pilates - Joi",
			bookingWeekDay:   "Wednesday",
			bookingStartTime: "16:30",
			class:            ClassSchedule{"410", "733205"},
			account:          []Credentials{alina},
		},
		Bookings{ // Zumba
			name:             "Zumba - Vineri",
			bookingWeekDay:   "Thursday",
			bookingStartTime: "16:30",
			class:            ClassSchedule{"410", "733244"},
			account:          []Credentials{alina},
		},
	}

	baseUrl, err := url.Parse("https://members.worldclass.ro")
	if err != nil {
		panic(err)
	}

	bookingAvailable := false
	for _, booking := range bookings {
		if booking.bookingWeekDay == weekday && booking.bookingStartTime == timeString {
			bookingAvailable = true
			for _, account := range booking.account {
				cookies := login(account, baseUrl)
				booked := schedule(cookies, booking.class, *baseUrl)
				if booked {
					log(fmt.Sprintf("Reserved class: %s, email: %s", booking.name, account.email))
				}
			}
		}
	}

	if !bookingAvailable {
		log("No booking found!")
	}
	// cookies := login(radu, baseUrl)
	// schedule(cookies, ClassSchedule{"455", "725731"}, *baseUrl)
}

func login(credentials Credentials, baseUrl *url.URL) []*http.Cookie {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.PostForm(
		baseUrl.JoinPath("_process_login.php").String(),
		url.Values{
			"email":           {credentials.email},
			"member_password": {credentials.password},
			"remember_me":     {"false"},
		})

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 302 && resp.Header.Get("Location") == baseUrl.JoinPath("dashboard.php").String() {
		return resp.Cookies()
	}

	panic("Invalid credentials")
}

func schedule(cookies []*http.Cookie, classToSchedule ClassSchedule, baseUrl url.URL) bool {
	scheduleUrl := baseUrl.JoinPath("_book_class.php")
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: cookieJar,
	}

	queryParams := url.Values{}
	queryParams.Set("id", classToSchedule.classId)
	queryParams.Set("clubid", classToSchedule.clubId)
	scheduleUrl.RawQuery = queryParams.Encode() // encode and attach the query string

	scheduleRequest, err := http.NewRequest("GET", scheduleUrl.String(), nil)
	if err != nil {
		panic(err)
	}

	for _, cookie := range cookies {
		scheduleRequest.AddCookie(cookie)
	}

	scheduleRequestResult, err := client.Do(scheduleRequest)
	if err != nil {
		panic(err)
	}

	defer scheduleRequestResult.Body.Close()

	if scheduleRequestResult.StatusCode == 302 && scheduleRequestResult.Header.Get("Location") == baseUrl.JoinPath("member-schedule.php").String() {
		return true
	}

	return false

}

func log(message string) {
	now := time.Now()
	fmt.Println(fmt.Sprintf("[%s] %s", now.Format(time.DateTime), message))
}
