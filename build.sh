export GOOS=linux
go build -o cub cub.go
export GOOS=darwin
go build -o cub_macos cub.go
export GOOS=windows
go build -o cub.exe cub.go
