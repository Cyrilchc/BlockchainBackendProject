# docker build -t mockblockchainservice .
# docker run -it --rm -p 5001:5001 --name mockblockchainservice mockblockchainservice

FROM golang:latest

WORKDIR /usr/src/MockBlockchainService
EXPOSE 5001

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/MockBlockchainService ./...

CMD ["MockBlockchainService"]