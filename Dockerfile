# Команды для запуска приложения в контейнере
# docker build -t sternritter/task-manager:v1 .
# docker run -d -p 7540:7540 -e TODO_PASSWORD='password' -v $(pwd)/scheduler.db:/scheduler.db --name task-manager sternritter/task-manager:v1

FROM golang:1.22.0

WORKDIR /app

ENV TODO_PORT=7540

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 7540
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /task-manager ./cmd/main.go

VOLUME /scheduler.db

CMD ["/task-manager"]