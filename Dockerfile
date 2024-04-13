# Команды для проверки
# docker build -t sternritter/task-manager:v1 .
# docker run -p 8080:8080 -d -e TODO_PASSWORD='password' sternritter/task-manager:v1

FROM golang:1.22.0

WORKDIR /app

ENV TODO_PORT=8080

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd
COPY internal internal
COPY resources resources
COPY web web
COPY scheduler.db .

EXPOSE 8080
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /task-manager ./cmd/main.go

CMD ["/task-manager"]