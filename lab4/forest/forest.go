package forest

import "image/color"

type State int

func (s State) IsFiring() bool {
	return s == Fire
}

func (s State) IsFlammable() bool {
	return s != Fire && s != Burnt
}

const (
	Fire State = iota + 1
	Burnt
	Grass
	RarefiedForests
	DenseForest
)

var States = []State{
	Fire, Burnt, Grass, RarefiedForests, DenseForest,
}

var StateColorMap = map[State]color.RGBA{
	Fire:            {255, 0, 0, 255},   // RED
	Burnt:           {165, 42, 42, 255}, // BROWN
	Grass:           {144, 238, 144, 255},
	RarefiedForests: {R: 34, G: 221, B: 34, A: 255},
	DenseForest:     {R: 25, G: 165, B: 25, A: 255},
}
