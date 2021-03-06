package main

import (
	"fmt"
	"strconv"
	"time"
)

type Mensa struct {
	Name  string
	Open  [7]OpeningHours
	Meals map[Date][]Meal
	Sides map[Date][]string
}

type Meal struct {
	Name  string
	Price Price
	Hints []string
	Icons []Icon
}

type Icon struct {
	Name  string
	Title string
}

type Price struct {
	Student int
	Staff   int
	Guest   int
}

func (p Price) Format(t string) string {
	var price int
	if t == "student" {
		price = p.Student
	} else if t == "staff" {
		price = p.Staff
	} else if t == "guest" {
		price = p.Guest
	} else {
		return "---"
	}

	return strconv.Itoa(price/100) + "," + fmt.Sprintf("%02d", price%100) + " €"
}

type OpeningHours struct {
	Opens  time.Duration
	Closes time.Duration
}

func (o OpeningHours) Format() string {
	if o.Closes <= o.Opens {
		return "geschlossen"
	}
	return strconv.Itoa(int(o.Opens.Hours())) + ":" + fmt.Sprintf("%02d", int(o.Opens.Minutes())%60) +
		" - " + strconv.Itoa(int(o.Closes.Hours())) + ":" + fmt.Sprintf("%02d", int(o.Closes.Minutes())%60) +
		" Uhr"
}

func (o OpeningHours) AlreadyClosed() bool {
	n := time.Now()
	h := time.Duration(n.Hour()) * time.Hour
	m := time.Duration(n.Minute()) * time.Minute
	return h+m >= o.Closes
}

func NewOpeningHours(oh, om, ch, cm int) OpeningHours {
	return OpeningHours{
		time.Duration(oh)*time.Hour + time.Duration(om)*time.Minute,
		time.Duration(ch)*time.Hour + time.Duration(cm)*time.Minute,
	}
}

type Date struct {
	Year  int
	Month time.Month
	Day   int
}

func (d Date) MonthInt() int {
	return int(d.Month)
}

func (d Date) Offset(days int) Date {
	t := time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.UTC).AddDate(0, 0, days)
	return Date{t.Year(), t.Month(), t.Day()}
}

func (d Date) IsToday() bool {
	n := time.Now()
	return d.Year == n.Year() && d.Month == n.Month() && d.Day == n.Day()
}

func (d Date) Weekday() int {
	t := int(time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.UTC).Weekday()) - 1
	if t < 0 {
		return 6
	}
	return t
}

func Today() Date {
	t := time.Now()
	return Date{t.Year(), t.Month(), t.Day()}
}

type TemplateData struct {
	Days     []Date
	Day      int
	Weekdays []string
	Menu     []Mensa
}

type Config struct {
	Mensen []struct {
		Name         string           `json:"name"`
		URL          string           `json:"url"`
		OpeningHours jsonOpeningHours `json:"openingHours"`
	} `json:"mensen"`
}

type jsonOpeningHours struct {
	Mo []int `json:"Mo"`
	Di []int `json:"Di"`
	Mi []int `json:"Mi"`
	Do []int `json:"Do"`
	Fr []int `json:"Fr"`
	Sa []int `json:"Sa"`
	So []int `json:"So"`
}
