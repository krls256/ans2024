package main

import (
	"fmt"
	"image"
	"image/gif"
	"nsa/lab4/illness"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	f := illness.BasicField(illness.Point{X: 25, Y: 25})

	images := []image.Image{}
	stats := [][3]int{}
	for range 250 {
		f = f.Next()
		stats = append(stats, f.Stat())

		images = append(images, f.ToPNG())
	}

	printStats(stats)

	gifFile := illness.ImagesToGif(images)
	file, err := os.Create(fmt.Sprintf("./lab4/results/illness_%.2f_%.2f_%v.gif",
		illness.SickProbability, illness.LoseImmunityProbability, time.Now().Format(time.DateTime)))
	if err != nil {
		panic(err)
	}

	gif.EncodeAll(file, gifFile)
}

func printStats(stats [][3]int) {
	h := []string{}
	s := []string{}
	i := []string{}

	for _, stat := range stats {
		h = append(h, strconv.Itoa(stat[0]))
		s = append(s, strconv.Itoa(stat[1]))
		i = append(i, strconv.Itoa(stat[2]))
	}

	fmt.Printf("h = [%v]\n", strings.Join(h, ", "))
	fmt.Printf("s = [%v]\n", strings.Join(s, ", "))
	fmt.Printf("i = [%v]\n", strings.Join(i, ", "))
}
