FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o moody-situation .

FROM scratch

COPY --from=builder /app/moody-situation /app/moody-situation
COPY conf.json .

EXPOSE 80
ENTRYPOINT ["/app/moody-situation"]