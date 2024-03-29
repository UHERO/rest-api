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
