package messages

type gameIdMessage struct {
	Gameid string
}

type ConnectionMessage struct {
	Gameid   string
	Playerid string
}

type GetPlayers gameIdMessage

type GetGameState gameIdMessage

type PlayerInput struct {
	Value     int
	Timestamp int
	Gameid    string
}
