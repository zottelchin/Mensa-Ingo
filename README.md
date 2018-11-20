# betterMensa-in-Go

The idea comes from @t_aus_m and the design was made by @moqmar. It makes a simple Webserver in golang to display our uni mensa page without the unnessasary things.

This Side was build with mainly in [GO](https://golang.org). We used [Gin](https://github.com/gin-gonic/gin) as Webframework.

Loads Mensa Settings from JSON-File named `config.json`
```json
{
    "mensen": [
        {
            "name": "UniCampus: Unterer Saal",
            "url": "https://www.studentenwerk-magdeburg.de/mensen-cafeterien/mensa-unicampus/speiseplan-unten/",
            "openingHours": {
                "Mo": [10,45,14,0],
                "Di": [10,45,14,0],
                "Mi": [10,45,14,0],
                "Do": [10,45,14,0],
                "Fr": [10,45,14,0],
                "Sa": [12,0,13,30],
                "So": [-1]
            }
        },
        {
            "name": "UniCampus: Oberer Saal",
            "url": "https://www.studentenwerk-magdeburg.de/mensen-cafeterien/mensa-unicampus/speiseplan-oben/",
            "openingHours": {
                "Mo": [10,45,15,15],
                "Di": [10,45,15,15],
                "Mi": [10,45,15,15],
                "Do": [10,45,15,15],
                "Fr": [10,45,14,30],
                "Sa": [-1],
                "So": [-1]
            }
        }
    ]
}
```
The Opening Hours are formated like [Opening Hour, Opening Minute, Closing Hour, Closing Minute]. If the Opening Hour is -1 The Mensa is closed on this Weekday.