# syntax=docker/dockerfile:1

FROM golang:latest

WORKDIR /app
VOLUME ./tablebases

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .

RUN go build -o ./fathom-http

EXPOSE 80

ENTRYPOINT ["/app/fathom-http", "--listen", ":80", "--poolSize", "0", "--tbDir", "./tablebases"]
CMD ["--maxTime", "60s", "--allowOrigin", "*"]
