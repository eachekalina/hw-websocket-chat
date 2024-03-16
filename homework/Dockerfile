FROM golang:1.21

WORKDIR /homework

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./server homework/cmd/server

EXPOSE 8080

CMD ["./server"]