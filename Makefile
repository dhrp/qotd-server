tag = $(shell git describe --tags --abbrev=0)
fulltag = $(shell git describe --tags --abbrev=7)


tag:
	@echo $(tag)

fulltag:
	@echo $(fulltag)

dockerize:
	docker buildx build -t dhrp/qotd:latest -t dhrp/qotd:$(tag) --platform=linux/arm64/v8,linux/amd64 --push .
    
push:
	docker push --all-tags dhrp/qotd:latest

run:
	docker run -p 8000:80 dhrp/qotd:$(tag)