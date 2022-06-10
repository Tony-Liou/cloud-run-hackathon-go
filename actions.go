package main

type Action int

const (
	MoveForward Action = iota
	TurnLeft
	TurnRight
	Throw
)

func (a Action) String() string {
	return [...]string{"F", "L", "R", "T"}[a]
}
