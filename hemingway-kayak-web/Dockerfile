# ---- Build ----
FROM golang:1.25 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/server ./cmd/server

# ---- Runtime ----
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/server /server

ENV PORT=8080
EXPOSE 8080

# Las credenciales llegan por variables de entorno (ver .env.example):
#   DB_SERVER, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, JWT_SECRET,
#   REDIS_ADDR, REDIS_PASSWORD, CORS_ORIGINS
ENTRYPOINT ["/server"]
