package players

import (
	"errors"
)

type Player struct {
	Id      string `json:"id"`
	Game_id string `json:"gameId"`
	Name    string `json:"name"`
	Online  bool   `json:"online"`
}

var playersList = []Player{
	Player{
		Id:      "6a5674",
		Game_id: "2222",
		Name:    "whatever",
		Online:  true,
	},
	Player{
		Id:      "067df5",
		Game_id: "2222",
		Name:    "somebody",
		Online:  false,
	},
}

func GetPlayers(game_id string) ([]Player, error) {
	var result []Player

	for _, p := range playersList {
		if p.Game_id == game_id {
			result = append(result, p)
		}
	}

	return result, nil
}

func GetPlayer(game_id, id string) (Player, error) {
	for _, p := range playersList {
		if p.Id == id && p.Game_id == game_id {
			return p, nil
		}
	}

	return Player{}, errors.New("not found")
}

func SavePlayer(player Player) error {
	return nil
}
