package executor

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

const (
	nasaUrl = "https://api.nasa.gov/planetary/apod"
)

type JobInterface interface {
	Execute()
}

type FetchNasaApiJob struct {
	requestParams string
	result        chan GetRequestJobResult
	url           string
}

func NewFetchNasaApiJob(requestParams string, result chan GetRequestJobResult) JobInterface {
	g := FetchNasaApiJob{requestParams, result, nasaUrl}
	return &g
}

func (f *FetchNasaApiJob) Execute() {

	res := GetRequestJobResult{"", nil}

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodGet, f.url+f.requestParams, nil)
	request.Header.Set("Content-Type", "application/json")

	if err != nil {
		res.Error = err
		f.result <- res
		return
	}

	// perform call to external API here
	response, err := client.Do(request)

	if err != nil {
		res.Error = err
		f.result <- res
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		res.Error = err
		f.result <- res
		return
	}

	if response.StatusCode != http.StatusOK {
		responseBody := NasaApiErrorResponse{}
		err = json.Unmarshal(body, &responseBody)

		if err != nil {
			res.Error = err
			f.result <- res
			return
		}

		res.Error = errors.New("NasaAPI error: " + responseBody.Error.Message)
		f.result <- res
		return
	}

	responseBody := NasaApiCorrectResponse{}
	err = json.Unmarshal(body, &responseBody)

	if err != nil {
		res.Error = err
		f.result <- res
		return
	}

	res.Url = responseBody.Url

	f.result <- res
}
