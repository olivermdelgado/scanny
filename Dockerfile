FROM golang:1.15
ADD . /go/src/redditScanner
WORKDIR /go/src/redditScanner
RUN go build -o reddit-scanner main.go
CMD ["reddit-scanner"]