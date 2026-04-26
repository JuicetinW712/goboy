# Build
FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=js GOARCH=wasm go build -o goboy.wasm ./cmd/server

# Final Container
FROM nginx:1.29.8-alpine AS final

# RUN adduser -D base
# USER base

WORKDIR /app
COPY --from=builder /app/goboy.wasm /usr/share/nginx/html/goboy.wasm
COPY --from=builder /app/cmd/server/index.html /usr/share/nginx/html/index.html
COPY --from=builder /app/cmd/server/game.html /usr/share/nginx/html/game.html
COPY --from=builder /app/cmd/server/wasm_exec.js /usr/share/nginx/html/wasm_exec.js
COPY --from=builder /app/cmd/server/nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
