package main

import (
	"fmt"
	"strconv"
	"time"
)

type Mensa struct {
	Name string
	Open [7]OpeningHours
	Food map[Date][]Dish
}

type Dish struct {
	Name  string
	Price Price
	Hints []string
	Icons []string
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

	return strconv.Itoa(price/100) + "," + strconv.Itoa(price%100) + " €"
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
