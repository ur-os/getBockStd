FROM golang:1.22
LABEL authors="urick0s"

WORKDIR /app

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .

RUN go mod tidy
RUN go build -o app/getBlock

CMD ["./app/getBlock"]