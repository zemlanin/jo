package players

import (
	"crypto/md5"
	"crypto/rand"
	"errors"
	"fmt"
	// "log"
)

type Player struct {
	Online bool   `json:"online"`
	Id     string `json:"id"`
	Gameid string `json:"gameId"`
}

var playersList = []Player{}

func GetPlayers(game_id string) ([]Player, error) {
	var result []Player

	for _, p := range playersList {
		if p.Gameid == game_id {
			result = append(result, p)
		}
	}

	return result, nil
}

func GetPlayer(game_id, id string) (Player, error) {
	for _, p := range playersList {
		if p.Id == id && p.Gameid == game_id {
			return p, nil
		}
	}

	return Player{}, errors.New("not found")
}

func GeneratePlayer(game_id string) Player {
	rand_bytes := make([]byte, 120)
	_, err := rand.Read(rand_bytes)
	if err != nil {
		panic(err)
	}
	hash := md5.Sum(rand_bytes)

	player_id := fmt.Sprintf("%x", hash[:4])
	player := Player{
		Id:     player_id,
		Gameid: game_id,
		Online: true,
	}

	playersList = append(playersList, player)

	return player
}

func DisconnectPlayer(playerId string) {
	var newPlayersList []Player

	for _, p := range playersList {
		if p.Id == playerId {
			p.Online = false
		}
		newPlayersList = append(newPlayersList, p)
	}

	playersList = newPlayersList
}

func ConnectPlayer(playerId string) {
	var newPlayersList []Player

	for _, p := range playersList {
		if p.Id == playerId {
			p.Online = true
		}
		newPlayersList = append(newPlayersList, p)
	}

	playersList = newPlayersList
}

func SavePlayer(player Player) error {
	return nil
}
