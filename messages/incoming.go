package messages

type GameIdMessage struct {
	Gameid string
}

type ConnectionMessage struct {
	Gameid   string
	Playerid string
}

type GetPlayers GameIdMessage

type GetGameState GameIdMessage

type PlayerInput struct {
	Value     int
	Timestamp int
	Gameid    string
}
