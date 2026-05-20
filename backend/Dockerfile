FROM golang:1.22-alpine AS base

WORKDIR /workspace/backend

COPY backend/go.mod ./
RUN go mod download

COPY backend ./

FROM base AS development

ENV HOST=0.0.0.0
ENV PORT=8080
ENV DB_PATH=/data/kanban.db

EXPOSE 8080

CMD ["go", "run", "./cmd/api"]

FROM base AS test

RUN go test ./...

FROM base AS build

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/api ./cmd/api

FROM alpine:3.20 AS runtime

WORKDIR /app

COPY --from=build /out/api /usr/local/bin/api

ENV HOST=0.0.0.0
ENV PORT=8080
ENV DB_PATH=/data/kanban.db

EXPOSE 8080

CMD ["api"]
