FROM golang:stretch 

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o moody-gateway main.go
WORKDIR /dist
RUN cp /build/moody-gateway .

CMD ["/dist/moody-gateway"]