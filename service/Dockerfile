FROM golang

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o moody-service .

EXPOSE 80
CMD ["/app/moody-service"]