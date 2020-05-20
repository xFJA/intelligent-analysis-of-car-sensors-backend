FROM golang:1.14.2-stretch

ENV PORT=8080
ENV GO111MODULE=on

WORKDIR /app 

# Download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy source code
COPY . .

# Create binary
RUN go build -o main .

EXPOSE 8080

CMD ["main"]