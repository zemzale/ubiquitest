FROM golang:1.23-alpine

WORKDIR /app

COPY ./server/go.mod .
COPY ./server/go.sum .

RUN go mod download

COPY ./server .

RUN go build -o /app/app

EXPOSE 8080

CMD ["/app/app"]
