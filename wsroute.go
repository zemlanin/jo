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
		private = lazyJson{
			"type": "PLAYER",
		}
		private["player"] = lazyJson{
			"playerId": "12",
			"gameId":   2222,
			"name":     "whatever",
			"online":   true,
		}

	case "GET_PLAYERS":
		private = lazyJson{
			"type": "PLAYERS",
		}
		private["player"] = nil
		private["players"] = []lazyJson{
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

	case "GET_GAME_STATE":
		public = lazyJson{
			"type": "GAME_STATE",
		}
		game_field := lazyJson{
			"x": 0,
			"y": 2,
		}
		public["gameState"] = lazyJson{
			"gameField": game_field,
			"gameId":    2222,
		}

	}

	return private, public, nil
}
