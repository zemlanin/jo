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
		player := lazyJson{
			"playerId": "12",
			"gameId":   2222,
			"name":     "whatever",
			"online":   true,
		}
		private = lazyJson{
			"type": "PLAYER",
			"payload": player,
		}

	case "GET_PLAYERS":
		private = lazyJson{
			"type": "PLAYERS",
		}
		players := []lazyJson{
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

		private = lazyJson{
			"type": "PLAYERS",
			"payload": players,
		}

	case "GET_GAME_STATE":
		game_field := lazyJson{
			"x": 0,
			"y": 2,
		}
		game_state := lazyJson{
			"gameField": game_field,
			"gameId":    2222,
		}
		public = lazyJson{
			"type": "GAME_STATE",
			"payload": game_state,
		}
	}

	return private, public, nil
}
