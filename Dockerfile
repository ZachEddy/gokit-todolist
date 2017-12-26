FROM alpine:latest

ADD cmd/server/main /bin/todolist
CMD todolist
