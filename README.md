## El Gourment Veloz - POC

This repository contains a POC implementation of an order processing system for a
theoretical restaurant. The system exposes a basic CRUD REST API.

You can find an OpenAPI spec for the service [here](docs/openapi.yml).

The service is implemented using [gin](https://gin-gonic.com/) for API routing
and [gorm](https://gorm.io/) for database storage. It follows an hexagonal
architecture approach with the following structure:

- *domain*: Contains logic related to order handling. It has no ties to the external
frameworks codebase (other than relying on gin's default validator)
- *internal/gin*: Code related to gin routing and endpoint handling
- *internal/gorm*: Code related to gorm data models and implementations of the data store
interfaces required in the domain's logic.

The current implementation requires a _postgres_ database to run. The easiest way to run
the server locally is using *Docker compose* with the following command:

```
$ docker compose up
```
>> This will run a postgres instance and expose it locally on port 5432, as well as build
an image with the server, run it and expose it locally on port 9001
 
If you wanna run the server locally, you can use compose to provide only the postgres
database with

```
$ docker compose start db
```

Then you can run the server with

```
$ go run main.go
```

There is a comprehensible set of unit tests in the project, written with ginkgo+gomega. To
run the tests, you can use one of the two commands:

```
$ go test ./...
```

If you have ginkgo installed locally, you can also use:

```
$ ginkgo -r
```

There are also a couple of basic load tests created with [k6](https://k6.io/). To run them,
you need to install k6 locally, start the server and run each test with:

```
$ k6 run k6/create.js
$ k6 run k6/index.js
```

### Covered use cases

- Create a new order
- Get the list of all orders
- Get the list of active orders
- Get a single order's details
- Update an order's status
- Update an order's list of dishes
- Store the history of statuses for a particular order
- VIP Prioritization with custom order sorting
- Cancel an order

### TODO

- Handle exact duplicates

