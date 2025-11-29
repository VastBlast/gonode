build:
	go build -ldflags="-s -w" gonode.go
	$(if $(shell command -v upx), upx gonode)

mac:
	GOOS=darwin go build -ldflags="-s -w" -o gonode-darwin gonode.go
	$(if $(shell command -v upx), upx gonode-darwin)

win:
	GOOS=windows go build -ldflags="-s -w" -o gonode.exe gonode.go
	$(if $(shell command -v upx), upx gonode.exe)

linux:
	GOOS=linux go build -ldflags="-s -w" -o gonode-linux gonode.go
	$(if $(shell command -v upx), upx gonode-linux)