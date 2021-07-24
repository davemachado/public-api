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

build: dep test html
	@${ENV_FLAGS} go build -o public-api
	@docker build -t ${IMG} . -f Dockerfile.scratch
	@docker tag ${IMG} ${LATEST}

push: login
	@docker push ${NAME}

login:
	@docker login -u ${DOCKER_USER} -p${DOCKER_PASS}

data:
	@curl -o /tmp/public-apis.md https://raw.githubusercontent.com/toddmotto/public-apis/master/README.md
	@./md2json /tmp/public-apis.md > entries.json
	@rm /tmp/public-apis.md

html:
	mkdir -p static
	pandoc --from markdown_github --to html --standalone README.md > static/index.html
