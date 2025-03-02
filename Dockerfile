FROM golang:1.22.2

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY main.go ./
COPY ./templates ./templates

RUN CGO_ENABLED=0 GOOS=linux go build -o ./weather-server

CMD ["./weather-server"]

