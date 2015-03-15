package main

import (
	"log"
)

type lazyJson map[string]interface{}
type wsMessage struct{
	Type string `json:"type"`
	Payload interface{} `json:"payload"`
}

func routeMessage(message wsMessage) (wsMessage, wsMessage, error) {
	var private wsMessage
	var public wsMessage

	log.Println(message)
	switch message.Type {
	case "CONNECT_PLAYER":
		player := lazyJson{
			"playerId": "12",
			"gameId":   2222,
			"name":     "whatever",
			"online":   true,
		}
		private = wsMessage{
			Type: "PLAYER",
			Payload: player,
		}

	case "GET_PLAYERS":
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

		private = wsMessage{
			Type: "PLAYERS",
			Payload: players,
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
		public = wsMessage{
			Type: "GAME_STATE",
			Payload: game_state,
		}
	}

	return private, public, nil
}
