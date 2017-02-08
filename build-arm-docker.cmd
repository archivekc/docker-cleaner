set GOOS=linux
set GOARCH=arm
set GOARM=7

go build -o docker-cleaner

REM docker build -f Dockerfile.arm -t jdfischer/docker-cleaner:rpi-latest .
