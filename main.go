package main

import (
	"image"
	"image/color"
	_ "image/png"
	"io/ioutil"
	"log"
	"math"
	"strings"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	input "github.com/quasilyte/ebitengine-input"
	"github.com/solarlune/resolv"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	mplusNormalFont font.Face
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
	moveSpd              = 4.0
	dpi                  = 72
)

var (
	lvl1_map = [...]string{
		"w w w w w w w w w w w w w w w",
		"w 0 0 0 0 0 0 w 0 0 0 c c c w",
		"w 0 0 0 0 0 0 w 0 0 0 0 c c w",
		"w 0 0 0 0 0 0 w 0 0 0 0 c c w",
		"w 0 0 0 0 0 0 w 0 0 0 0 0 c w",
		"w 0 0 0 0 0 0 w 0 0 0 0 0 0 w",
		"w 0 0 0 0 0 0 w 0 0 0 0 0 0 w",
		"w 0 0 0 0 0 0 0 0 0 0 0 0 0 w",
		"w 0 0 0 0 0 0 0 0 0 0 0 0 0 w",
		"w 0 0 0 0 0 0 w 0 0 0 0 0 0 w",
		"w 0 0 c 0 0 0 w 0 0 0 0 0 0 w",
		"w 0 0 0 0 0 0 w 0 0 0 0 0 0 w",
		"w 0 0 0 0 0 0 w 0 0 0 0 0 0 w",
		"w 0 0 0 0 0 0 w 0 0 0 0 0 0 w",
		"w 0 0 0 0 0 0 w 0 0 0 0 0 0 w",
		"w 0 0 0 0 0 0 w 0 0 0 0 0 0 w",
		"w 0 0 0 0 0 0 w 0 0 0 0 0 0 w",
		"w 0 0 0 0 0 0 w 0 0 0 0 0 0 w",
		"w 0 0 0 0 0 0 w 0 0 0 0 0 0 w",
		"w w w w w w w w w w w w w w w",
	}
	lvl2_map = [...]string{
		"w w w w w w w w w w w w w w w",
		"w 0 0 0 0 0 0 w 0 0 0 0 0 c w",
		"w 0 0 0 0 0 0 w 0 0 0 0 c c w",
		"w 0 0 0 0 0 0 w 0 0 0 0 c c w",
		"w 0 0 0 0 0 0 w 0 0 0 0 0 c w",
		"w 0 0 0 0 0 0 w 0 0 0 0 0 0 w",
		"w 0 0 0 0 0 0 w 0 0 0 0 0 0 w",
		"w 0 0 0 0 0 0 0 0 0 0 0 0 0 w",
		"w 0 0 0 0 0 0 0 0 0 0 0 0 0 w",
		"w 0 0 0 0 0 0 w 0 0 0 0 0 0 w",
		"w 0 0 c 0 0 0 w 0 0 0 0 0 0 w",
		"w 0 0 0 0 0 0 w 0 0 0 0 c 0 w",
		"w 0 0 0 0 0 0 0 0 0 0 0 c 0 w",
		"w 0 0 0 0 0 0 0 0 0 0 0 c 0 w",
		"w 0 0 0 0 0 0 0 0 0 0 0 0 0 w",
		"w 0 0 0 0 0 0 w 0 0 0 0 c 0 w",
		"w 0 0 0 0 c 0 w 0 0 0 0 0 0 w",
		"w 0 0 0 0 0 0 w 0 0 0 0 0 0 w",
		"w 0 0 0 0 0 0 w 0 0 0 0 0 0 w",
		"w w w w w w w w w w w w w w w",
	}
)

type lvl_data struct {
	npces   []*npc
	lvl_map []string
	enemes  []*enemy
}

var (
	lvl1_data = lvl_data{
		lvl_map: lvl1_map[:],
		npces: []*npc{
			&npc{
				startPosX: 200,
				startPosY: 70,
				sprite:    playerImage,
				model:     resolv.NewObject(200, 70+frameHeight/2, frameWidth, frameHeight/2),
			},
			&npc{
				startPosX: 200,
				startPosY: 150,
				sprite:    playerImage,
				model:     resolv.NewObject(200, 150+frameHeight/2, frameWidth, frameHeight/2),
			},
		},
	}
	lvl2_data = lvl_data{
		lvl_map: lvl2_map[:],
		npces: []*npc{
			&npc{
				startPosX: 200,
				startPosY: 70,
				sprite:    playerImage,
				model:     resolv.NewObject(200, 70+frameHeight/2, frameWidth, frameHeight/2),
			},
			&npc{
				startPosX: 200,
				startPosY: 150,
				sprite:    playerImage,
				model:     resolv.NewObject(200, 150+frameHeight/2, frameWidth, frameHeight/2),
			},
		},
	}
)

const (
	ActionUnknown input.Action = iota
	ActionMoveLeft
	ActionMoveRight
	ActionMoveTop
	ActionMoveDown
	ActionUnbound
	ActionInteract
)

var (
	playerImage    *ebiten.Image
	playerFace     *ebiten.Image
	bacgroundImage *ebiten.Image
	npcImage       *ebiten.Image
	npcFace        *ebiten.Image
	Kostya_dialog  [][]string
)

type Game struct {
	inputSystem input.System
	p           *player
	n           *npc
	space       *resolv.Space
	lvl_map     []string
}

type player struct {
	startPosX int
	startPosY int
	input     *input.Handler
	sprite    *ebiten.Image
	face      *ebiten.Image
	model     *resolv.Object
	count     int
	frameOX   int
	frameOY   int
	isLocked  bool
}
type npc struct {
	startPosX int
	startPosY int
	model     *resolv.Object
	sprite    *ebiten.Image
	dialog    [][]string
	isActive  bool
	state     int
	count     int
	frameOX   int
	frameOY   int
	face      *ebiten.Image
}
type enemy struct {
	startPosX int
	startPosY int
	model     *resolv.Object
	sprite    *ebiten.Image
	count     int
	frameOX   int
	frameOY   int
}
type Object struct {
	OX       int
	OY       int
	isObject bool
}

var mapObjects = map[string]Object{
	"w": Object{OX: 16, OY: 3 * 16, isObject: true},
	"c": Object{OX: 16, OY: 14 * 16, isObject: true},
	"0": Object{OX: 2 * 16, OY: 15 * 16, isObject: false},
}

func (g *Game) Update() error {

	if g.p.input.ActionIsJustPressed(ActionInteract) {
		if math.Abs(g.p.model.Position.X-g.n.model.Position.X) < 32 && math.Abs(g.p.model.Position.Y-g.n.model.Position.Y) < 32 {
			if g.n.isActive {
				if g.n.state < len(g.n.dialog)-1 {
					g.n.state++

				} else {
					g.n.state = 0
					g.n.isActive = false
					g.p.isLocked = false
				}
			} else {
				g.n.isActive = true
				g.p.isLocked = true
			}
		}
	}
	g.p.count++
	g.n.count++
	g.inputSystem.Update()
	g.p.Update()
	g.n.Update()

	return nil
}

func roomLvl(lvl_data *lvl_data) *Game {
	g := &Game{}
	g.inputSystem.Init(input.SystemConfig{
		DevicesEnabled: input.AnyDevice,
	})
	tt, _ := opentype.Parse(fonts.MPlus1pRegular_ttf)

	mplusNormalFont, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    12,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})

	keymap := input.Keymap{
		ActionMoveLeft:  {input.KeyGamepadLeft, input.KeyLeft, input.KeyA},
		ActionMoveRight: {input.KeyGamepadRight, input.KeyRight, input.KeyD},
		ActionMoveTop:   {input.KeyGamepadLeft, input.KeyUp, input.KeyW},
		ActionMoveDown:  {input.KeyGamepadRight, input.KeyDown, input.KeyS},
		ActionInteract:  {input.KeyGamepadA, input.KeyE},
		ActionUnbound:   {},
	}
	g.p = &player{
		startPosX: 70,
		startPosY: 70,
		input:     g.inputSystem.NewHandler(0, keymap),
		frameOX:   0,
		frameOY:   0,
		sprite:    playerImage,
		face:      playerFace,
		model:     resolv.NewObject(70, 70+frameHeight/2, frameWidth, frameHeight/2),
	}
	g.n = &npc{
		startPosX: 170,
		startPosY: 170,
		frameOX:   0,
		frameOY:   0,
		sprite:    npcImage,
		dialog:    Kostya_dialog[:][:],
		face:      npcFace,
		model:     resolv.NewObject(170, 170+frameHeight/2, frameWidth, frameHeight/2),
	}
	//for i := 0; i < len(lvl_data.npces); i++ {
	//	g.n = &npc{
	//		startPosX: lvl_data.npces[i].startPosX,
	//		startPosY: lvl_data.npces[i].startPosY,
	//		frameOX:   0,
	//		frameOY:   0,
	//		sprite:    lvl_data.npces[i].sprite,
	//		model:     lvl_data.npces[i].model,
	//	}
	//}
	g.lvl_map = lvl_data.lvl_map
	g.space = resolv.NewSpace(screenWidth, screenHeight, bacgroundTextureSize, bacgroundTextureSize)
	g.space.Add(g.p.model)
	g.space.Add(g.n.model)

	for i := 0; i < len(g.lvl_map); i++ {
		line := strings.Split(g.lvl_map[i], " ")
		for j := 0; j < len(line); j++ {
			if mapObjects[line[j]].isObject {
				g.space.Add(resolv.NewObject(float64(i*bacgroundTextureSize), float64(j*bacgroundTextureSize), bacgroundTextureSize, bacgroundTextureSize))
			}
		}
	}
	return g
}

func Loader(path string) *ebiten.Image {
	asset, _, err := ebitenutil.NewImageFromFile("_assets/" + path)
	if err != nil {
		log.Fatal(err)
	}
	ebitenImage := ebiten.NewImageFromImage(asset)
	return ebitenImage
}

func DialogLoader(path string) [][]string {
	content, err := ioutil.ReadFile("_assets/" + path)
	if err != nil {
		log.Fatal(err)
	}
	res1 := strings.Split(string(content), "\n")
	var dialog [][]string
	for i := 0; i < len(res1); i++ {
		line := strings.Split(res1[i], ":")
		dialog = append(dialog, [][]string{line}...)
	}
	return dialog
}
func (g *Game) Draw(screen *ebiten.Image) {
	for i := 0; i < len(g.lvl_map); i++ {
		line := strings.Split(g.lvl_map[i], " ")
		for j := 0; j < len(line); j++ {
			op := &ebiten.DrawImageOptions{}
			sx, sy := mapObjects[line[j]].OX, mapObjects[line[j]].OY
			op.GeoM.Translate(float64(i*bacgroundTextureSize), float64(j*bacgroundTextureSize))
			screen.DrawImage(bacgroundImage.SubImage(image.Rect(sx, sy, sx+bacgroundTextureSize, sy+bacgroundTextureSize)).(*ebiten.Image), op)
		}
	}
	g.p.Draw(screen)
	g.n.Draw(screen)
	if g.n.isActive {
		g.n.Dialog(screen, g.p.face)
	}

}

func (p *player) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.model.Position.X), float64(p.model.Position.Y-frameHeight/2))
	i := (p.count / 12) % frameCount
	sx, sy := p.frameOX+i*frameWidth, p.frameOY
	screen.DrawImage(p.sprite.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
}

func (n *npc) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(n.model.Position.X), float64(n.model.Position.Y-frameHeight/2))
	i := (n.count / 12) % frameCount
	sx, sy := n.frameOX+i*frameWidth, n.frameOY
	screen.DrawImage(n.sprite.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (n *npc) Dialog(screen, face *ebiten.Image) {
	textBoxBoeder := ebiten.NewImage(screenWidth, screenHeight/3)
	textBox := ebiten.NewImage(screenWidth-8, screenHeight/3-8)
	textBoxBoeder.Fill(color.RGBA{0, 0, 0, 255})
	textBox.Fill(color.RGBA{255, 255, 255, 255})
	textBoxBoeder_op := &ebiten.DrawImageOptions{}
	textBoxBoeder_op.GeoM.Translate(0, screenHeight*0.75)
	textBox_op := &ebiten.DrawImageOptions{}
	textBox_op.GeoM.Translate(4, screenHeight*0.75+4)
	face_op := &ebiten.DrawImageOptions{}
	screen.DrawImage(textBoxBoeder, textBoxBoeder_op)
	screen.DrawImage(textBox, textBox_op)

	if n.dialog[n.state][0] == "n" {
		face_op.GeoM.Translate(screenWidth-128, screenHeight-128)
		screen.DrawImage(n.face, face_op)
		text.Draw(screen, n.dialog[n.state][1], mplusNormalFont, 8, screenHeight*0.75+24, color.Black)
	} else {
		face_op.GeoM.Translate(0, screenHeight-128)
		screen.DrawImage(face, face_op)
		text.Draw(screen, n.dialog[n.state][1], mplusNormalFont, 8+128, screenHeight*0.75+24, color.Black)
	}
}

func (p *player) Update() {
	p.frameOY = 0
	p.frameOX = 0
	dx, dy := 0.0, 0.0
	moveSpd := 4.0

	if p.input.ActionIsPressed(ActionMoveLeft) {
		dx = -moveSpd
		p.frameOY = 32
		p.frameOX = 32 * 6
	}
	if p.input.ActionIsPressed(ActionMoveRight) {
		dx += moveSpd
		p.frameOY = 32
		p.frameOX = 0
	}
	if p.input.ActionIsPressed(ActionMoveDown) {
		dy += moveSpd
		p.frameOY = 32
		p.frameOX = 32 * 9
	}
	if p.input.ActionIsPressed(ActionMoveTop) {
		dy = -moveSpd
		p.frameOY = 32
		p.frameOX = 32 * 3
	}
	if collision := p.model.Check(dx, 0); collision != nil {
		dx = 0
	}
	if collision := p.model.Check(0, dy); collision != nil {
		dy = 0
	}
	if p.isLocked {
		dx, dy = 0, 0
		p.frameOY = 0
		p.frameOX = 0
	}
	p.model.Position.X += dx
	p.model.Position.Y += dy
	p.model.Update()
}
func (n *npc) Update() {
}
func main() {
	// Decode an image from the image file's byte slice.
	playerImage = Loader("Vano.png")
	bacgroundImage = Loader("Textures.png")
	npcImage = Loader("Kostya.png")
	npcFace = Loader("Kostya_face.png")
	playerFace = Loader("Vano_face.png")
	Kostya_dialog = DialogLoader("Kostya_dialogs")
	ebiten.SetWindowTitle("Simulyator Stoyaniya V Uglu")
	if err := ebiten.RunGame(roomLvl(&lvl1_data)); err != nil {
		log.Fatal(err)
	}
}
