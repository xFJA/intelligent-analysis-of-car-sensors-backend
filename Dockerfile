FROM golang:1.14.2-stretch

ENV PORT=8080
ENV GO111MODULE=on

WORKDIR /build 

# Download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy source code
COPY . .

# Create binary
RUN go build -o main .
WORKDIR /dist
RUN cp /build/main .

EXPOSE 8080

CMD ["/dist/main"]