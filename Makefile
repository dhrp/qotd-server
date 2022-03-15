image:
    docker buildx build -t dhrp/qotd --platform=linux/amd64 .
    
push:
    docker push dhrp/qotd