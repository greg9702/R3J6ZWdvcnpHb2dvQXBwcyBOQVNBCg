package utils

import (
	"time"
)

func GetListOfDate(startDate time.Time, endDate time.Time) []string {

	var dateStringList []string

	dateStringList = append(dateStringList, "?date=2020-07-01")
	dateStringList = append(dateStringList, "?date=2020-07-02")
	dateStringList = append(dateStringList, "?date=2020-07-03")
	dateStringList = append(dateStringList, "?date=2020-07-04")
	dateStringList = append(dateStringList, "?date=2020-07-05")
	dateStringList = append(dateStringList, "?date=2020-07-06")
	dateStringList = append(dateStringList, "?date=2020-07-07")

	return dateStringList
}
