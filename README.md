# intelligent-analysis-of-car-sensors-backend

This is a RESTful API service created in Go using Gin (a HTTP web framework). It is used to store vehicular sensor parameters collected through an OBD-II port of a vehicle stored in `csv` files.

## Getting started

### Technologies

- Go
- Gin (HTTP web framework for Go)
- Gorm (Go ORM)

### Deploy

`docker-compose` is used to simplify the deployment and avoid to install dependencies.

#### Requirements

- docker
- docker-compose

#### Run

From `root`:

```bash
docker-compose up -d
```
