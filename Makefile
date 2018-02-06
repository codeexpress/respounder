all:
	GOOS=windows GOARCH=386 go build -o binaries/respounder-x86.exe respounder.go
	GOOS=windows GOARCH=amd64 go build -o binaries/respounder-x64.exe respounder.go
	GOOS=linux GOARCH=386 go build  -o binaries/respounder-x86 respounder.go
	GOOS=linux GOARCH=amd64 go build -o binaries/respounder-x64 respounder.go
	GOOS=darwin GOARCH=386 go build -o binaries/respounder-osx respounder.go
	GOOS=darwin GOARCH=amd64 go build -o binaries/respounder-osx64 respounder.go
