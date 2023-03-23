# Se establece la imagen base de golang con Alpine 3.17 como sistema operativo.
FROM golang:alpine3.17 AS builder

# Se actualiza el índice de paquetes de Alpine y se instala git.
RUN apk update && apk add --no-cache git

# Se establece el GOPROXY en la variable de entorno para que Go utilice un proxy específico.
ENV GOPROXY=https://proxy.golang.org

# Se establece el directorio de trabajo en el contenedor.
WORKDIR /go/src/app

# Se copia el contenido actual del directorio de trabajo al contenedor.
COPY . .

# Se habilita el uso de módulos Go.
RUN go env -w GO111MODULE=on

# Se inicializa el módulo dockergo.
RUN go mod init dockergo

# Se descargan las dependencias necesarias para el proyecto.
RUN go get -d -v ./...

# Se descargan e instalan las dependencias específicas del proyecto.
RUN go get github.com/0x9ef/openai-go
RUN go get github.com/go-telegram-bot-api/telegram-bot-api

# Se instala la aplicación dentro del contenedor.
RUN go install -v ./...

# Se compila la aplicación para que se ejecute en un sistema operativo Linux.
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/app main.go

# Se establece la imagen base desde cero para reducir el tamaño del contenedor.
FROM scratch

# Se copia el archivo binario de la aplicación compilada desde la imagen anterior.
COPY --from=builder /go/bin/app/ /go/bin/app

# Se copia el certificado CA para establecer conexiones seguras.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Se establece el comando de inicio de la aplicación.
ENTRYPOINT ["go/bin/app"]
