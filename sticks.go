package main

import (
	"log"
	"strconv"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"golang.org/x/image/font/gofont/gomedium"

	"image"
	"time"

	"math/rand"

	"github.com/fogleman/gg"
)

const G_Sticks int = 50
const G_TwoPlayersMode int = 1
const G_Space int = 50

var G_SticksCnt int = G_Sticks
var showTaken bool = false
var G_PlayerToMove int = 1
var timeStart int64 = 0
var Take int = 0

func play() {
	if G_SticksCnt%4 == 1 {
		if G_SticksCnt == 1 {
			Take = 1
			make_move(Take, G_PlayerToMove)
		} else {
			Rnd := 1 + rand.Intn(3)
			Take = Rnd
			make_move(Take, G_PlayerToMove)
		}
	} else {
		if G_SticksCnt%4 == 0 {
			Take = 3
			make_move(Take, G_PlayerToMove)
		} else {
			Take = (G_SticksCnt % 4) - 1
			make_move(Take, G_PlayerToMove)
		}
	}
	showTaken = true
	timeStart = time.Now().UnixNano()
}

func make_move(Get, Who int) {
	if Get <= G_SticksCnt {
		G_PlayerToMove = 3 - G_PlayerToMove
		G_SticksCnt = G_SticksCnt - Get
	}
}

func StickImage(showLastMove bool) image.Image {
	const (
		Width   = 1000
		Height  = 620
		CenterX = Width / 2
		CenterY = Height / 2
	)
	dc := gg.NewContext(Width, Height)

	dc.SetHexColor("#ffffff")
	dc.Clear()
	dc.Push()
	dc.SetHexColor("#000000")
	f, _ := truetype.Parse(gomedium.TTF)
	face := truetype.NewFace(f, &truetype.Options{Size: 20})
	dc.SetFontFace(face)
	dc.DrawString("Remaining: "+strconv.Itoa(G_SticksCnt), 50, 50)
	dc.RotateAbout(gg.Radians(-90), float64(CenterX), float64(CenterY))
	dc.Stroke()
	for i := G_SticksCnt; i > 0; i-- {
		dc.Push()
		dc.SetLineWidth(4)
		dc.MoveTo(700, float64((30 * i)))
		dc.LineTo(550, float64((30 * i)))
		dc.Stroke()
		dc.Pop()
	}
	var cnt int = G_SticksCnt + Take
	if showLastMove && G_SticksCnt > 0 {
		dc.SetHexColor("#ff0000")
		for i := Take; i > 0; i-- {
			dc.Push()
			dc.SetLineWidth(4)
			dc.MoveTo(700, float64((30 * cnt)))
			dc.LineTo(550, float64((30 * cnt)))
			dc.Stroke()
			dc.Pop()
			cnt -= 1
		}
	}

	if G_SticksCnt == 0 {
		dc.RotateAbout(gg.Radians(90), 500, 500)
		dc.DrawString("Player "+strconv.Itoa(G_PlayerToMove)+" wins", 500, 500)
	}
	return dc.Image()
}

type Game struct {
	stop bool
}

func (g *Game) Update(screen *ebiten.Image) error {
	button := 0
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		button = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		button = 2
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		button = 3
	}
	if button != 0 {
		make_move(button, G_PlayerToMove)
		if G_TwoPlayersMode > 0 {
			play()
		}
	}
	if showTaken && (time.Now().UnixNano()-timeStart) >= 1000000000 {
		showTaken = false
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	m := StickImage(showTaken)
	em, _ := ebiten.NewImageFromImage(m, ebiten.FilterDefault)
	screen.DrawImage(em, &ebiten.DrawImageOptions{})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1000, 320
}

func main() {
	game := &Game{}
	rand.Seed(time.Now().UnixNano())
	ebiten.SetWindowSize(1000, 320)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
