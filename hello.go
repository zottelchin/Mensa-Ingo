package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var url string = "https://www.studentenwerk-magdeburg.de/mensen-cafeterien/heute-in-unseren-mensen/"
var lastUpdateh int
var lastUpdated int

func handler(w http.ResponseWriter, r *http.Request) {
	content := gethtml()
	_, err := fmt.Fprintf(w, "%s", content)
	if err != nil {
		panic(err)
	}
}

func gethtml() string {
	if lastUpdated < time.Now().Day() {
		update()
		fmt.Println("Data updatet")
	} else {
		fmt.Println("Data not updatet")
	}
	cont, _ := ioutil.ReadFile("menu.txt")
	content := string(cont[:])
	return content
}

func main() {
	update()
	fmt.Println("Got Mensa Page. Starting Server")
	port := ":8080"
	http.HandleFunc("/", handler)
	http.ListenAndServe(port, nil)
}
func update() {
	fmt.Println("Updating file!")
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	cont, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	read, _ := ioutil.ReadFile("start.html")
	htmlStart := string(read[:])
	read, _ = ioutil.ReadFile("End.html")
	htmlEnd := string(read[:])
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

	cont = []byte(htmlStart + content + htmlEnd)
	ioutil.WriteFile("menu.txt", cont, 0600)

}
