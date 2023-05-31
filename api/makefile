APPNAME="go-twitch-api"
depend:
	@go mod tidy

run: depend
	go install github.com/cosmtrek/air@v1.29.0
	air -c .air.toml