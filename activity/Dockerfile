FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o moody-activity .
RUN sh conf.sh
RUN apk add git && git clone https://github.com/vishnubob/wait-for-it && \ 
    [ -z DB_PORT ] && DB_PORT=5432

FROM scratch

COPY --from=builder /app/moody-activity /
COPY --from=builder /app/data/ /data
COPY --from=builder /app/wait-for-it/wait-for-it.sh /

EXPOSE 80
ENTRYPOINT ["./wait-for-it.sh", "-p", "$DB_PORT", "--", "./moody-activity"]