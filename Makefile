default: build

build:
	@go build -o cmd/server/main cmd/server/main.go

clean:
	@rm cmd/server/main

alpine:
	@GOOS=linux GOARCH=amd64 go build -o cmd/server/main cmd/server/main.go

docker: alpine
	@docker build -t $(IMG) --no-cache .
	@docker tag $(IMG) $(LATEST)

push:
	@docker push $(NAME)

NAME 		:= zacheddy/gokit-todolist
TAG			:= $(shell cat VERSION)-$(shell git rev-parse --short HEAD)
IMG 		:= $(NAME):$(TAG)
LATEST 	:= $(NAME):latest
