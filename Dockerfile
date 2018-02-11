FROM golang:1.9 as builder

WORKDIR /go/src/github.com/radu-matei/events-operator

COPY . .

RUN go get -u github.com/golang/dep/...
RUN dep ensure

RUN go build


FROM ubuntu

COPY --from=builder /go/src/github.com/radu-matei/events-operator/events-operator .

CMD ["./events-operator"]