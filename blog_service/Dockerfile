FROM golang:1.25.5-alpine

RUN adduser --disabled-password --gecos '' gouser

RUN mkdir -p /home/gouser/project

WORKDIR /home/gouser/project

COPY . .

RUN chown -R gouser:gouser /home/gouser/project

USER gouser

ENV SERVER_PORT=8000

RUN go mod tidy
RUN go build -o build/server cmd/main.go

EXPOSE 8000

CMD ["./build/server"]
 