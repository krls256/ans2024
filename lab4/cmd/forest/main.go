package main

import (
	"fmt"
	"image"
	"image/gif"
	"nsa/lab4/forest"
	"os"
	"time"
)

func main() {
	f := forest.BasicField(forest.Point{X: 50, Y: 50})

	images := []image.Image{}
	for range 250 {
		f = f.Next()

		images = append(images, f.ToPNG())
	}

	gifFile := forest.ImagesToGif(images)
	file, err := os.Create(fmt.Sprintf("./lab4/results/forest_%.4f_%v.gif",
		forest.BurnProbability, time.Now().Format(time.DateTime)))
	if err != nil {
		panic(err)
	}

	gif.EncodeAll(file, gifFile)
}
