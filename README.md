[![Build Status](https://travis-ci.org/UHERO/rest-api.svg?branch=master)](https://travis-ci.org/UHERO/rest-api)
# rest-api
REST API for UHERO Time Series

## Mock API

In the `mock-api` folder, run `npm install`. Then start the mock-api server with

```
$ NODE_ENV=development node index.js
```

This will start up the mock-api with stubs at http://localhost:8080

You can look at the API docs at http://localhost:8080/docs

When looking at examples in the documentation, you will need to replace `api.uhero.hawaii.edu` with `localhost:8080` to get valid results.


## Environment Variables

* GITHUB_KEY
* GITHUB_SECRET
* DB_USER
* DB_PASSWORD

### Session Secret
* SESSION_SECRET example:
```
export SESSION_SECRET=`openssl rand -base64 32`
```

## Config
`common/config.json` allows you to change the database connection string.

## Key files to put in the `key` folder
```
openssl genrsa -out app.rsa 1024
openssl rsa -in app.rsa -pubout > app.rsa.pub
```
