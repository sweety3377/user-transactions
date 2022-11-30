# syntax=docker/dockerfile=1

FROM golang:1.19-alpine
WORKDIR /user-transactions
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o ./app ./cmd/app
RUN echo "user-transaction service started"
EXPOSE 8080
CMD ["/user-transactions/app"]