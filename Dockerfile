# --- STAGE 1: Build the React Frontend ---
FROM node:24-alpine AS frontend-builder
WORKDIR /app
# Copy only package files to cache dependencies
COPY frontend/package*.json ./
RUN npm install
# Copy the rest of the frontend
COPY frontend/ ./
RUN npm run build
RUN echo "--- LISTING DIRECTORIES ---" && ls -F /app

# --- STAGE 2: Build the Go Backend ---
FROM golang:1.26-alpine AS backend-builder
WORKDIR /app
# Copy only Go module files to cache
COPY backend/go.mod backend/go.sum ./
RUN go mod download
# Copy the rest of the backend
COPY backend/ ./
# Copy the built React assets from Stage 1 into the backend folder
COPY --from=frontend-builder /dist ./dist
# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o obligo main.go

# --- STAGE 3: Final Tiny Runtime ---
FROM alpine:latest
WORKDIR /root/
# Copy the binary and the dist folder
COPY --from=backend-builder /app/obligo .
COPY --from=backend-builder /app/dist ./dist

# Set the environment variable so Go knows where to look
ENV DIST_PATH=/root/dist
ENV APP_PORT=1234

EXPOSE 1234
CMD ["./obligo"]
