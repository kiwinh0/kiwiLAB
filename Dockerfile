# Etapa 1: Construcción
FROM golang:1.21-alpine AS builder

# Instalar dependencias necesarias para CGO (necesario para SQLite)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código
COPY . .

# Compilar la aplicación habilitando CGO para SQLite
RUN CGO_ENABLED=1 GOOS=linux go build -o kiwilab ./cmd/kiwilab/main.go

# Etapa 2: Imagen final ligera
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Copiar el binario y las carpetas necesarias desde la etapa de construcción
COPY --from=builder /app/kiwilab .
COPY --from=builder /app/ui ./ui

# Exponer el puerto
EXPOSE 8080

# Comando para ejecutar
CMD ["./kiwilab"]
