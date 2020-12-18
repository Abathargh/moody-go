FROM golang:alpine AS builder
# TODO pick max go 1.15 OR implement altSubjects in certs

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