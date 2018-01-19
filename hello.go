package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
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
	}

	cont = []byte(htmlStart + content + htmlEnd)
	ioutil.WriteFile("menu.txt", cont, 0600)

}

func schedultUpdate() {
	for {
		time.Sleep(2 * time.Hour)
		update()
		fmt.Println("Data updated.")
	}
}
