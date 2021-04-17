package utils

import (
	"errors"
	"time"
)

func GetListOfDate(startDate, endDate time.Time) ([]string, error) {

	var dateStringList []string

	if endDate.Before(startDate) {
		return dateStringList, errors.New("End date is before start date")
	}

	for rd := rangeDate(startDate, endDate); ; {
		date := rd()
		if date.IsZero() {
			break
		}
		dateStringList = append(dateStringList, date.Format("2006-01-02"))
	}
	return dateStringList, nil
}

func rangeDate(start, end time.Time) func() time.Time {
	y, m, d := start.Date()
	start = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	y, m, d = end.Date()
	end = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	return func() time.Time {
		if start.After(end) {
			return time.Time{}
		}
		date := start
		start = start.AddDate(0, 0, 1)
		return date
	}
}
