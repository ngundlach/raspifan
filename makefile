build:
	env GOARCH=arm64 GOOS=linux go build -x -o out/
release:
	env GOARCH=arm64 GOOS=linux go build -ldflags "-w -s" -trimpath -x -o out/
