FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

COPY receipt.json /app/receipt.json

RUN go build -o receipt-processor .

ENV PORT=5000

EXPOSE $PORT

CMD ["./receipt-processor"]

