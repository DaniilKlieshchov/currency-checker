FROM golang:latest

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

EXPOSE 1313

RUN go build -o ./currency-checker ./cmd

CMD ["./currency-checker"]

