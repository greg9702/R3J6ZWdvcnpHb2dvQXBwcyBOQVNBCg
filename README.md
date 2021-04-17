### __Overview__
For the project used [Gin Web Framework](https://github.com/gin-gonic/gin). <br>
In tests used [Testify](https://github.com/stretchr/testify) for creating mocks.

### __Project structure__
```
├── docker-compose.yaml   
├── README.md
└── [url-collector]
    ├── [api] - keeps the pictures API handler
    ├── Dockerfile
    ├── [executor] - provides a executor for executing tasks concurrently
    ├── [fetcher] - responsible for coordinating a process of fetching external data
    ├── go.mod
    ├── go.sum
    ├── main.go - main file, entry point
    ├── [models] - contains models required by the application
    └── [utils] - contains util tools and helpers
```
### __Run__

To start the project:
```
docker-compose up --build
```

### __Tests__
Prepared tests for `GET /pictures` handler. <br>
To run tests:
```
cd url-collector
go test ./... url-collector/url-collector
```

### __Input validation__
Both parameters `start_date`, `end_date` are required. <br>
Both must be in `2006-01-02` format and cannot be further than the current date. <br>
Also `start_date` must be earlier or equal to `end_date`.

### __Endpoints description__
#### __Get images__

> GET /pictures HTTP/1.1

Get images from the time interval.

##### Request Headers
- Content-Type: application/json.

##### Status Codes
- 200 - all images were fetched successfully
- 400 - invalid input data or missing Content-Type: application/json header
- 500 - server error

#### Query Parameters
- start_date - date in format `2006-01-02`
- end_date - date in format `2006-01-02`

##### __Example Request__
```
curl -X GET "http://localhost:8080/pictures?start_date=2021-04-14&end_date=2021-04-15" -H "Content-Type: application/json"
```

##### __Example Response__
```
{
   "urls":[
      "https://apod.nasa.gov/apod/image/example/1.png",
      "https://apod.nasa.gov/apod/image/example/2.png"
   ]
}
```

##### __Example Error Response__
```
{
   "error":"..."
}
```
