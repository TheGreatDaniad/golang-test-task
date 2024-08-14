FROM golang:1.18-alpine

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o queryservice .

EXPOSE 8081

CMD ["./queryservice"]
