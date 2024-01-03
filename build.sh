export GOOS=linux
go build -o cub
export GOOS=darwin
go build -o cub_macos
export GOOS=windows
go build -o cub.exe
