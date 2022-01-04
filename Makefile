.PHONY: build deploy gomodgen

build:
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/lambda lambda/main.go

deploy: build
	npx sls deploy --verbose

upload-assets:
	aws s3 sync --delete assets s3://in.fkn.space/i/stickers/
