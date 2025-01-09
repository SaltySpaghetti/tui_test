package models

type State struct {
	Step     Step
	ApiState ApiState
}

type ApiState struct {
	Status Status
	Error  error
}

type Status int

const (
	Idle Status = iota
	Loading
	Success
	Error
)

type Step int

const (
	Searching Step = iota
	Downloading
)
