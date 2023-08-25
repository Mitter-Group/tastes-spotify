# Usar la imagen oficial de Go
FROM golang:1.20.6

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar los archivos de tu proyecto al contenedor
COPY . .

# Construir tu aplicación
RUN go build -o main ./cmd/api

# Exponer el puerto que tu aplicación usará
EXPOSE 8080

# Ejecutar tu aplicación cuando el contenedor se inicie
CMD ["./main"]
