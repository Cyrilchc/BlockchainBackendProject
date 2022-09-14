FROM golang:1.19.1-alpine

WORKDIR /usr/src/MockBlockchainService
EXPOSE 5001

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/MockBlockchainService ./...

CMD ["MockBlockchainService"]