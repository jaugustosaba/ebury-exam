# EBURY-EXAM

World tour helper CLI and REST applications. 

## Requirements

* linux kernel >= 6.8
* go >= 1.22
* docker >= 24.0.7
* curl >= 8.5.0

## Layout

| package | description |
| - | - |
| tour | tour core datatype abstraction and methods |
| tour/api | tour functionality service handling |
| tour/api/web | tour funcionality as a simple REST service |
| tour/cli | tour functionaly as a command line REPL |

## Local execution

Go's default run command could be used for running the command line or web server applications.

### CLI

Basic form:

```shell
go run tour/cli/main.go [INPUT-CSVs]
```

Example:

```shell
go run tour/cli/main.go data/input-routes.csv
```

### WEB

Basic form:

```shell
go run tour/api/web/main.go [-port NUMBER = 8080] [INPUT-CSVs]
```

Example:

```shell
go run tour/api/web/main.go data/input-routes.csv
```

## Docker build and run

Two containers available: a command line and a REST web server. Both containers contains a copy of data/input-routes.csv used as argument for sample tour loading. Both are alpine based (no MUSLC compatibility problems expected).

### CLI

```shell
docker build -f Dockerfile.cli . -t eburycli
docker run -it eburycli
```

### WEB

```shell
docker build -f Dockerfile.web . -t eburyweb
docker run -d -p 8080:8080 eburyweb
```

## REST API

Basic REST API description. Assuming local execution at port 8080. 

OBS: The content type header must be set to ***application/json***.

### Response Common Data

All responses share a commom outer object with the following fields:

| field   | type | description |
| ------  | ----- | ----------- |
| status  | string | "ok" or "error" |
| reason  | string | a message describing the error |
| response | object | the response payload (abscent on errors) |

### POST /route/add 

Adds a new route to the current tour.

#### Request format

| field   | type | description |
| ------  | ----- | ----------- |
| origin  | string | origin city name |
| destiny | string | destiny city name |
| cost    | int    | cost of this route |

#### Response format

An empty object is returned in "response" to inform success.

#### Example using CURL

Request:

```json
{
    "origin": "A",
    "destiny": "B",
    "cost": 100
}
```


Example:

```shell
curl -v localhost:8080/route/add -H 'Content-type: application/json' --data '{"origin"
    : "A", "destiny": "B", "cost": 100 }'
```

Response:

```json
{
    "status":"ok",
    "response":{}
}
```

### POST /route/shortest

Computes shortest route bewteen two cities

#### Request format

| field   | type | description |
| ------  | ----- | ----------- |
| origin  | string | origin city name |
| destiny | string | destiny city name |

#### Response format

| field   | type | description |
| ------  | ----- | ----------- |
| shortestRoute  | string[] | array with each city name on route |
| cost | int | shortest route cost value |


#### Example using CURL

Request:

```json
{
    "origin": "A",
    "destiny": "B"
}
```

Example:

```shell
curl -v localhost:8080/route/shortest -H 'Content-type: application/json' --data '{"origin": "GRU", "destiny": "CDG"}'
```

Response:

```json
{
    "status":"ok",
    "response":{
        "shortestRoute":["GRU","BRC","SCL","ORL","CDG"],
        "cost":40
    }
}
```

## Internals

The tour datatype is a graph implementation where cities are nodes and routes are edges. The shortest path finder algorithm used is the classic ***Djikstra Algorithm***.

A mutex is used on web server to avoid concurrent access to out in memory tour datatype.
