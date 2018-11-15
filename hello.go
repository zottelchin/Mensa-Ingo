package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

// Load Uni Mensa Magdeburg on default
var urls = []string{
	"https://www.studentenwerk-magdeburg.de/mensen-cafeterien/mensa-unicampus/speiseplan-unten/",
	"https://www.studentenwerk-magdeburg.de/mensen-cafeterien/mensa-unicampus/speiseplan-oben/",
	"https://www.studentenwerk-magdeburg.de/mensen-cafeterien/mensa-kellercafe/speiseplan/",
}
var menu = []Mensa{
	Mensa{
		"UniCampus: Unterer Saal",
		[7]OpeningHours{
			NewOpeningHours(10, 45, 14, 00), NewOpeningHours(10, 45, 14, 00), NewOpeningHours(10, 45, 14, 00), NewOpeningHours(10, 45, 14, 00), NewOpeningHours(10, 45, 14, 00),
			NewOpeningHours(12, 00, 13, 30),
			OpeningHours{},
		},
		map[Date][]Meal{},
		map[Date][]string{},
	},
	Mensa{
		"UniCampus: Oberer Saal",
		[7]OpeningHours{
			NewOpeningHours(10, 45, 15, 15), NewOpeningHours(10, 45, 15, 15), NewOpeningHours(10, 45, 15, 15), NewOpeningHours(10, 45, 15, 15),
			NewOpeningHours(10, 45, 14, 30),
			OpeningHours{}, OpeningHours{},
		},
		map[Date][]Meal{},
		map[Date][]string{},
	},
	Mensa{
		"Kellercafé Zschokkestraße",
		[7]OpeningHours{
			NewOpeningHours(10, 45, 14, 00), NewOpeningHours(10, 45, 14, 00), NewOpeningHours(10, 45, 14, 00), NewOpeningHours(10, 45, 14, 00), NewOpeningHours(10, 45, 14, 00),
			OpeningHours{}, OpeningHours{},
		},
		map[Date][]Meal{},
		map[Date][]string{},
	},
}

var headerRegex = regexp.MustCompile(`^.*\s(\d+\.\d+\.\d+)\s*$`)

func updateMenu() {
	log.Println("Updating menu...")
	for id := 0; id < len(urls); id++ {
		log.Printf("Now processing: %s at %s", menu[id].Name, urls[id])
		doc, err := goquery.NewDocument(urls[id])
		if err != nil {
			log.Print(err)
			continue
		}

		// Each day in the table
		doc.Find(".entry-content .mensa table").Each(func(i int, e *goquery.Selection) {
			// Parse date to integer
			datestring := strings.Split(headerRegex.ReplaceAllString(e.Find("thead td").Text(), "$1"), ".")
			y, e1 := strconv.Atoi(datestring[2])
			m, e2 := strconv.Atoi(datestring[1])
			d, e3 := strconv.Atoi(datestring[0])
			if e1 != nil || e2 != nil || e3 != nil {
				log.Printf("Error processing date: %s\n %s; %s; %s", datestring, e1, e2, e3)
				return
			}
			date := Date{y, time.Month(m), d}

			// Each meal in the table
			results := e.Find("tbody tr")
			menu[id].Meals[date] = make([]Meal, results.Length()-1)
			results.Each(func(i int, e *goquery.Selection) {
				if i >= results.Length()-1 {
					menu[id].Sides[date] = strings.Split(strings.TrimPrefix(e.Get(0).FirstChild.FirstChild.Data, "Beilagen: "), ", ")
				} else {
					// Parse price
					price := strings.Split(strings.Replace(strings.TrimSpace(e.Get(0).FirstChild.LastChild.Data), ",", "", -1), " | ")
					student, e1 := strconv.Atoi(price[0])
					staff, e2 := strconv.Atoi(price[1])
					guest, e3 := strconv.Atoi(price[2])
					if e1 != nil || e2 != nil || e3 != nil {
						log.Printf("Error processing pricing: %s\n %s; %s; %s", price, e1, e2, e3)
						return
					}

					// Parse icons
					iconResults := e.Find("img")
					icons := make([]Icon, iconResults.Length())
					iconResults.Each(func(i int, e *goquery.Selection) {
						icons[i] = Icon{
							strings.TrimSuffix(strings.TrimPrefix(e.AttrOr("src", ""), "/wp-content/themes/swmd2012/mensasym/mensasym_"), ".png"),
							strings.TrimPrefix(e.AttrOr("alt", ""), "Symbol "),
						}
					})

					// Parse name
					name := e.Find("strong").Get(0).FirstChild
					for name.Data == "span" {
						name = name.FirstChild
					}

					// Parse hints
					hintElement := e.Get(0).LastChild.LastChild
					hints := []string{}
					if hintElement.FirstChild != nil {
						hints = strings.Split(strings.Trim(hintElement.FirstChild.Data, "() "), ") (")
					}

					// Insert
					menu[id].Meals[date][i] = Meal{
						name.Data,
						Price{student, staff, guest},
						hints,
						icons,
					}
				}
			})
		})
	}
}

func scheduleUpdate() {
	for {
		time.Sleep(1 * time.Hour)
		updateMenu()
	}
}

func isMensaOpenOn(d Date) bool {
	for _, m := range menu {
		if m.Open[d.Weekday()].Format() != "geschlossen" {
			return true
		}
	}
	return false
}
func isMensaStillOpen() bool {
	for _, m := range menu {
		if !m.Open[time.Now().Weekday()].AlreadyClosed() {
			return true
		}
	}
	return false
}

func handler(c *gin.Context) {
	offset := c.Param("offset")
	day := 0
	if offset != "" {
		var err error
		day, err = strconv.Atoi(offset)
		if !(err == nil && day >= 1 && day <= 6) {
			notFound(c)
			return
		}
	}

	t := Today()
	days := make([]Date, 7)
	days[0] = t
	if !isMensaStillOpen() {
		days[0] = days[0].Offset(1)
	}
	for !isMensaOpenOn(days[0]) {
		days[0] = days[0].Offset(1)
	}

	for i := 1; i < 7; i++ {
		days[i] = days[i-1].Offset(1)
		for !isMensaOpenOn(days[i]) {
			days[i] = days[i].Offset(1)
		}
	}

	c.HTML(http.StatusOK, "mensa.html", TemplateData{
		days,
		day,
		[]string{"Mo", "Di", "Mi", "Do", "Fr", "Sa", "So"},
		menu,
	})
}

func notFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.html", nil)
}

func loadConfig() {
	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Errorf("%v", err)
		return
	}
	defer configFile.Close()

	byteValue, _ := ioutil.ReadAll(configFile)
	config := Config{}
	json.Unmarshal(byteValue, &config)

	//parse Mensen to global vars
	urls = []string{}
	menu = []Mensa{}
	for _, m := range config.Mensen {
		urls = append(urls, m.URL)
		tmp := Mensa{
			m.Name,
			doSomeMagic(m.OpeningHours),
			map[Date][]Meal{},
			map[Date][]string{},
		}
		menu = append(menu, tmp)
	}
}

func doSomeMagic(open jsonOpeningHours) [7]OpeningHours {
	Mo := OpeningHours{}
	if open.Mo[0] != -1 {
		Mo = NewOpeningHours(open.Mo[0], open.Mo[1], open.Mo[2], open.Mo[3])
	}
	Di := OpeningHours{}
	if open.Di[0] != -1 {
		Di = NewOpeningHours(open.Di[0], open.Di[1], open.Di[2], open.Di[3])
	}
	Mi := OpeningHours{}
	if open.Mi[0] != -1 {
		Mi = NewOpeningHours(open.Mi[0], open.Mi[1], open.Mi[2], open.Mi[3])
	}
	Do := OpeningHours{}
	if open.Do[0] != -1 {
		Do = NewOpeningHours(open.Do[0], open.Do[1], open.Do[2], open.Do[3])
	}
	Fr := OpeningHours{}
	if open.Fr[0] != -1 {
		Fr = NewOpeningHours(open.Fr[0], open.Fr[1], open.Fr[2], open.Fr[3])
	}
	Sa := OpeningHours{}
	if open.Sa[0] != -1 {
		Sa = NewOpeningHours(open.Sa[0], open.Sa[1], open.Sa[2], open.Sa[3])
	}
	So := OpeningHours{}
	if open.So[0] != -1 {
		So = NewOpeningHours(open.So[0], open.So[1], open.So[2], open.So[3])
	}
	return [7]OpeningHours{Mo, Di, Mi, Do, Fr, Sa, So}
}

func main() {
	// Load Config
	loadConfig()
	// Load Menu
	updateMenu()
	go scheduleUpdate()

	// Start Server
	r := gin.Default()
	r.LoadHTMLGlob("*.html")

	r.GET("/", handler)
	r.GET("/+:offset", handler)

	filepath.Walk("static", func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			r.StaticFile(strings.TrimPrefix(path, "static/"), path)
		}
		return nil
	})

	r.NoRoute(notFound)
	r.Run(":8700")
}
