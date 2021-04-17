package fetcher

import (
	"fmt"
	"os"
	"time"
	"url-collector/url-collector/executor"
	"url-collector/url-collector/models"
	"url-collector/url-collector/utils"
)

const (
	defaultApiKey = "DEMO_KEY"
)

type FetcherInterface interface {
	FetchData(interface{}) ([]string, error)
}

// nasaFetcher is the implementation of FetcherInterface
type nasaFetcher struct {
	exec executor.ExecutorInterface
}

func NewNasaFetcher(exec executor.ExecutorInterface) FetcherInterface {
	nf := &nasaFetcher{exec}
	return nf
}

func (nf *nasaFetcher) FetchData(object interface{}) ([]string, error) {
	returnList := make([]string, 0)

	picturesDateRange := object.(*models.PicturesToBeFetched)

	requestsParams, err := nf.prepareRequestArguments(picturesDateRange.StartDate, picturesDateRange.EndDate)

	if err != nil {
		return returnList, err
	}

	results := make(chan executor.GetRequestJobResult, 100)

	go func() {
		for i := 0; i < len(requestsParams); i++ {
			job := executor.NewFetchNasaApiJob(requestsParams[i], results)
			nf.exec.AddNewJob(job)
		}
	}()

	// task collector
	for i := 0; i < len(requestsParams); i++ {
		jobResult := <-results
		if jobResult.Error != nil {
			// do we want to handle it in this way?
			fmt.Println("Executor returned error")
			return returnList, jobResult.Error
		} else {
			returnList = append(returnList, jobResult.Url)
		}
	}
	return returnList, nil
}

// prepareRequestArguments is very flexible, if we want to add more params to request, we can do it here
func (nf *nasaFetcher) prepareRequestArguments(startDate time.Time, endDate time.Time) ([]string, error) {
	var resultParamList []string
	dateList, err := utils.GetListOfDate(startDate, endDate)

	if err != nil {
		return resultParamList, err
	}

	apiKey := defaultApiKey

	apiKeyString := os.Getenv("API_KEY")
	if apiKeyString != "" {
		apiKey = apiKeyString
	}

	for i := 0; i < len(dateList); i++ {
		resultParamList = append(resultParamList, fmt.Sprintf("?api_key=%s&date=%s", apiKey, dateList[i]))
	}

	return resultParamList, nil
}
