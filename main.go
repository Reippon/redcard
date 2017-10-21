package main

import (
	"flag"

	"github.com/gramework/gramework"
	"github.com/parnurzeal/gorequest"
)

var (
	apiKey = flag.String("apiKey", "", "apifootball api key")
)

func main() {
	app := gramework.New()
	app.RegFlags()
	flag.Parse()

	if *apiKey == "" {
		app.Logger.Fatal("apiKey is required")
	}

	app.GET("/", func(ctx *gramework.Context) {
		_, body, errs := gorequest.New().
			Get("https://apifootball.com/api/?action=get_events&from=2017-10-20&to=2017-10-21&APIkey=" + *apiKey).
			EndBytes()

		if errs != nil {
			ctx.JSONError(errs)
			return
		}

		var eventsResult []interface{}
		_, err := ctx.UnJSONBytes(body, &eventsResult)
		if err != nil {
			ctx.JSONError(err)
			return
		}

		cards := make([]interface{}, 0)

		for _, v := range eventsResult {
			event, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			ctx.SetContentType("application/json")
			cards = append(cards, event["cards"])
		}

		ctx.SetContentType("application/json")
		ctx.JSON(cards)
	})

	app.ListenAndServe()
}
