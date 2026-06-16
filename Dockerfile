# --- STAGE 1: Build the React Frontend ---
FROM node:24-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# --- STAGE 2: Build the Go Backend ---
FROM golang:1.26-alpine AS backend-builder
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
COPY --from=frontend-builder /app/dist ./dist
RUN CGO_ENABLED=0 GOOS=linux go build -o obligo main.go

# --- STAGE 3: Final Tiny Runtime ---
FROM alpine:latest
WORKDIR /app
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
COPY --from=backend-builder /app/obligo .
COPY --from=backend-builder /app/dist ./dist
RUN chown -R appuser:appgroup /app

EXPOSE 1234
ENV DIST_PATH=/root/dist
ENV APP_PORT=1234

CMD ["./obligo"]
