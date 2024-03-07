package main

import (
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	input "github.com/quasilyte/ebitengine-input"
)

const (
	screenWidth          = 320
	screenHeight         = 240
	backgroundOY         = 16 * 15
	backgroundOX         = 16 * 2
	bacgroundTextureSize = 16
	frameWidth           = 16
	frameHeight          = 32
	frameCount           = 6
)
const (
	ActionUnknown input.Action = iota
	ActionMoveLeft
	ActionMoveRight
	ActionMoveTop
	ActionMoveDown
	ActionUnbound
)

var (
	runnerImage    *ebiten.Image
	bacgroundImage *ebiten.Image
)

type Game struct {
	inputSystem input.System
	p           *player
}
type player struct {
	input   *input.Handler
	pos     image.Point
	count   int
	frameOX int
	frameOY int
}

func (g *Game) Update() error {
	g.p.count++
	g.inputSystem.Update()
	g.p.Update()
	return nil
}

func newExampleGame() *Game {
	g := &Game{}
	g.inputSystem.Init(input.SystemConfig{
		DevicesEnabled: input.AnyDevice,
	})
	keymap := input.Keymap{
		ActionMoveLeft:  {input.KeyGamepadLeft, input.KeyLeft, input.KeyA},
		ActionMoveRight: {input.KeyGamepadRight, input.KeyRight, input.KeyD},
		ActionMoveTop:   {input.KeyGamepadLeft, input.KeyUp, input.KeyW},
		ActionMoveDown:  {input.KeyGamepadRight, input.KeyDown, input.KeyS},
		ActionUnbound:   {},
	}
	g.p = &player{
		input:   g.inputSystem.NewHandler(0, keymap),
		pos:     image.Point{X: 96, Y: 96},
		frameOX: 0,
		frameOY: 0,
	}

	return g
}
func (g *Game) Draw(screen *ebiten.Image) {
	for i := 0; i < screenWidth/bacgroundTextureSize; i++ {
		for j := 0; j < screenHeight/bacgroundTextureSize; j++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(i*bacgroundTextureSize), float64(j*bacgroundTextureSize))
			//op.GeoM.Translate(-float64(bacgroundTextureSize)/2, -float64(bacgroundTextureSize)/2)
			screen.DrawImage(bacgroundImage.SubImage(image.Rect(backgroundOX, backgroundOY, backgroundOX+bacgroundTextureSize, backgroundOY+bacgroundTextureSize)).(*ebiten.Image), op)
		}
	}
	g.p.Draw(screen)
}
func (p *player) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	op.GeoM.Translate(float64(p.pos.X), float64(p.pos.Y))
	i := (p.count / 6) % frameCount
	sx, sy := p.frameOX+i*frameWidth, p.frameOY
	screen.DrawImage(runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
func (p *player) Update() {
	p.frameOY = 0
	p.frameOX = 0
	if p.input.ActionIsPressed(ActionMoveLeft) {
		p.pos.X -= 4
		p.frameOY = 32
		p.frameOX = 32 * 6
	}
	if p.input.ActionIsPressed(ActionMoveRight) {
		p.pos.X += 4
		p.frameOY = 32
		p.frameOX = 0
	}
	if p.input.ActionIsPressed(ActionMoveDown) {
		p.pos.Y += 4
		p.frameOY = 32
		p.frameOX = 32 * 9
	}
	if p.input.ActionIsPressed(ActionMoveTop) {
		p.pos.Y -= 4
		p.frameOY = 32
		p.frameOX = 32 * 3
	}

}
func main() {
	// Decode an image from the image file's byte slice.
	pers, _, err := ebitenutil.NewImageFromFile("_assets/Adam.png")
	background, _, err := ebitenutil.NewImageFromFile("_assets/Textures.png")
	if err != nil {
		log.Fatal(err)
	}
	runnerImage = ebiten.NewImageFromImage(pers)
	bacgroundImage = ebiten.NewImageFromImage(background)
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Animation (Ebitengine Demo)")
	if err := ebiten.RunGame(newExampleGame()); err != nil {
		log.Fatal(err)
	}
}
