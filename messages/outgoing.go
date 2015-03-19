package messages

import "jo/players"

type GetPlayersOut struct {
	Gameid  string           `json:"gameId"`
	Players []players.Player `json:"players"`
}
