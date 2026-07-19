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
RUN CGO_ENABLED=0 GOOS=linux go build -o obligo main.go

# --- STAGE 3: Final Tiny Runtime ---
FROM scratch AS product
WORKDIR /app
COPY --from=frontend-builder /app/dist ./dist
COPY --from=backend-builder /app/obligo .

EXPOSE 1234
ENV DIST_PATH=/app/dist
ENV APP_PORT=1234

CMD ["./obligo"]
