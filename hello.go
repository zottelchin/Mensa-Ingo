package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
<<<<<<< HEAD
=======
	"strconv"
	"strings"
>>>>>>> 302b2881c04aedf42f6434a9cfee23e11963ae6a
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

<<<<<<< HEAD
var url string = "https://www.studentenwerk-magdeburg.de/mensen-cafeterien/heute-in-unseren-mensen/"
=======
// TODO: Mensen nicht hardcoden sondern aus Config-Datei laden
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
>>>>>>> 302b2881c04aedf42f6434a9cfee23e11963ae6a

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
					if name.Data == "span" {
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

<<<<<<< HEAD
func gethtml() string {
	cont, _ := ioutil.ReadFile("menu.txt")
	content := string(cont[:])
	return content
}

func main() {
	update()
	go schedultUpdate()
	fmt.Println("Got Mensa Page. Starting Server")
	port := ":8080"
	http.HandleFunc("/", handler)
	http.ListenAndServe(port, nil)
=======
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
>>>>>>> 302b2881c04aedf42f6434a9cfee23e11963ae6a
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
<<<<<<< HEAD
	re4 := regexp.MustCompile("<span\\sclass='grau(?:(?U).*)<\\/strong>")
	for i := 0; i < len(match); i++ {
		match[i][0] = re4.ReplaceAllString(match[i][0], "</strong>")
	}
	re5 := regexp.MustCompile("(style=\"(?:(?U).*)\")|(<\\/tr><\\/thead><tbody>)|(<\\/tr><\\/table>)|(<\\/table>)")
	for i := 0; i < len(match); i++ {
		match[i][0] = re5.ReplaceAllString(match[i][0], "")
	}
	re6 := regexp.MustCompile("<span\\sclass='grau(?:(?U).*)<\\/span>")
	for i := 0; i < len(match); i++ {
		match[i][0] = re6.ReplaceAllString(match[i][0], "")
	}
	re7 := regexp.MustCompile("(<tr\\s?><td\\s(?:(?U).*)>)|(<table><thead><tr><td(?:(?U).*)>)")
	for i := 0; i < len(match); i++ {
		match[i][0] = re7.ReplaceAllString(match[i][0], "<li>")
	}
	re8 := regexp.MustCompile("<\\/td>")
	for i := 0; i < len(match); i++ {
		match[i][0] = re8.ReplaceAllString(match[i][0], "</li>")
	}
	re9 := regexp.MustCompile("<div(?:(?U).*)>")
	for i := 0; i < len(match); i++ {
		match[i][0] = re9.ReplaceAllString(match[i][0], "<div class='mensa'><ul>")
	}
	re10 := regexp.MustCompile("<\\/div>")
	for i := 0; i < len(match); i++ {
		match[i][0] = re10.ReplaceAllString(match[i][0], "</ul></div>")
	}
	re11 := regexp.MustCompile("\\|\\s\\d,\\d\\d")
	for i := 0; i < len(match); i++ {
		match[i][0] = re11.ReplaceAllString(match[i][0], "")
	}
	re12 := regexp.MustCompile("<br\\s\\/>")
	for i := 0; i < len(match); i++ {
		match[i][0] = re12.ReplaceAllString(match[i][0], " ")
	}
	// re13 := regexp.MustCompile("\\d,\\d\\d")
	// for i := 0; i < len(match); i++ {
	// match2 := re13.FindAllStringSubmatch(match[0][0], -1)
	// for j := 0; j < len(match2); j++ {
	// 	match[0][0] = re13.ReplaceAllString(match[0][0], match2[j][0]+"&euro;")
	// 	re13.
	// }
	// }
	content = "<h1>baremetal Mensa Plan</h1><h4>Mensa UniCampus Magdeburg, unterer Saal</h4>"
	for i := 0; i < len(match); i++ {
		content += match[i][0]
=======
	for !isMensaOpenOn(days[0]) {
		days[0] = days[0].Offset(1)
	}

	for i := 1; i < 7; i++ {
		days[i] = days[i-1].Offset(1)
		for !isMensaOpenOn(days[i]) {
			days[i] = days[i].Offset(1)
		}
>>>>>>> 302b2881c04aedf42f6434a9cfee23e11963ae6a
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

func main() {
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
	r.Run()
}

func schedultUpdate() {
	for {
		time.Sleep(2 * time.Hour)
		update()
		fmt.Println("Data updated.")
	}
}
