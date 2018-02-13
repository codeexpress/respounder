all:
	GOOS=windows GOARCH=386 go build -o binaries/respounder-x86.exe respounder.go computernames.go
	GOOS=windows GOARCH=amd64 go build -o binaries/respounder-x64.exe respounder.go computernames.go
	GOOS=linux GOARCH=386 go build  -o binaries/respounder-x86 respounder.go computernames.go
	GOOS=linux GOARCH=amd64 go build -o binaries/respounder-x64 respounder.go computernames.go
	GOOS=darwin GOARCH=386 go build -o binaries/respounder-osx respounder.go computernames.go
	GOOS=darwin GOARCH=amd64 go build -o binaries/respounder-osx64 respounder.go computernames.go
