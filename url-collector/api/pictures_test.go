package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type ExpectedOkResponse struct {
	Urls []string `json:"urls"`
}

type ExpectedErrorResponse struct {
	Error string `json:"error"`
}

type FetcherMock struct {
	mock.Mock
}

func (fm *FetcherMock) FetchData(object interface{}) ([]string, error) {
	args := fm.Called(object)
	return args.Get(0).([]string), args.Error(1)
}

func TestGetImages(t *testing.T) {

	todayDate := time.Now()
	dateLaterThanToday := todayDate.AddDate(0, 0, 5)
	dateLaterThanTodayString := dateLaterThanToday.Format("2006-01-02")

	cases := []struct {
		requestParams      map[string]string
		cName              string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			requestParams:      map[string]string{},
			cName:              "No query params",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"Key: 'PicturesToBeFetched.StartDate' Error:Field validation for 'StartDate' failed on the 'required' tag\nKey: 'PicturesToBeFetched.EndDate' Error:Field validation for 'EndDate' failed on the 'required' tag"}`,
		},
		{
			requestParams:      map[string]string{"start_date": "2020-07-02"},
			cName:              "Missing end_date param",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"Key: 'PicturesToBeFetched.EndDate' Error:Field validation for 'EndDate' failed on the 'required' tag"}`,
		},
		{
			requestParams:      map[string]string{"end_date": "2020-07-02"},
			cName:              "Missing start_date param",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"Key: 'PicturesToBeFetched.StartDate' Error:Field validation for 'StartDate' failed on the 'required' tag"}`,
		},
		{
			requestParams:      map[string]string{"start_date": "123", "end_date": "2020-07-02"},
			cName:              "Invalid start_date format",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"parsing time \"123\" as \"2006-01-02\": cannot parse \"123\" as \"2006\""}`,
		},
		{
			requestParams:      map[string]string{"start_date": "2020-07-02", "end_date": "123"},
			cName:              "Invalid end_date format",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"parsing time \"123\" as \"2006-01-02\": cannot parse \"123\" as \"2006\""}`,
		},
		{
			requestParams:      map[string]string{"start_date": "2020-07-02", "end_date": "2020-07-01"},
			cName:              "end_date before start_date",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"End date is before start date"}`,
		},
		{
			requestParams:      map[string]string{"start_date": "2020-07-02", "end_date": dateLaterThanTodayString},
			cName:              "end_date after today",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"End date is after today's date"}`,
		},
		{
			requestParams:      map[string]string{"start_date": dateLaterThanTodayString, "end_date": dateLaterThanTodayString},
			cName:              "start_date older than today",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"End date is after today's date"}`,
		},
		{
			requestParams:      map[string]string{"start_date": "2020-07-02", "end_date": "2020-07-03"},
			cName:              "Any fetcher error",
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"error":"Fetcher error"}`,
		},
		{
			requestParams:      map[string]string{"start_date": "2020-07-02", "end_date": "2020-07-03"},
			cName:              "Ok request",
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"urls":["www.testurl.com/1","www.testurl.com/2"]}`,
		},
		{
			requestParams:      map[string]string{"start_date": "2020-07-02", "end_date": "2020-07-03"},
			cName:              "Ok request, empty urls list",
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"urls":[]}`,
		},
	}

	url := "/pictures"

	emptyStringList := make([]string, 0)

	var nonEmptyStringList []string
	nonEmptyStringList = append(nonEmptyStringList, "www.testurl.com/1")
	nonEmptyStringList = append(nonEmptyStringList, "www.testurl.com/2")

	mockFetcher := new(FetcherMock)
	mockFetcher.On("FetchData", mock.Anything).Return(emptyStringList, errors.New("Fetcher error")).Once().
		On("FetchData", mock.Anything).Return(nonEmptyStringList, nil).Once().
		On("FetchData", mock.Anything).Return(emptyStringList, nil).Once()

	controller := NewPicturesController(mockFetcher)

	// create test router
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(
		gin.Recovery(),
	)
	router.Use(cors.Default())

	router.GET(url, controller.GetImages)

	as := assert.New(t)
	for idx, tc := range cases {

		request, err := http.NewRequest(http.MethodGet, url, nil)
		q := request.URL.Query()

		for k, v := range tc.requestParams {
			q.Add(k, v)
		}

		request.URL.RawQuery = q.Encode()

		request.Header.Set("Content-Type", "application/json")
		require.NoError(t, err)

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)
		as.Equal(tc.expectedStatusCode, recorder.Code, "Expected test case [%d] %s: to return %d; got %d instead", idx, tc.cName, tc.expectedStatusCode, recorder.Code)
		as.Equal(tc.expectedBody, recorder.Body.String(), "Expected test case [%d] %s: to return %s; got %s instead", idx, tc.cName, tc.expectedBody, recorder.Body.String())
	}
}
