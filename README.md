# MockApiHub

[![Go Report Card](https://goreportcard.com/badge/github.com/wcsanders1/MockApiHub)](https://goreportcard.com/report/github.com/wcsanders1/MockApiHub)
[![BuildStatus](https://travis-ci.org/wcsanders1/MockApiHub.svg?branch=master)](https://travis-ci.org/wcsanders1/MockApiHub)
[![codecov](https://codecov.io/gh/wcsanders1/MockApiHub/branch/master/graph/badge.svg)](https://codecov.io/gh/wcsanders1/MockApiHub)

This server allows you to create a collection of mock APIs.

This is a work in progress and a release version is not yet complete. However, it is minimally functional, and you may build and use it. To build this application, you must have a Go deveopment environment set up. To set up a Go development environment, follow the instructions here: golang.org/doc/install.

After creating a Go development environment, follow these steps to build and run the MockApiHub:

* Clone this repository into your `GOPATH`.
* Install `govendor` by running `go get -u github.com/kardianos/govendor`
* From the root of the MockApiHub repo, run `govendor sync`, which should populate the `vendor` file with the packages needed to build the MockApiHub application.
* Run `go build mockApiHub.go`, which should produce an executable file that you can run.

After you have the MockApiHub running, you should be able to send a request to the endpoints of the example mock APIs that are part of the repository. Their files are contained in `api/apis`. For example, if you send a request from the same machine on which that the MockApiHub is running to `http://localhost:5001/customersapi/accounts`, you should receive a response with the following JSON:

``` json
[
    {
        "customerId": 1,
        "balance": 4.5
    },
    {
        "customerId": 2,
        "balance": 0
    }
]
```

Note that the JSON returned is the JSON in the `accounts.json` file in `api/apis/customersApi`. Also, note that the endpoints and their URLs are configured in `customersApi.toml` in that directory.

Adding mock APIs is easy. To do so, follow these steps:

* Add a directory to `api/apis`. **The directory must end with the letters `api`.**
* Add a configuration `toml` file, also ending with the letters `api`.
* To configure the API, follow the examples in the example config files in the repo.
* Note: If you want a route to have a parameter, add `:` to the route fragment. For example, the `getCustomers` endpoint in the `customersApi` example is configured as follows: `customers/:id/balances`. This means that the endpoint could be hit with the following URL: `http://localhost:5001/customersapi/customers/12345/balances`, where `12345` is the `id`.

If you change or add a mock API, make the change to `api/apis`, then send a `POST` request to the main server, configured by default to listen on port 5000, with the following path: `refresh-all-mock-apis`.

To see all registered mock APIs, send this `GET` request to the main server: `show-all-registered-mock-apis`.