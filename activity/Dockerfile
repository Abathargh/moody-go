FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o moody-activity .
RUN sh conf.sh


FROM scratch

COPY --from=builder /app/data/conf.json /data/
COPY --from=builder /app/moody-activity /

EXPOSE 80
ENTRYPOINT ["/moody-activity"]