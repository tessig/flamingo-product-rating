# flamingo-product-rating [![Build Status](https://travis-ci.org/tessig/flamingo-product-rating.svg?branch=master)](https://travis-ci.org/tessig/flamingo-product-rating)
A simple [Flamingo](https://www.flamingo.me/) application as showcase for a DevOps approach.

## Building the app

Needs at least go 1.11.4 and uses [go modules](https://github.com/golang/go/wiki/Modules).

```bash
go build -o rating .
```

## Development setup

Database and productservice for development can be found in devenv directory.

Simply run `docker-compose up` from within devenv.

On first or new setup, the database will be empty. Run 

* `CONTEXT=dev go run main.go migrate up` to create the schema
* `CONTEXT=dev go run main.go seed` to import some test data

Then start the app via `CONTEXT=dev go run main.go serve`

The config in `config/config_dev.yml` matches the docker-compose setup.

The app will be under http://localhost:3322/
The metrics endpoint will be under http://localhost:13210/metrics

## Standalone example

In docker-compose-standalone directory you can find a complete operational example. 
