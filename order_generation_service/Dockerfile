FROM golang:latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download -x
COPY .env ./
COPY . .
RUN go build -o main .
EXPOSE 8081
CMD ["./main"]