# MockApiHub

[![Go Report Card](https://goreportcard.com/badge/github.com/wcsanders1/MockApiHub)](https://goreportcard.com/report/github.com/wcsanders1/MockApiHub)
[![BuildStatus](https://travis-ci.org/wcsanders1/MockApiHub.svg?branch=master)](https://travis-ci.org/wcsanders1/MockApiHub)
[![codecov](https://codecov.io/gh/wcsanders1/MockApiHub/branch/master/graph/badge.svg)](https://codecov.io/gh/wcsanders1/MockApiHub)
[![License](https://img.shields.io/badge/license-mit-blue.svg)](https:/githubusercontent.com/wcsanders1/MOckApiHub/master/LICENSE)
[![GoDoc](https://img.shields.io/badge/go-documentation-darkblue.svg)](https://godoc.org/github.com/wcsanders1/MockApiHub)
![Version](https://img.shields.io/badge/version-0.1.0-darkred.svg)

## Overview

This application allows you to create a collection of mock APIs. The APIs are easily configurable and you can reconfigure them--i.e., change their URLs, HTTP methods, or return data--without restarting the hub server. This application is useful when you need to test a service that makes HTTP requests and you want to control the data being returned from those requests.

## Getting Started

Download the appropriate binary from [here](https://github.com/wcsanders1/MockApiHub/releases). The binary with the `exe` extension is for Windows and the other is for Linux.

Before starting the application, you must configure the hub server in a file named `app_config.toml`, which must be in the same directory level as the executable. The file must be in `toml` format. You must provide the port on which you want the hub server to listen. If you want the server to use TLS, you must enable TLS and provide the paths for the certificate and key files. In addition, you can configure logging in this file. (Note that logging configuration for the mock APIs is separate, discussed below.) Here is an example of a valid `app_config.toml`:

```toml
[http]
port = 5000
useTLS = false
certFile = ""
keyFile = ""

[log]
loggingEnabled = true
fileName = "testLogs/mockApiHub/default.log"
maxFileDaysAge = 3
formatAsJSON = true
prettyJSON = true
```

After configuring the hub server, you need to configure your mock APIs and provide files containing the data you want them to return. The API configuration files and data files must be placed in `api/apis`, whose root must be the directory of the executable. Each mock API must have its own directory as a subdirectory of `api/apis`, which must itself end in the letters `Api`. Each mock API must have it's own configuration file, which must end in the letters `api` and must be in `toml` format. (I apologize for the complexity and constraint of how the configuration files must be. I plan to simplify this in a later release.)

This application does not limit the amount of mock APIs that can run at once; however, each mock API must listen on a distinct port, and none of them can listen on the same port as the hub server.

Here is an example of a valid mock API configuration file, named `customersApi.toml`, which has two endpoints:

```toml
baseUrl = "customersApi"

[log]
loggingEnabled = true
fileName = "testLogs/customersApi/default.log"
maxFileDaysAge = 3
formatAsJSON = true
prettyJSON = true

[http]
port = 5001
useTLS = false
certFile = ""
keyFile = ""

[endpoints]

    [endpoints.getAllAccounts]
    path = "accounts"
    file = "accounts.json"
    method = "GET"

    [endpoints.getCustomerBalances]
    path = "customers/:id/balances"
    file = "customers.json"
    method = "GET"
```

As shown above, logging for a mock API is configured the same way as for the hub server. If you want a mock API to use TLS and the hub server is also using TLS, you can leave the certificate and key file entries in the mock API configuration empty, in which case the certificate and key of the hub server will be used. Files containing the data that a mock API returns need to be placed in the same directory as the mock API's configuration file.

A mock API can have a `baseURL`, which applies to all of its endpoints. Each endpoint of a mock API must have an entry in the `endpoints` section of the configuration file. Files containing the data that a mock API returns need to be placed in the same directory as the mock API's configuration file. In the example above (assuming the hub server is running on `localhost`) the following `GET` request will return whatever is in the `accounts.json` file: `http://localhost:5001/customersapi/accounts`. If you want to enforce valid JSON for a particular endpoint, you can add `enforceValidJSON = true` to that endpoint's configuration.

If you want to have a variable as part of a mock API's URL, just put a colon in front of the route fragment. Also, a query string can be added to any HTTP request to a mock API and its values and keys will be logged. For example, the `getCustomerBalances` endpoint in the configuration above can be hit with the following URL, `http://localhost:5001/customersapi/customers/12345/balances?page=2&size=50`, which will return whatever is in `customers.json`. If logging is enabled, the request will be logged like this (note the logging of the `id` variable in `params`, as well as the keys and values in the query string):

```json
{
  "baseURL": "customersApi",
  "certFile": "",
  "func": "github.com/wcsanders1/MockApiHub/api.(*API).ServeHTTP",
  "keyFile": "",
  "level": "debug",
  "msg": "handler exists for this path",
  "params": {
    "id": "12345"
  },
  "path": "customersapi/customers/:id/balances",
  "pkg": "api",
  "port": 5001,
  "query": {
    "page": [
      "2"
    ],
    "size": [
      "50"
    ]
  },
  "time": "2018-11-03 20:40:36",
  "useTLS": false
}
```

This application does not cache the contents of the files that the mock APIs serve, so if you want to change the content of the files, you can do so without restarting or reloading anything.

If you make a change to a mock API's configuration file, simply save the changes and send a `POST` request to the hub server with the path `refresh-all-mock-apis`; e.g., `http://localhost:5000/refresh-all-mock-apis`. which will apply any changes you made to the mock APIs.

A `GET` request to the hub server with the path `show-all-registered-mock-apis` will return all of the registered mock APIs with their configurations; e.g., `http://localhost:5000/show-all-registered-mock-apis`.

## License

[MIT](https://github.com/wcsanders1/MOckApiHub/master/LICENSE)