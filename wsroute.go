package main

import (
	"log"
)

type lazyJson map[string]interface{}

func routeMessage(message lazyJson) (lazyJson, lazyJson, error) {
	var private lazyJson = nil
	var public lazyJson = nil

	log.Println(message)
	switch message["type"] {
	case "CONNECT_PLAYER":
		message["type"] = "PLAYER"
		message["player"] = lazyJson{
			"playerId": "12",
			"gameId":   2222,
			"name":     "whatever",
			"online":   true,
		}

		private = message

	case "GET_PLAYERS":
		message["type"] = "PLAYERS"
		message["players"] = []lazyJson{
			{
				"gameId": 2222,
				"name":   "whatever",
				"online": true,
			},
			{
				"gameId": 2222,
				"name":   "another",
				"online": false,
			},
		}
		private = message

	case "GET_GAME_STATE":
		message["type"] = "GAME_STATE"
		game_field := lazyJson{
			"x": 0,
			"y": 2,
		}
		message["gameState"] = lazyJson{
			"gameField": game_field,
			"gameId":    2222,
		}
		public = message

	}

	return private, public, nil
}
