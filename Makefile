all:
	GOOS=windows GOARCH=386 go build -o binaries/respounder-win32.exe respounder.go
	GOOS=windows GOARCH=amd64 go build -o binaries/respounder-win64.exe respounder.go
	GOOS=linux GOARCH=386 go build  -o binaries/respounder-linux32 respounder.go
	GOOS=linux GOARCH=amd64 go build -o binaries/respounder-linux64 respounder.go
	GOOS=darwin GOARCH=386 go build -o binaries/respounder-osx32 respounder.go
	GOOS=darwin GOARCH=amd64 go build -o binaries/respounder-osx64 respounder.go

