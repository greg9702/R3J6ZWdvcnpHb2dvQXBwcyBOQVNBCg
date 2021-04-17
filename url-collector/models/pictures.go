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
	// we do not let inserting dates higher than today, Nasa API doesnt like it
	if p.EndDate.After(time.Now()) {
		return errors.New("End date is after today's date")
	}

	return nil
}

type UrlList struct {
	Urls []string
}
