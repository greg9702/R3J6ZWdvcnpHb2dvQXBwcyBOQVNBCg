package executor

type GetRequestJobResult struct {
	Url   string
	Error error
}

type NasaApiCorrectResponse struct {
	Date           string `json:"date"`
	Explanation    string `json:"explanation"`
	MediaType      string `json:"media_type"`
	ServiceVersion string `json:"service_version"`
	Title          string `json:"title"`
	Url            string `json:"url"`
}

type NasaApiErrorResponse struct {
	Error NasaError `json:"error"`
}

type NasaError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
