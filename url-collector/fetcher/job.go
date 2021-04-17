package fetcher

import "time"

type JobInterface interface {
	Execute()
}

type GetRequestJobResult struct {
	Url   string
	Error error
}

type FetchNasaApiJob struct {
	requestParams string
	result        chan GetRequestJobResult
}

func NewFetchNasaApiJob(requestParams string, result chan GetRequestJobResult) JobInterface {
	g := FetchNasaApiJob{requestParams, result}
	return &g
}

func (g *FetchNasaApiJob) Execute() {

	time.Sleep(time.Second)

	// perform call to external API here

	timeStr := "www.example.com/"
	timeStr += g.requestParams

	res := GetRequestJobResult{timeStr, nil}
	g.result <- res
}
