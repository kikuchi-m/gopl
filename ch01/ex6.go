package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"os"
	"time"
)

var palette = []color.Color{
	color.Black,
	color.RGBA{0x00, 0x40, 0x80, 0xff},
	color.RGBA{0x10, 0x50, 0x90, 0xff},
	color.RGBA{0x20, 0x60, 0xa0, 0xff},
	color.RGBA{0x30, 0x70, 0xb0, 0xff},
	color.RGBA{0x40, 0x80, 0xc0, 0xff},
	color.RGBA{0x50, 0x90, 0xd0, 0xff},
	color.RGBA{0x60, 0xa0, 0xe0, 0xff},
	color.RGBA{0x70, 0xb0, 0xf0, 0xff},
	color.RGBA{0x80, 0xc0, 0x00, 0xff},
	color.RGBA{0x90, 0xd0, 0x10, 0xff},
	color.RGBA{0xa0, 0xe0, 0x20, 0xff},
	color.RGBA{0xb0, 0xf0, 0x30, 0xff},
	color.RGBA{0xc0, 0x00, 0x40, 0xff},
	color.RGBA{0xd0, 0x10, 0x50, 0xff},
	color.RGBA{0xe0, 0x20, 0x60, 0xff},
	color.RGBA{0xf0, 0x30, 0x70, 0xff},
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	lissajous(os.Stdout)
}

func lissajous(out io.Writer) {
	const (
		cycles  = 5
		res     = 0.001
		size    = 100
		nframes = 64
		delay   = 8
	)
	freq := rand.Float64() * 3.0
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), uint8(i%16+1))
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}
