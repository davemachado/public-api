NAME        := davemachado/public-api
TAG         := $$(git log -1 --pretty=%h)
IMG         := ${NAME}:${TAG}
LATEST      := ${NAME}:latest

all: build

dep:
	@go get -v ./...

test:
	@go test -race -coverprofile=coverage.txt -covermode=atomic -v ./...

build: dep test html
	@docker build -t ${IMG} .
	@docker tag ${IMG} ${LATEST}

data: html
	@curl -o /tmp/public-apis.md https://raw.githubusercontent.com/public-apis/public-apis/master/README.md
	@./md2json /tmp/public-apis.md > entries.json
	@rm /tmp/public-apis.md

html:
	mkdir -p static
	pandoc --from markdown_github --to html --standalone README.md > static/index.html
