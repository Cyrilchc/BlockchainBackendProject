# docker build -t backend .
# docker run -it --rm -p 5000:5000 --name backend backend

FROM golang:latest

WORKDIR /usr/src/Backend
EXPOSE 5000

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/Backend ./...

RUN go test
RUN go test -bench .

# Install Performance Testing tool
RUN wget https://github.com/grafana/k6/releases/download/v0.40.0/k6-v0.40.0-linux-amd64.deb
RUN dpkg -i k6-v0.40.0-linux-amd64.deb

# Run tests
#RUN k6 run sample_test.js
#RUN k6 run stress_test.js
#RUN k6 run spike_test.js

CMD ["Backend"]
