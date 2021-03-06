package game

import (
	// "log"
	"jo/messages"
)

type lazyJson map[string]interface{}
type gameField map[string]int

var game_field = gameField{
	"x": 0,
	"y": 0,
}

func GetState(game_id string) lazyJson {
	game_state := lazyJson{
		"gameField": game_field,
		"gameId":    game_id,
	}

	return game_state
}

func InterpretInput(input messages.PlayerInput) {
	switch input.Value {
	case 0:
		game_field["y"] = game_field["y"] + 1

	case 1:
		game_field["y"] = game_field["y"] - 1

	case 2:
		game_field["x"] = game_field["x"] - 1

	case 3:
		game_field["x"] = game_field["x"] + 1

	}
}
