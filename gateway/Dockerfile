FROM golang:1.15-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o moody-gateway .

FROM scratch

ENV GODEBUG="x509ignoreCN=0"


COPY --from=builder /app/moody-gateway /

EXPOSE 7000
ENTRYPOINT ["/moody-gateway"]