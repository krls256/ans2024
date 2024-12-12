package forest

import (
	cRand "crypto/rand"
	"github.com/samber/lo"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"math/rand/v2"
)

var rng *rand.Rand

func init() {
	seed := [32]byte{}
	_, err := cRand.Read(seed[:])
	if err != nil {
		panic(err)
	}

	rng = rand.New(rand.NewChaCha8(seed))
}

func BasicField(maxPoint Point) Field {
	state := map[Point]Cell{}

	for x := 0; x <= maxPoint.X; x++ {
		for y := 0; y <= maxPoint.Y; y++ {
			itr := rng.IntN(DefaultFieldMaxAge)

			state[Point{X: x, Y: y}] = Cell{
				Iteration: itr,
				State:     IterationToState(itr),
			}
		}
	}

	return Field{
		MaxPoint: maxPoint,
		State:    state,
	}
}

const BurnProbability = 0.5
const DefaultFieldMaxAge = 70

func IterationToState(iteration int) State {
	if iteration == 0 {
		return Fire
	}

	if iteration < 5 {
		return Burnt
	}

	if iteration < 10 {
		return Grass
	}

	if iteration < 50 {
		return RarefiedForests
	}

	return DenseForest
}

type Cell struct {
	State     State
	Iteration int
}

func (c Cell) Next() Cell {
	return Cell{
		State:     IterationToState(c.Iteration + 1),
		Iteration: c.Iteration + 1,
	}
}

/*
	(0, 0) (1, 0)
	(0, 1) (1, 1)
*/

type Point struct {
	X, Y int
}

func (p Point) Top() Point {
	return Point{X: p.X, Y: p.Y - 1}
}

func (p Point) Bottom() Point {
	return Point{X: p.X, Y: p.Y + 1}
}

func (p Point) Left() Point {
	return Point{X: p.X - 1, Y: p.Y}
}

func (p Point) Right() Point {
	return Point{X: p.X + 1, Y: p.Y}
}

type Field struct {
	MaxPoint Point

	State map[Point]Cell
}

func (f Field) Next() Field {
	newState := map[Point]Cell{}

	for p, cell := range f.State {
		if f.HasBurningNeighbor(p) && cell.State.IsFlammable() && rng.Float64() < BurnProbability {
			newState[p] = Cell{State: Fire, Iteration: 0}

			continue
		}

		newState[p] = cell.Next()
	}

	return Field{
		MaxPoint: f.MaxPoint,
		State:    newState,
	}
}

func (f Field) HasBurningNeighbor(p Point) bool {
	top := p.Top()
	if f.ValidPoint(top) && f.State[top].State.IsFiring() {
		return true
	}

	bottom := p.Bottom()
	if f.ValidPoint(bottom) && f.State[bottom].State.IsFiring() {
		return true
	}

	left := p.Left()
	if f.ValidPoint(left) && f.State[left].State.IsFiring() {
		return true
	}

	right := p.Right()
	if f.ValidPoint(right) && f.State[right].State.IsFiring() {
		return true
	}

	return false
}

func (f Field) ValidPoint(p Point) bool {
	return p.X >= 0 && p.Y >= 0 && p.X <= f.MaxPoint.X && p.Y <= f.MaxPoint.Y
}

const wantImageSize = 500

func (f Field) ToPNG() image.Image {
	xCellSize := wantImageSize / (f.MaxPoint.X + 1)
	yCellSize := wantImageSize / (f.MaxPoint.Y + 1)

	xRealImageSize := xCellSize * (f.MaxPoint.X + 1)
	yRealImageSize := yCellSize * (f.MaxPoint.Y + 1)

	img := image.NewRGBA(image.Rect(0, 0, xRealImageSize, yRealImageSize))

	for x := 0; x <= f.MaxPoint.X; x++ {
		for y := 0; y <= f.MaxPoint.Y; y++ {
			rect := image.Rect(
				x*xCellSize,
				y*yCellSize,
				(x+1)*xCellSize,
				(y+1)*yCellSize,
			)
			draw.Draw(img, rect, &image.Uniform{StateColorMap[f.State[Point{X: x, Y: y}].State]}, image.Point{}, draw.Src)
		}
	}

	return img
}

func ImagesToGif(images []image.Image) *gif.GIF {
	outGIF := &gif.GIF{}

	for _, img := range images {
		bounds := img.Bounds()
		palettedImage := image.NewPaletted(bounds, color.Palette(
			lo.Map(lo.Values(StateColorMap), func(item color.RGBA, index int) color.Color {
				return item
			}),
		))

		// Используем FloydSteinberg дизеринг для лучшего качества
		draw.FloydSteinberg.Draw(palettedImage, bounds, img, image.Point{})

		// Добавляем кадр
		outGIF.Image = append(outGIF.Image, palettedImage)
		outGIF.Delay = append(outGIF.Delay, 5)
	}

	return outGIF
}
