FROM golang:stretch 

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o moody-server ./cmd/server/moody-server.go
WORKDIR /dist
RUN cp /build/moody-server .
RUN mkdir /root/.moody
RUN cp /build/config/* $HOME/.moody/


EXPOSE 1883
CMD ["/dist/moody-server"]