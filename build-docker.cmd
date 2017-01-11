set GOOS=linux
set GOARCH=386

go build -o docker-cleaner

REM docker build -f Dockerfile -t jdfischer/docker-cleaner:latest .
