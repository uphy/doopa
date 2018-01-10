FROM golang:1.9

COPY . /go/src/github.com/uphy/doopa/
CMD go run /go/src/github.com/uphy/doopa/app/web/main.go
