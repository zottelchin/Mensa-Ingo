package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var url = "https://www.studentenwerk-magdeburg.de/mensen-cafeterien/heute-in-unseren-mensen/"
var menu = ""

func gethtml() string {
	// Get content
	cont := menu
	content := string(cont[:])

	re := regexp.MustCompile("(<div\\sclass='mensa'>.*</div>)|(<h4>.*</h4>)")
	match := re.FindAllStringSubmatch(content, -1)
	re3 := regexp.MustCompile("(<td\\sstyle='t(?:(?U).*)</td>)|(<a(?:(?U).*)>)|(</a>)")
	for i := 0; i < len(match); i++ {
		match[i][0] = re3.ReplaceAllString(match[i][0], "")
	}
	re4 := regexp.MustCompile("<span\\sclass='grau(?:(?U).*)</strong>")
	for i := 0; i < len(match); i++ {
		match[i][0] = re4.ReplaceAllString(match[i][0], "</strong>")
	}
	lastUpdateh = time.Now().Hour()
	lastUpdated = time.Now().Day()

	content = "<p>Letztes Update der Daten am " + strconv.Itoa(lastUpdated) + ". um " + strconv.Itoa(lastUpdateh) + " Uhr</p><h4>Mensa UniCampus Magdeburg, unterer Saal</h4>"
	for i := 0; i < len(match); i++ {
		content += match[i][0]
	}

	return content
}

func updateMenu() error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	menu = string(content)
	return nil
}

func scheduleUpdate() {
	for {
		time.Sleep(1 * time.Hour)
		updateMenu()
	}
}

func handler(c *gin.Context) {
	offset := c.Param("offset")
	if offset != "" {
		i, err := strconv.Atoi(offset)
		if !(err == nil && i >= 1 && i <= 6) {
			notFound(c)
			return
		}
	}
	c.HTML(http.StatusOK, "mensa.html", map[string]interface{}{
		"content": template.HTML(gethtml()),
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
