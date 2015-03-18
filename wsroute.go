package main

import (
	"github.com/mitchellh/mapstructure"
	"jo/game"
	"jo/players"
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
		var player players.Player
		var c_message players.ConnectionMessage

		if err := mapstructure.Decode(message.Payload, &c_message); err != nil {
			panic(err)
		}

		if c_message.Playerid == "" {
			player = players.GeneratePlayer(c_message.Gameid)
		} else if p, err := players.GetPlayer(c_message.Gameid, c_message.Playerid); err == nil {
			player = p
		} else {
			player = players.GeneratePlayer(c_message.Gameid)
		}

		private = wsMessage{
			Type:    "PLAYER",
			Payload: player,
		}

		players, err := players.GetPlayers(c_message.Gameid)
		if err != nil {
			panic(err)
		}
		public = wsMessage{
			Type:    "PLAYERS",
			Payload: players,
		}

	case "GET_PLAYERS":
		players, err := players.GetPlayers(game_id)
		if err != nil {
			panic(err)
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
		if err := mapstructure.Decode(message.Payload, &input); err != nil {
			panic(err)
		}
		game.InterpretInput(input)
		public = wsMessage{
			Type:    "GAME_STATE",
			Payload: game.GetState(input.Gameid),
		}
	}

	return private, public, nil
}
