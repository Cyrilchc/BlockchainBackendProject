# docker build -t backend .
# docker run -it --rm -p 5000:5000 --name backend backend

FROM golang:latest

WORKDIR /usr/src/Backend
EXPOSE 5000

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/Backend ./...

CMD ["Backend"]