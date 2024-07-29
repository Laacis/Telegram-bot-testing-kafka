FROM golang:latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download -x
COPY . .
RUN go build -o main .
EXPOSE 8082
CMD ["./main"]