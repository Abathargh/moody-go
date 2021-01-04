FROM node:12-alpine AS builder
WORKDIR /app
COPY . .
RUN arch=$(uname -m); \
        if [ "$arch" = "aarch64" ] || [ "$arch" = "armv7l" ]; then \
            apk update && apk add python make g++; \
        fi
RUN npm install && npm run build

FROM node:12-alpine
RUN npm install -g serve
WORKDIR /app
COPY --from=builder /app/build .
EXPOSE 3000

CMD ["serve", "-l", "tcp://0.0.0.0:3000", "-s", "."]
