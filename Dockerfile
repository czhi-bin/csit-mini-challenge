FROM golang:1.20-alpine

WORKDIR /app

ENV GIN_MODE=release

COPY go.mod .
RUN go mod download

COPY . .

EXPOSE 8080

RUN go build -o /csit-mini-challenge

CMD [ "/csit-mini-challenge" ]