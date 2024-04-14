# yamr

Yet another meter recorder

## Service

### REST routes

* `GET /kinds`: returns a list kinds available in the recorder
* `GET /readings`: returns a list of all recordings
* `GET /readings?kind=<kind>`: returns a list of all recordings of kind `<kind>`
* `POST /reading`: add a new reading

### Development

#### Requirements

* go (for local building and/or running)
* Docker (for containerized building and/or running)

#### Local

From `service/` start the app with

```bash
go run .
```

Build the executable with

```bash
go build
```

### Docker

Build the image

```bash
docker build -f Dockerfile.service --tag cybernop/yamr-service .
```
