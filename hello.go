package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var url = "https://www.studentenwerk-magdeburg.de/mensen-cafeterien/heute-in-unseren-mensen/"
var menu = []Mensa{
	Mensa{
		"Unterer Saal der Mensa UniCampus",
		[7]OpeningHours{
			OpeningHours{10*time.Hour + 45*time.Minute, 14 * time.Hour},
			OpeningHours{10*time.Hour + 45*time.Minute, 14 * time.Hour},
			OpeningHours{10*time.Hour + 45*time.Minute, 14 * time.Hour},
			OpeningHours{10*time.Hour + 45*time.Minute, 14 * time.Hour},
			OpeningHours{10*time.Hour + 45*time.Minute, 14 * time.Hour},
			OpeningHours{12 * time.Hour, 13*time.Hour + 30*time.Minute},
			OpeningHours{},
		},
		map[Date][]Dish{
			Date{2018, time.January, 19}: []Dish{
				Dish{
					"Aktion: Gemüse-Knusperschnitzel Gärtnerin Art mit Sauce Hollandaise",
					Price{135, 240, 310},
					[]string{"a1", "c", "g", "i"},
					[]string{"vegetarisch", "knoblauch"},
				},
			},
		},
	},
}

func updateMenu() error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	_, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	//menu = string(content)
	return nil
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

func main() {
	var err error

	// Load Menu
	if err = updateMenu(); err != nil {
		panic(err)
	}
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
