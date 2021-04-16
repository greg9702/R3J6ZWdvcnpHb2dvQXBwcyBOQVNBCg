package models

import (
	"errors"
	"time"
)

type PicturesToBeFetched struct {
	StartDate time.Time `form:"start_date" binding:"required" time_format:"2006-01-02"`
	EndDate   time.Time `form:"end_date" binding:"required" time_format:"2006-01-02"`
}

func (p *PicturesToBeFetched) Validate() error {
	if p.EndDate.Before(p.StartDate) {
		return errors.New("End date is before start date")
	}
	return nil
}

type UrlList struct {
	Urls []string
}
