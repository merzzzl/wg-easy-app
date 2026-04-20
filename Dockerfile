FROM golang:1.24.6-alpine AS builder-backend

WORKDIR /app/backend

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .

ENV CGO_ENABLED=0
RUN go build -o /out/wg-easy-app ./cmd/app

FROM node:24-alpine AS builder-frontend

WORKDIR /app/frontend

COPY frontend/package*.json ./
RUN npm ci

COPY frontend/ .
RUN npm run build

FROM alpine:3.22

WORKDIR /app

COPY --from=builder-backend /out/wg-easy-app /app/wg-easy-app
COPY --from=builder-frontend /app/frontend/dist /app/static

EXPOSE 8080

CMD ["/app/wg-easy-app"]
