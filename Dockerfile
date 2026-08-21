# syntax=docker/dockerfile:1.7
FROM golang:1.26-alpine AS api-build
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/stockpilot ./cmd/server

FROM node:22-alpine AS web-build
WORKDIR /src/web
COPY web/package*.json ./
RUN npm install --no-audit --no-fund
COPY web ./
RUN npm run build

FROM alpine:3.22
RUN addgroup -S stockpilot && adduser -S -G stockpilot stockpilot
WORKDIR /app
COPY --from=api-build /out/stockpilot /app/stockpilot
COPY --from=web-build /src/web/dist /app/web/dist
COPY migrations /app/migrations
USER stockpilot
EXPOSE 8080
ENTRYPOINT ["/app/stockpilot"]
