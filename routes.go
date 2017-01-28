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
		muxRoute{"Stats", "GET", "/stats", statsHandlerV1},
		muxRoute{"Live", "GET", "/live", liveHandlerV1},
		muxRoute{"WS", "GET", "/ws", websocketHandlerV1},
		muxRoute{"PreTrigger", "POST", "/trigger/{action}", preTriggerHandlerV1},
		muxRoute{"NukeTrigger", "POST", "/nuketrigger", nukeTriggerHandlerV1},
	},
}
