package main

import (
	"net/http"

	"github.com/robfig/cron"
)

func main() {
	c := cron.New()
	c.AddFunc("0 */10 * * * *", func() {
		resp, err := http.Get("https://line-anime-bot.herokuapp.com/ping")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
	})
	c.Start()
	select {}
}
