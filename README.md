[![Build Status](https://travis-ci.org/UHERO/rest-api.svg?branch=master)](https://travis-ci.org/UHERO/rest-api)
# rest-api
REST API for UHERO Time Series

## Mock API

In the `mock-api` folder, run `npm install`. Then start the mock-api server with

```
npm run dev
```

This will start up the mock-api with stubs at http://localhost:8080

You can look at the API docs at http://localhost:8080/docs

When looking at examples in the documentation, you will need to replace `api.UHERO.hawaii.edu` with `localhost:8080` to get valid results.

## Generating Documentation

```
npm install -g bootprint
npm install -g bootprint-swagger
bootprint swagger mock-api/api/swagger.yaml public/docs
```


## Environment Variables

* GITHUB_KEY
* GITHUB_SECRET
* GITHUB_CALLBACK
* DB_USER
* DB_PASSWORD
* DB_HOST
* SESSION_SECRET
* JIRA_USER
* JIRA_PASSWORD

### Session Secret
* SESSION_SECRET example:
```
export SESSION_SECRET=`openssl rand -base64 32`
```

## Key files to put in the `key` folder
```
openssl genrsa -out app.rsa 1024
openssl rsa -in app.rsa -pubout > app.rsa.pub
```

## Testing the API

End-to-end testing of the API can be accomplished using [Postman](https://www.getpostman.com/ "Postman").
Test collections (JSON files) can be
found in the `tests` directory. In Postman, click the *Import* button in the top left, and drop
or select one of the test collections to upload. In the left sidebar area, select tab *Collections*.
Choose a test collection to run, and click the expansion arrow. You should see a blue *Run* button.
Clicking this will open Postman's Collection Runner. Select the collection to run, if it's not already
highlighted, and click the blue *Start Test* button.

## Testing the UI

The API developer portal requires polymer to build the assets in the `public` folder.

Assuming you already have npm and bower installed, to install the Polymer CLI run the following command:
```
npm install -g polymer-cli
```

You can then run the following command from the root of the project to build the UI, build the server, and start the server:
```
cd public && polymer build && cd .. && go run main.go
```

## API Usage Metrics

The API includes built-in metrics tracking to monitor usage patterns and popular data series without requiring external dependencies.

### Viewing Metrics

**Current Day Metrics:**
```
GET /v1/metrics
```

**Historical Metrics (up to 30 days):**
```
GET /v1/metrics/historical?days=7
```

*Note: Metrics endpoints require valid API authentication.*

### Metrics Collected

- **Endpoint Usage**: Request counts, response times, and error rates by endpoint
- **Series Popularity**: Most requested series by ID and name
- **API Key Activity**: Usage patterns per API key (anonymized for privacy)
- **Temporal Patterns**: Hourly request distribution
- **Performance**: Average response times and error rates

### Metrics Storage

- Metrics are stored in JSON files in the `./metrics/` directory (configurable via `METRICS_DIR` environment variable)
- Daily files are automatically created and rotated at midnight
- In-memory counters provide real-time tracking with minimal performance impact

### Sample Output

```json
{
  "period": "2025-01-24",
  "summary": {
    "total_requests": 1340,
    "unique_api_keys": 12,
    "avg_response_time_ms": 145.32,
    "error_rate": 0.0149
  },
  "top_endpoints": [
    {"path": "/v1/series", "requests": 780},
    {"path": "/v1/series/siblings", "requests": 420}
  ],
  "top_series": [
    {"identifier": "name:GDP@HI", "requests": 67},
    {"identifier": "id:12345", "requests": 45}
  ],
  "hourly_distribution": [12, 8, 15, 42, ...]
}
```
