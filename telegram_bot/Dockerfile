# Getting latest
FROM golang:latest

# Woking dir inside docker /app
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the .env file and rest of files
COPY .env ./
COPY . .

# Build the bot that is bot.go
RUN go build -o bot .

# Command to run the executable that is bot
CMD ["./bot"]
