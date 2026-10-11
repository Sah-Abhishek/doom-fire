package main

import (
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	H = 168
	W = 320
)

type Game struct {
	fire   []uint8
	pixels []byte
}

func positionInList(x, y int) (int, int, int, int) {
	return ((y * W) + x) * 4, ((y*W)+x)*4 + 1, ((y*W)+x)*4 + 2, ((y*W)+x)*4 + 3
}

func (g *Game) setPixel(x, y int, red, green, blue byte) {
	r, gr, b, a := positionInList(x, y)
	g.pixels[r] = red
	g.pixels[gr] = green
	g.pixels[b] = blue
	g.pixels[a] = 0xFF
}

// spreadGradient is my experiment of rising gradient
func (g *Game) spreadGradient(x, y int) {
	from := y*W + x
	to := (y-1)*W + x

	heat := g.fire[from]
	if heat == 0 {
		g.fire[to] = 0
		return
	}
	g.fire[to] = uint8((y) * MaxHeat / (H - 1))
}

func (g *Game) spread(x, y int) {
	from := y*W + x
	heat := g.fire[from]

	if heat == 0 {
		g.fire[(y-1)*W+x] = 0
		return
	}
	drift := rand.Intn(3) - 1
	nx := (x + drift + W) % W
	if nx < 0 {
		nx = 0
	}
	if nx > W-1 {
		nx = W - 1
	}
	to := (y-1)*W + nx
	cool := uint8(rand.Intn(2))
	g.fire[to] = heat - cool
}

func (g *Game) Update() error {
	for y := 1; y < H; y++ {
		for x := 0; x < W; x++ {
			g.spread(x, y)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for x := 0; x < W; x++ {
		for y := 0; y < H; y++ {

			// What pixel it is
			i := y*W + x
			heat := g.fire[i]

			color := RedPalette[heat]

			g.pixels[i*4] = color[0]
			g.pixels[i*4+1] = color[1]
			g.pixels[i*4+2] = color[2]
			g.pixels[i*4+3] = 255
		}
	}
	screen.WritePixels(g.pixels)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return W, H
}

func main() {
	g := &Game{
		fire:   make([]uint8, W*H),
		pixels: make([]byte, W*H*4),
	}

	for x := 0; x < W; x++ {
		i := W*(H-1) + x
		g.fire[i] = uint8(MaxHeat)
	}

	ebiten.SetWindowSize(W*3, H*3)
	ebiten.SetWindowTitle("Doom Fire")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
