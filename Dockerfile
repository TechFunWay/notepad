# Stage 1: Build frontend
FROM node:20-alpine AS web-build
WORKDIR /app/web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

# Stage 2: Build backend
FROM golang:1.22-alpine AS server-build
RUN apk add --no-cache git
WORKDIR /app/server
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ ./
COPY --from=web-build /app/web/dist ./static/dist
ARG VERSION=dev
ARG BUILD_TIME=unknown
ARG GIT_COMMIT=unknown
RUN CGO_ENABLED=0 go build -ldflags "-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}" -o /notepad ./main.go

# Stage 3: Runtime
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=server-build /notepad /app/notepad
EXPOSE 8904
VOLUME ["/app/data"]
ENV PORT=8904 DB_PATH=/app/data/notepad.db DATA_DIR=/app/data
HEALTHCHECK --interval=30s --timeout=3s CMD wget -q --spider http://localhost:8904/api/health || exit 1
ENTRYPOINT ["/app/notepad"]
