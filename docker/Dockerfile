FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN mkdir bin

RUN go build -o bin/app src/app/main.go

EXPOSE 4242

CMD ["./bin/app"]
