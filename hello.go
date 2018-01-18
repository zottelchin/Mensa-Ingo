package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

var url string = "https://www.studentenwerk-magdeburg.de/mensen-cafeterien/heute-in-unseren-mensen/"

func handler(w http.ResponseWriter, r *http.Request) {
	content := gethtml()
	_, err := fmt.Fprintf(w, "%s", content)
	if err != nil {
		panic(err)
	}
}

func gethtml() string {
	read, _ := ioutil.ReadFile("start.html")
	htmlStart := string(read[:])
	read, _ = ioutil.ReadFile("End.html")
	htmlEnd := string(read[:])
	cont, _ := ioutil.ReadFile("menu.txt")
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
	content = "<h4>Mensa UniCampus Magdeburg, unterer Saal</h4>"
	for i := 0; i < len(match); i++ {
		content += match[i][0]
	}
	return htmlStart + content + htmlEnd
}

func main() {
	// startup()
	// port := ":8080"
	// http.HandleFunc("/", handler)
	// http.ListenAndServe(port, nil)
	// fmt.Println(port)
	// fmt.Print(gethtml())
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
}
func startup() {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	cont, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile("menu.txt", cont, 0600)
}
