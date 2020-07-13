build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/line-anime-bot main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/gocron cron/cronjob.go
