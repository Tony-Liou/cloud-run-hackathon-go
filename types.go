package main

type ArenaUpdate struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
	Arena ArenaInfo `json:"arena"`
}

type ArenaInfo struct {
	Dimensions []int                  `json:"dims"`
	State      map[string]PlayerState `json:"state"`
}

type PlayerState struct {
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Direction string `json:"direction"`
	WasHit    bool   `json:"wasHit"`
	Score     int    `json:"score"`
}
