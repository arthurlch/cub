export GOOS=linux
go build -o cub cub.go
export GOOS=Windows
go build -o cub.exe cub.go