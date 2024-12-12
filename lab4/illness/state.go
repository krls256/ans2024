package illness

import "image/color"

type State int

func (s State) IsSick() bool {
	return s == Sick
}

func (s State) IsHealthy() bool {
	return s == Healthy
}

func (s State) IsImmune() bool {
	return s == Immunity
}

const (
	SickPeriod              = 5
	ImmunityPeriod          = -5
	SickProbability         = 0.2
	LoseImmunityProbability = 0.2
)

const (
	Healthy State = iota + 1
	Sick
	Immunity
)

var States = []State{
	Healthy, Sick, Immunity,
}

var StateColorMap = map[State]color.RGBA{
	Healthy:  {120, 120, 120, 255},
	Sick:     {0, 255, 0, 255},
	Immunity: {0, 0, 255, 255},
}
