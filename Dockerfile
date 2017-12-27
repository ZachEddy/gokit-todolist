FROM alpine:latest

ADD cmd/server/main /bin/todolist
EXPOSE 8080 8080

CMD todolist
