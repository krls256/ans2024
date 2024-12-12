package illness

import (
	cRand "crypto/rand"
	"github.com/samber/lo"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"math/rand/v2"
)

func IterationToState(iteration int) State {
	if iteration == 0 {
		return Healthy
	}

	if iteration < 0 {
		return Immunity
	}

	return Sick
}

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
			state[Point{X: x, Y: y}] = Cell{
				Iteration: 0,
				State:     IterationToState(0),
			}
		}
	}

	state[Point{X: maxPoint.X / 2, Y: maxPoint.Y / 2}] = Cell{
		Iteration: SickPeriod,
		State:     IterationToState(SickPeriod),
	}

	return Field{
		MaxPoint: maxPoint,
		State:    state,
	}
}

type Field struct {
	MaxPoint Point

	State map[Point]Cell
}

func (f Field) Neighbors(p Point) []Point {
	return []Point{
		p.Top(f.MaxPoint),
		p.TopLeft(f.MaxPoint),
		p.Left(f.MaxPoint),
		p.BottomLeft(f.MaxPoint),
		p.Bottom(f.MaxPoint),
		p.BottomRight(f.MaxPoint),
		p.Right(f.MaxPoint),
		p.TopRight(f.MaxPoint),
	}
}
func (f Field) HasSickNeighbor(p Point) bool {
	for _, n := range f.Neighbors(p) {
		if f.State[n].State.IsSick() {
			return true
		}
	}

	return false
}

func (f Field) Next() Field {
	newState := map[Point]Cell{}

	for p, cell := range f.State {
		if f.HasSickNeighbor(p) && cell.State.IsHealthy() && rng.Float64() < SickProbability {
			newState[p] = Cell{State: IterationToState(SickPeriod), Iteration: SickPeriod}

			continue
		}

		newState[p] = cell.Next()
	}

	return Field{
		MaxPoint: f.MaxPoint,
		State:    newState,
	}
}

func (f Field) Stat() [3]int {
	// h/s/i
	stat := [3]int{}

	for _, s := range f.State {
		if s.State.IsHealthy() {
			stat[0]++
		}

		if s.State.IsSick() {
			stat[1]++
		}

		if s.State.IsImmune() {
			stat[2]++
		}
	}

	return stat
}

type Cell struct {
	State     State
	Iteration int
}

func (c Cell) Next() Cell {
	if c.State.IsHealthy() {
		return c
	}

	if c.State.IsSick() {
		cell := Cell{
			State:     IterationToState(c.Iteration - 1),
			Iteration: c.Iteration - 1,
		}

		if cell.State.IsHealthy() {
			cell = Cell{
				State:     IterationToState(ImmunityPeriod),
				Iteration: ImmunityPeriod,
			}
		}

		return cell
	}

	if rand.Float64() < LoseImmunityProbability {
		return Cell{
			State:     IterationToState(c.Iteration + 1),
			Iteration: c.Iteration + 1,
		}
	}

	return c
}

type Point struct {
	X, Y int
}

func (p Point) Round(maxPoint Point) Point {
	return Point{X: (p.X + maxPoint.X) % (maxPoint.X + 1), Y: (p.Y + maxPoint.Y) % (maxPoint.Y + 1)}
}

func (p Point) Top(maxPoint Point) Point {
	return Point{X: p.X, Y: p.Y - 1}.Round(maxPoint)
}

func (p Point) TopLeft(maxPoint Point) Point {
	return Point{X: p.X - 1, Y: p.Y - 1}.Round(maxPoint)
}

func (p Point) TopRight(maxPoint Point) Point {
	return Point{X: p.X + 1, Y: p.Y - 1}.Round(maxPoint)
}

func (p Point) Bottom(maxPoint Point) Point {
	return Point{X: p.X, Y: p.Y + 1}.Round(maxPoint)
}

func (p Point) BottomLeft(maxPoint Point) Point {
	return Point{X: p.X - 1, Y: p.Y + 1}.Round(maxPoint)
}

func (p Point) BottomRight(maxPoint Point) Point {
	return Point{X: p.X + 1, Y: p.Y + 1}.Round(maxPoint)
}

func (p Point) Left(maxPoint Point) Point {
	return Point{X: p.X - 1, Y: p.Y}.Round(maxPoint)
}

func (p Point) Right(maxPoint Point) Point {
	return Point{X: p.X + 1, Y: p.Y}.Round(maxPoint)
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
