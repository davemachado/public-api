NAME        := davemachado/public-api
TAG         := $$(git log -1 --pretty=%h)
IMG         := ${NAME}:${TAG}
LATEST      := ${NAME}:latest

ENV_FLAGS   := CGO_ENABLED=0 GOOS=linux

all: build

dep:
	@go get -v ./...

test:
	@go test -race -coverprofile=coverage.txt -covermode=atomic -v ./...

build: dep test
	@${ENV_FLAGS} go build
	@docker build -t ${IMG} . -f Dockerfile.scratch
	@docker tag ${IMG} ${LATEST}

push: login
	@docker push ${NAME}

login:
	@docker login -u ${DOCKER_USER} -p${DOCKER_PASS}

data:
	@curl -o public-apis.md https://raw.githubusercontent.com/toddmotto/public-apis/master/README.md
	@./md2json.py public-apis.md

