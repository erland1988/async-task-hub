FROM node:20-alpine AS vue-builder

WORKDIR /app

COPY . .

RUN npm config set registry https://registry.npmmirror.com

ENV NODE_OPTIONS="--max-old-space-size=4096"

RUN cd vue/frontend && npm install && npm run build && npm cache clean --force

RUN cd vue/backend && npm install && npm run build && npm cache clean --force


FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY --from=vue-builder /app/. .

RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download

RUN go build -o main .


FROM alpine:3.18
WORKDIR /app

RUN apk add --no-cache tzdata

COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/public ./public

EXPOSE 9090

CMD ["./main"]
