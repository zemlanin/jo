package main

import (
	"github.com/mitchellh/mapstructure"
	"jo/game"
	"log"
)

type lazyJson map[string]interface{}
type wsMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func routeMessage(message wsMessage) (wsMessage, wsMessage, error) {
	var private wsMessage
	var public wsMessage
	var ok bool
	var incoming map[string]interface{}
	var game_id string

	if incoming, ok = message.Payload.(map[string]interface{}); !ok {
		return private, public, nil
	}

	if game_id, ok = incoming["gameId"].(string); !ok {
		return private, public, nil
	}

	log.Println(message)
	switch message.Type {
	case "CONNECT_PLAYER":
		player := lazyJson{
			"playerId": "12",
			"gameId":   game_id,
			"name":     "whatever",
			"online":   true,
		}
		private = wsMessage{
			Type:    "PLAYER",
			Payload: player,
		}

	case "GET_PLAYERS":
		players := []lazyJson{
			{
				"gameId": game_id,
				"name":   "whatever",
				"online": true,
			},
			{
				"gameId": game_id,
				"name":   "another",
				"online": false,
			},
		}

		private = wsMessage{
			Type:    "PLAYERS",
			Payload: players,
		}

	case "GET_GAME_STATE":
		public = wsMessage{
			Type:    "GAME_STATE",
			Payload: game.GetState(game_id),
		}

	case "NEW_PLAYERS_INPUT":
		var input game.PlayerInput
		if err := mapstructure.Decode(incoming, &input); err != nil {
			panic(err)
		}
		game.InterpretInput(input)
		public = wsMessage{
			Type:    "GAME_STATE",
			Payload: game.GetState(game_id),
		}
	}

	return private, public, nil
}
