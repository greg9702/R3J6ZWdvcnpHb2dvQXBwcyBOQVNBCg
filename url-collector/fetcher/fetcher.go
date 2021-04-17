package fetcher

import (
	"fmt"
	"time"
	"url-collector/url-collector/models"
	"url-collector/url-collector/utils"
)

type FetcherInterface interface {
	FetchData(interface{}) ([]string, error)
}

// nasaFetcher is the implementation of FetcherInterface
type nasaFetcher struct {
	executor ExecutorInterface
}

func NewNasaFetcher(executor ExecutorInterface) FetcherInterface {
	nf := &nasaFetcher{executor}
	return nf
}

func (nf *nasaFetcher) FetchData(object interface{}) ([]string, error) {
	var returnList []string

	dateRange := object.(*models.PicturesToBeFetched)

	requestsParams, err := nf.prepareRequestArguments(dateRange.StartDate, dateRange.EndDate)

	if err != nil {
		return returnList, err
	}

	results := make(chan GetRequestJobResult, 100)

	// spawn task generator (date, channel)
	go func() {
		for i := 0; i < len(requestsParams); i++ {
			job := NewFetchNasaApiJob(requestsParams[i], results)
			nf.executor.AddNewJob(job)
		}
	}()

	// task collector
	for i := 0; i < len(requestsParams); i++ {
		jobResult := <-results
		if jobResult.Error != nil {
			fmt.Println("executor returned error")
		} else {
			returnList = append(returnList, jobResult.Url)
		}
	}
	return returnList, nil
}

// prepareRequestArguments is very flexible, if we want to add more params to request, we can do it here
func (nf *nasaFetcher) prepareRequestArguments(startDate time.Time, endDate time.Time) ([]string, error) {
	dateList := utils.GetListOfDate(startDate, endDate)
	return dateList, nil
}
