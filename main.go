package main

import (
	"image"
	"image/color"
	"image/png"
	_ "image/png"
	"io/ioutil"
	"log"
	"math"
	"os"
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
	screenWidth          = 512
	screenHeight         = 256
	bacgroundTextureSize = 16
	frameWidth           = 16
	frameHeight          = 32
	frameCount           = 6
	moveSpd              = 4.0
	dpi                  = 72
)

var (
	lvl1_map = [...]string{
		"b b b b b b b b b b b b b b b b",
		"b f f f f f f f f f f f f f f b",
		"b f f f f f f f f f f f f f f b",
		"b f f f b_1 b_2 f f f f f f f f f b",
		"b x f f f f f f f f f f f f f b",
		"b x f f f f f f f f f f f f f b",
		"b f f f f f f f f f f f f f f b",
		"b f f f f f f f f f f f f f f b",
		"b f f f f f f f f f f f f f f b",
		"b b b b b b f f f f b b b b b b",
		"b f f f f f f f f f f f f f f b",
		"b f f f f f f f f f f f f f f b",
		"b x f f f f f f f f f f f f f b",
		"b x f f f f f f f f f f f f f b",
		"b f f f v f f f f f s_1 f f f f b",
		"b f f f f f f f f f s_2 f f f f b",
		"b f f f f f f f f f f f f f f b",
		"b f f f f f f f f f f f f f f b",
		"b f f f f f f f f f f f f f f b",
		"b f f f f f f f f f f f f f f b",
		"b f f f f f f f f f f f f f f b",
		"b f f f f f f f f f f f f f f b",
		"b b b b b b f f f f b b b b b b",
		"b f f f f f f f f f f f f f x b",
		"b f f f f f f f f f f f f f x b",
		"b f f f f f f f f f f f f f f b",
		"b f f f f f f f f f f f f f f b",
		"b f f f f f f f f f f f f f f b",
		"b f f f f f f f f f f f f f f d",
		"b f f f f f f f f f f f f f f b",
		"b f f f f f f f f f f f f f f b",
		"b b b b b b b b b b b b b b b b",
	}
	lvl2_map = [...]string{
		"w w w w w w w w w w w w w w w w",
		"w a a a a a a a a a a a a a a w",
		"w a a a a a a a a a a a a a a w",
		"w a a a a a a a a a cr_1 cr_2 a a a w",
		"w a a a a a a a a a a a a a a w",
		"w a a a a a a a a a a a a a a w",
		"w a a a a_b g g g g g g g g g_b a w",
		"w a a a a_b g g g g g g g g g_b a w",
		"w a a a a_b g g g t t g g g g_b a w",
		"w a a a a_b g g g t t g g g g_b a w",
		"w a a a a_b g g g g g g g g g_b a w",
		"w a a a a_b g g g g g g g g g_b a w",
		"w a a a a_b g g g t t g g g g_b a w",
		"w a a a a_b g g t t t t g g g_b a w",
		"w a a a a_b g g t t t t g g g_b a w",
		"w a a a a_b g g g t t g g g g_b a w",
		"w a a a a_b g g g g g g g g g_b a w",
		"w a a a a_b g g g g g g g g g_b a w",
		"w a a a a_b g g g g g g g g g_b a w",
		"w a a a a_b g g g g g t g g g_b a w",
		"w a a a a_b g g g g t t t g g_b a w",
		"w a a a a_b g g g g t t t g g_b a w",
		"w a a a a_b g g g g g t g g g_b a w",
		"w a a a a_b g g g g g g g g g_b a w",
		"w a a a a_b g g g g g g g g g_b a w",
		"w a a a a_b g g g g g g g g g_b a w",
		"w a a a a_b g g g g g g g g g_b a w",
		"w a a a a a a a a a a a a a a w",
		"w a a a a a a a a a a a a a a w",
		"w a a cb_1 cb_2 a a a a a a a a a a w",
		"w a a a a a a a a a a_h a a a a w",
		"w w w w w w w w w w w w w w w w",
	}
)

type lvl_data struct {
	npces    []*npc_data
	lvl_map  []string
	enemies  []*enemy
	lvl_type string
	exitPosX int
	exitPosY int
}

var (
	playerImage    *ebiten.Image
	playerFace     *ebiten.Image
	bacgroundImage *ebiten.Image
	logo           []image.Image
)
var (
	lvl1_data = lvl_data{
		lvl_map:  lvl1_map[:],
		lvl_type: "room",
		exitPosX: 448,
		exitPosY: 240,
		npces: []*npc_data{
			&npc_data{
				startPosX:    200,
				startPosY:    130,
				sprite_asset: "Mama.png",
				face_asset:   "Mama_face.png",
				dialog_asset: "Mama_dialogs",
			},
			&npc_data{
				startPosX:    420,
				startPosY:    150,
				sprite_asset: "Vitalik.png",
				face_asset:   "Vitalik_face.png",
				dialog_asset: "Vitalik_dialogs",
			},
		},
	}
	lvl2_data = lvl_data{
		lvl_map:  lvl2_map[:],
		lvl_type: "room",
		exitPosX: 456,
		exitPosY: 160,
		npces: []*npc_data{
			&npc_data{
				startPosX:    100,
				startPosY:    32,
				sprite_asset: "Kostya.png",
				face_asset:   "Kostya_face.png",
				dialog_asset: "Kostya_dialogs",
			},
			&npc_data{
				startPosX:    200,
				startPosY:    32,
				sprite_asset: "Fil.png",
				face_asset:   "Fil_face.png",
				dialog_asset: "Fil_dialogs",
			},
			&npc_data{
				startPosX:    400,
				startPosY:    32,
				sprite_asset: "Artem.png",
				face_asset:   "Artem_face.png",
				dialog_asset: "Artem_dialogs",
			},
		},
	}
	lvl3_data = lvl_data{
		lvl_type: "rythm",
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

type Game struct {
	inputSystem input.System
	p           *player
	n           []*npc
	space       *resolv.Space
	lvl_map     []string
	curentLvl   int
	exitPosX    int
	exitPosY    int
	lvl_type    string
	count       int
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
type npc_data struct {
	startPosX    int
	startPosY    int
	sprite_asset string
	dialog_asset string
	face_asset   string
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
	"g":    Object{OX: 0, OY: 0, isObject: false},
	"g_b":  Object{OX: 0, OY: 16, isObject: true},
	"a":    Object{OX: 16, OY: 0, isObject: false},
	"a_b":  Object{OX: 16, OY: 16, isObject: true},
	"a_h":  Object{OX: 16 * 7, OY: 0, isObject: false},
	"t":    Object{OX: 16 * 2, OY: 0, isObject: true},
	"f":    Object{OX: 16 * 3, OY: 0, isObject: false},
	"w":    Object{OX: 16 * 4, OY: 0, isObject: true},
	"d":    Object{OX: 16 * 5, OY: 0, isObject: true},
	"b":    Object{OX: 16 * 6, OY: 0, isObject: true},
	"cb_1": Object{OX: 0, OY: 16 * 2, isObject: true},
	"cb_2": Object{OX: 0, OY: 16 * 3, isObject: true},
	"cr_1": Object{OX: 16, OY: 16 * 2, isObject: true},
	"cr_2": Object{OX: 16, OY: 16 * 3, isObject: true},
	"b_1":  Object{OX: 16 * 2, OY: 16 * 2, isObject: true},
	"b_2":  Object{OX: 16 * 2, OY: 16 * 3, isObject: true},
	"s_1":  Object{OX: 16 * 3, OY: 16 * 2, isObject: true},
	"s_2":  Object{OX: 16 * 4, OY: 16 * 2, isObject: true},
	"x":    Object{OX: 16 * 5, OY: 16 * 2, isObject: true},
	"v":    Object{OX: 16 * 6, OY: 16 * 2, isObject: true},
}

var lvls = map[int]*lvl_data{
	1: &lvl1_data,
	2: &lvl2_data,
	3: &lvl3_data,
}

func (g *Game) Update() error {
	lvlNumber := g.curentLvl
	if g.p.input.ActionIsJustPressed(ActionInteract) {
		for i := 0; i < len(g.n); i++ {
			if math.Abs(g.p.model.Position.X-g.n[i].model.Position.X) < 32 && math.Abs(g.p.model.Position.Y-g.n[i].model.Position.Y) < 32 {
				if g.n[i].isActive {
					if g.n[i].state < len(g.n[i].dialog)-1 {
						g.n[i].state++
					} else {
						g.n[i].state = 0
						g.n[i].isActive = false
						g.p.isLocked = false
					}
				} else {
					g.n[i].isActive = true
					g.p.isLocked = true
				}
			}
		}
		if math.Abs(g.p.model.Position.X-float64(g.exitPosX)) < 32 && math.Abs(g.p.model.Position.Y-float64(g.exitPosY)) < 32 {
			lvlNumber++
		}
	}
	g.p.count++
	g.count++
	g.inputSystem.Update()
	g.p.Update()
	for i := 0; i < len(g.n); i++ {
		g.n[i].Update()
		g.n[i].count++
	}
	if lvlNumber > g.curentLvl {
		g.curentLvl = lvlNumber
		g.ClearLevel()
		g.SetUpLevel(lvls[g.curentLvl])
		g.SetUpPlayer(lvls[g.curentLvl])
		if g.lvl_type == "room" {
			g.SetUpNpces(lvls[g.curentLvl])

		}
	}
	return nil
}

func startGame() *Game {
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
	g.curentLvl = 1
	g.SetUpLevel(&lvl1_data)
	g.SetUpPlayer(&lvl1_data)
	g.SetUpNpces(&lvl1_data)
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
	switch {
	case g.lvl_type == "room":
		g.DrawRoomBacground(screen)
		g.p.Draw(screen)
	case g.lvl_type == "rythm":
		g.DrawRythmBacground(screen)
	}

	for i := 0; i < len(g.n); i++ {
		g.n[i].Draw(screen)
	}
	for i := 0; i < len(g.n); i++ {
		if g.n[i].isActive {
			g.n[i].Dialog(screen, g.p.face)
		}
	}
}
func (g *Game) DrawRoomBacground(screen *ebiten.Image) {
	for i := 0; i < len(g.lvl_map); i++ {
		line := strings.Split(g.lvl_map[i], " ")
		for j := 0; j < len(line); j++ {
			op := &ebiten.DrawImageOptions{}
			sx, sy := mapObjects[line[j]].OX, mapObjects[line[j]].OY
			op.GeoM.Translate(float64(i*bacgroundTextureSize), float64(j*bacgroundTextureSize))
			screen.DrawImage(bacgroundImage.SubImage(image.Rect(sx, sy, sx+bacgroundTextureSize, sy+bacgroundTextureSize)).(*ebiten.Image), op)
		}
	}
}
func (g *Game) DrawRythmBacground(screen *ebiten.Image) {
	i := (g.count / 48) % 3
	ugolImage := Loader("ugol.png")
	filImage := Loader("Fil_tall.png")
	kostyaImage := Loader("Kostya_tall.png")
	artemImage := Loader("Artem_tall.png")
	vanoImage := Loader("Vano_tall.png")
	op_ugol := &ebiten.DrawImageOptions{}
	op_fil := &ebiten.DrawImageOptions{}
	op_kostya := &ebiten.DrawImageOptions{}
	op_artem := &ebiten.DrawImageOptions{}
	op_vano := &ebiten.DrawImageOptions{}
	op_fil.GeoM.Translate(128, 103)
	op_kostya.GeoM.Translate(192, 103)
	op_artem.GeoM.Translate(256, 103)
	op_vano.GeoM.Translate(320, 103)
	sx, sy := 0+i*64, 0
	screen.DrawImage(ugolImage, op_ugol)
	screen.DrawImage(filImage.SubImage(image.Rect(sx, sy, sx+64, sy+153)).(*ebiten.Image), op_fil)
	screen.DrawImage(kostyaImage.SubImage(image.Rect(sx, sy, sx+64, sy+153)).(*ebiten.Image), op_kostya)
	screen.DrawImage(artemImage.SubImage(image.Rect(sx, sy, sx+64, sy+153)).(*ebiten.Image), op_artem)
	screen.DrawImage(vanoImage.SubImage(image.Rect(sx, sy, sx+64, sy+153)).(*ebiten.Image), op_vano)
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

func (g *Game) ClearLevel() *Game {
	g.space = nil
	g.n = nil
	g.p = nil
	return g
}

func (g *Game) SetUpLevel(lvl_data *lvl_data) *Game {
	g.lvl_map = lvl_data.lvl_map
	g.lvl_type = lvl_data.lvl_type
	g.space = resolv.NewSpace(screenWidth, screenHeight, bacgroundTextureSize, bacgroundTextureSize)
	for i := 0; i < len(g.lvl_map); i++ {
		line := strings.Split(g.lvl_map[i], " ")
		for j := 0; j < len(line); j++ {
			if mapObjects[line[j]].isObject {
				g.space.Add(resolv.NewObject(float64(i*bacgroundTextureSize), float64(j*bacgroundTextureSize), bacgroundTextureSize, bacgroundTextureSize))
			}
		}
	}
	g.exitPosX, g.exitPosY = lvl_data.exitPosX, lvl_data.exitPosY
	return g
}

func (g *Game) SetUpNpces(lvl_data *lvl_data) *Game {
	for i := 0; i < len(lvl_data.npces); i++ {
		g.n = append(g.n, &npc{
			startPosX: lvl_data.npces[i].startPosX,
			startPosY: lvl_data.npces[i].startPosY,
			frameOX:   0,
			frameOY:   0,
			sprite:    Loader(lvl_data.npces[i].sprite_asset),
			dialog:    DialogLoader(lvl_data.npces[i].dialog_asset),
			face:      Loader(lvl_data.npces[i].face_asset),
			model:     resolv.NewObject(float64(lvl_data.npces[i].startPosX), float64(lvl_data.npces[i].startPosY+frameHeight/2), frameWidth, frameHeight/2),
		})
	}
	for i := 0; i < len(g.n); i++ {
		g.space.Add(g.n[i].model)
	}
	return g
}

func (g *Game) SetUpPlayer(lvl_data *lvl_data) *Game {
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
	g.space.Add(g.p.model)
	return g
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
	bacgroundImage = Loader("TexturePack.png")
	playerFace = Loader("Vano_face.png")
	logo_file, err := os.Open("_assets/Vano_face.png")
	if err != nil {
	}
	defer logo_file.Close()
	logo_img, _ := png.Decode(logo_file)
	logo = append(logo, logo_img)
	ebiten.SetWindowTitle("Simulyator Stoyaniya V Uglu")
	ebiten.SetWindowIcon(logo)
	if err := ebiten.RunGame(startGame()); err != nil {
		log.Fatal(err)
	}
}
