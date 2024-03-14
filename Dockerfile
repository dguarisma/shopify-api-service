FROM amazonlinux:2.0.20211001.0

ARG GO_VERSION=1.19

RUN yum update -y && \
    yum install -y git tar gzip make gcc && \
    curl -sL https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz | tar -C /usr/local -xzf -

ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /go/src/app

COPY . .

RUN go mod download && \
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main

CMD ["./main"]
