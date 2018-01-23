NAME        := davemachado/public-api
TAG         := $$(git log -1 --pretty=%h)
IMG         := ${NAME}:${TAG}
LATEST      := ${NAME}:latest

ENV_FLAGS   := CGO_ENABLED=0 GOOS=linux

all: test build

test:
	@go test -v ./...

build:
	@${ENV_FLAGS} go build
	@docker build -t ${IMG} . -f Dockerfile.scratch
	@docker tag ${IMG} ${LATEST}

push: login
	@docker push ${NAME}

login:
	@docker login -u ${DOCKER_USER} -p${DOCKER_PASS}
