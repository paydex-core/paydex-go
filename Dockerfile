FROM golang:1.13
WORKDIR /go/src/github.com/paydex-core/paydex-go

COPY . .
ENV GO111MODULE=on
RUN go install github.com/paydex-core/paydex-go/tools/...
RUN go install github.com/paydex-core/paydex-go/services/...
