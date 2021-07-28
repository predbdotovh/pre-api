package main

import (
	"net/http"
)

type muxRoute struct {
	Name    string
	Method  string
	Pattern string
	Handler http.HandlerFunc
}

type muxRoutes []muxRoute
type jsonMuxRoutes map[string]muxRoutes

var jsonRoutes = jsonMuxRoutes{
	"v1": muxRoutes{
		muxRoute{"Root", "GET", "/", rootHandlerV1},
		muxRoute{"Teams", "GET", "/teams", teamsHandlerV1},
		muxRoute{"Stats", "GET", "/stats", statsHandlerV1},
		muxRoute{"Live", "GET", "/live", liveHandlerV1},
		muxRoute{"Rss", "GET", "/rss", rssHandlerV1},
		muxRoute{"WS", "GET", "/ws", websocketHandlerV1},
	},
}

var triggerRoutes = muxRoutes{
	muxRoute{"NukeTrigger", "POST", "/nuke", nukeTriggerHandlerV1},
	muxRoute{"PreTrigger", "POST", "/{action}", preTriggerHandlerV1},
}
