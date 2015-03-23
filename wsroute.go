package main

import (
	"github.com/mitchellh/mapstructure"
	"jo/game"
	"jo/messages"
	joplayers "jo/players"
	"log"
)

type lazyJson map[string]interface{}
type wsMessage struct {
	gameId  string
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func broadcastPlayers(gameId string) {
	players, err := joplayers.GetPlayers(gameId)
	if err != nil {
		panic(err)
	}
	h.broadcast <- wsMessage{
		gameId: gameId,
		Type:   "PLAYERS",
		Payload: messages.GetPlayersOut{
			Gameid:  gameId,
			Players: players,
		},
	}
}

func broadcastGameState(gameId string) {
	h.broadcast <- wsMessage{
		gameId:  gameId,
		Type:    "GAME_STATE",
		Payload: game.GetState(gameId),
	}
}

func routeMessage(message wsMessage, c *connection) (wsMessage, wsMessage, error) {
	var private wsMessage
	var public wsMessage

	log.Println(message)
	switch message.Type {
	case "CONNECT_PLAYER":
		var player joplayers.Player
		var c_message messages.ConnectionMessage

		if err := mapstructure.Decode(message.Payload, &c_message); err != nil {
			panic(err)
		}

		if c_message.Playerid == "" {
			player = joplayers.GeneratePlayer(c_message.Gameid)
		} else if p, err := joplayers.GetPlayer(c_message.Gameid, c_message.Playerid); err == nil {
			player = p
			joplayers.ConnectPlayer(player.Id)
		} else {
			player = joplayers.GeneratePlayer(c_message.Gameid)
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

		go broadcastPlayers(c_message.Gameid)

	case "GET_PLAYERS":
		var p_message messages.GetPlayers

		if err := mapstructure.Decode(message.Payload, &p_message); err != nil {
			panic(err)
		}
		players, err := joplayers.GetPlayers(p_message.Gameid)
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
		go broadcastGameState(s_message.Gameid)

	case "NEW_PLAYERS_INPUT":
		var i_message messages.PlayerInput
		if err := mapstructure.Decode(message.Payload, &i_message); err != nil {
			panic(err)
		}
		game.InterpretInput(i_message)
		go broadcastGameState(i_message.Gameid)

	case "CHANGE_PLAYER_CONTROLLERS":
		var c_message messages.ControllersMessage
		if err := mapstructure.Decode(message.Payload, &c_message); err != nil {
			panic(err)
		}
		joplayers.ChangePlayerControllers(c_message.Playerid, c_message.Controllers)
		go broadcastPlayers(c_message.Gameid)
		// TODO: send new player_state to target player
	}

	return private, public, nil
}
