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
interfaces required in the domain's logic

To keep running the service simple, _sqlite_ is currently the default database.
To run the server, simply clone the project locally and run:

```
$ go run main.go
```

This should create an sqlite db file, run the migrations and start the server on port *9001*.

### Covered use cases

- Create a new order
- Get the list of all orders
- Get the list of active orders
- Get a single order's details
- Update an order's status
- Store the history of statuses for a particular order
- Cancel an order

