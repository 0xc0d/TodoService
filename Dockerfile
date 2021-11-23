##################
# Build Stage 1/2   -> Dependencies
##################
FROM golang:1.17 AS dependencies

WORKDIR $GOPATH/src/todoService

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/app

##################
# Build Stage 2/2   -> Production image
##################
FROM scratch AS production

COPY --from=dependencies /go/bin/app /app

ENTRYPOINT ["/app"]