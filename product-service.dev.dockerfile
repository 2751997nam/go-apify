
FROM golang:alpine

WORKDIR /app

COPY . /app


RUN go mod download && go mod verify

RUN go install github.com/cosmtrek/air@latest

# CMD ["air"]