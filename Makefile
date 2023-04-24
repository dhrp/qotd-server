tag = $(shell git describe --tags --abbrev=0)

tag:
	@echo $(tag)

dockerize:
	docker buildx build -t dhrp/qotd:latest -t dhrp/qotd:$(tag) --platform=linux/amd64 .
    
push:
	docker push --all-tags dhrp/qotd:latest

