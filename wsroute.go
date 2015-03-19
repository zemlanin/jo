package main

import (
	"github.com/mitchellh/mapstructure"
	"jo/game"
	"jo/messages"
	"jo/players"
	"log"
)

type lazyJson map[string]interface{}
type wsMessage struct {
	gameId  string
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func routeMessage(message wsMessage, c *connection) (wsMessage, wsMessage, error) {
	var private wsMessage
	var public wsMessage

	log.Println(message)
	switch message.Type {
	case "CONNECT_PLAYER":
		var player players.Player
		var c_message messages.ConnectionMessage

		if err := mapstructure.Decode(message.Payload, &c_message); err != nil {
			panic(err)
		}

		if c_message.Playerid == "" {
			player = players.GeneratePlayer(c_message.Gameid)
		} else if p, err := players.GetPlayer(c_message.Gameid, c_message.Playerid); err == nil {
			player = p
			players.ConnectPlayer(player.Id)
		} else {
			player = players.GeneratePlayer(c_message.Gameid)
		}

		if c.playerId == "" {
			c.playerId = c_message.Playerid
		}

		if c.gameId == "" {
			c.gameId = c_message.Gameid
		}

		private = wsMessage{
			gameId:  c_message.Gameid,
			Type:    "PLAYER",
			Payload: player,
		}

		players, err := players.GetPlayers(c_message.Gameid)
		if err != nil {
			panic(err)
		}
		public = wsMessage{
			gameId: c_message.Gameid,
			Type:   "PLAYERS",
			Payload: messages.GetPlayersOut{
				Gameid:  c_message.Gameid,
				Players: players,
			},
		}
		// send state on connect

	case "GET_PLAYERS":
		var p_message messages.GetPlayers

		if err := mapstructure.Decode(message.Payload, &p_message); err != nil {
			panic(err)
		}
		players, err := players.GetPlayers(p_message.Gameid)
		if err != nil {
			panic(err)
		}

		private = wsMessage{
			gameId: p_message.Gameid,
			Type:   "PLAYERS",
			Payload: messages.GetPlayersOut{
				Gameid:  p_message.Gameid,
				Players: players,
			},
		}

	case "GET_GAME_STATE":
		var s_message messages.GetGameState

		if err := mapstructure.Decode(message.Payload, &s_message); err != nil {
			panic(err)
		}
		public = wsMessage{
			gameId:  s_message.Gameid,
			Type:    "GAME_STATE",
			Payload: game.GetState(s_message.Gameid),
		}

	case "NEW_PLAYERS_INPUT":
		var i_message messages.PlayerInput
		if err := mapstructure.Decode(message.Payload, &i_message); err != nil {
			panic(err)
		}
		game.InterpretInput(i_message)
		public = wsMessage{
			gameId:  i_message.Gameid,
			Type:    "GAME_STATE",
			Payload: game.GetState(i_message.Gameid),
		}
	}

	return private, public, nil
}
