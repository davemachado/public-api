NAME    := davemachado/public-api
TAG     := $$(git log -1 --pretty=%h)
IMG     := ${NAME}:${TAG}
LATEST  := ${NAME}:latest

all: test build

test:
	@go test -v ./...

build:
	@docker build -t ${IMG} . -f Dockerfile.scratch
	@docker tag ${IMG} ${LATEST}

push:
	@docker push ${NAME}

login:
	@docker login -u ${DOCKER_USER} -p${DOCKER_PASS}

