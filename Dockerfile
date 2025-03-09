FROM golang:1.23

WORKDIR /app

COPY ./server/go.mod .
COPY ./server/go.sum .

RUN go mod download

COPY ./server .

RUN CGO_ENABLED=true go build -o /app/app

EXPOSE 8080

CMD ["/app/app"]
