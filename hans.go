package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/gabstv/ebiten-imgui/renderer"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/inkyblackness/imgui-go/v4"
	camera "github.com/melonfunction/ebiten-camera"
	input "github.com/quasilyte/ebitengine-input"
	"github.com/sqweek/dialog"
)

type G struct {
	mgr *renderer.Manager

	drawing_type int
	duration     string
	freq         string
	filepath     string

	picture Picture

	cam *camera.Camera

	retina bool
	w, h   int

	inputHandler *input.Handler
	inputSystem  input.System
}

var (
	white = ebiten.NewImage(1, 1)
)

const (
	ActionUnknown input.Action = iota
	ActionScrollVertical
	ActionMouseClick
	zoom_speed = 0.1
)

func init() {
	white.Fill(color.White)
}

func (g *G) Draw(screen *ebiten.Image) {
	//	ebitenutil.DebugPrint(screen, msg)
	// bílé pozadí
	g.cam.Surface.Fill(color.White)

	ebiten_image := ebiten.NewImageFromImage(g.picture.canvas.Image())
	g.cam.Surface.DrawImage(ebiten_image, nil)

	// vykreslení GUI
	g.cam.Blit(screen)
	g.mgr.Draw(screen)
}

func (g *G) Update() error {
	// Začátek práce s GUI
	g.mgr.Update(1.0 / 60.0)
	g.mgr.BeginFrame()
	{
		default_window_size := imgui.Vec2{
			X: 280,
			Y: 150,
		}

		// Zde je první podokno. To se stará o vstup a výstup.
		imgui.SetNextWindowSize(default_window_size)
		imgui.Begin("INPUT & OUTPUT")
		{
			if imgui.Button("Select image") {
				var err error
				g.filepath, err = dialog.File().Filter("PNG image", "png").Load()
				if err != nil {
					fmt.Println("Chyba: Nebyl vybrán žádný soubor!")
				}
				fmt.Println("Načten soubor:", g.filepath)
				ebiten.SetWindowTitle(fmt.Sprintf("Hans — %s", g.filepath))

				//imgimg := OpenImage(g.filepath)
				input_pic := OpenImage(g.filepath)
				bounds := input_pic.Bounds()
				pic_width := float64(bounds.Dx())
				pic_height := float64(bounds.Dy())
				g.picture = InitializePicture(g.filepath, pic_width, pic_height, 45)
				g.picture = DrawBitmap(g.picture, input_pic)
				g.picture = DrawRectangle(g.picture)
			}
			imgui.SameLine()
			imgui.Text(g.filepath)

			imgui.Spacing()

			imgui.InputText("Frequency", &g.freq)
			imgui.InputText("Duration", &g.duration)

			imgui.Spacing()

			if imgui.Button("Save image") {
				output_filepath, err := dialog.File().Filter("PNG image", "png").Save()
				if err != nil {
					fmt.Println("Chyba: Nebylo zvoleno, kam uložit soubor!")
				}
				// tady zavolat funkci, co soubor doopravdy uloží
				fmt.Println("Uložen soubor:", output_filepath)
			}
		}

		imgui.End()

		// Tady je druhé podokno. V něm se vybírá, co se bude kreslit (čáry či písmena).
		imgui.SetNextWindowSize(default_window_size)
		imgui.Begin("DRAWING SELECTOR")
		{
			imgui.Text("Select what you want to draw:")
			imgui.RadioButtonInt("Line", &g.drawing_type, 0)
			imgui.RadioButtonInt("Letter", &g.drawing_type, 1)
		}
		imgui.End()

		imgui.Begin("DEBUG")
		{
			imgui.Text(fmt.Sprintf("Frekvence:\t%s", g.freq))
			imgui.Text(fmt.Sprintf("Trvání:\t%s", g.duration))
			imgui.Text(fmt.Sprintf("Vybraný spetrogram:\t%s", g.filepath))
			if imgui.Button("Print stuff") {
				fmt.Println("Frekvence:", g.freq)
				fmt.Println("Trvání:", g.duration)
				fmt.Println("Vybraný spetrogram", g.filepath)
			}
		}
		imgui.End()
	}
	g.mgr.EndFrame()
	// Konec práce s GUI

	// Ovládání
	// zoom při skrolování
	g.inputSystem.Update()
	if info, ok := g.inputHandler.JustPressedActionInfo(ActionScrollVertical); ok {
		if info.Pos.Y > 0 {
			g.cam.Zoom(1.1)
		} else if info.Pos.Y < 0 {
			g.cam.Zoom(0.9)
		}
	}

	// klikání
	return nil
}

func (g *G) Layout(outsideWidth, outsideHeight int) (int, int) {
	if g.retina {
		m := ebiten.DeviceScaleFactor()
		g.w = int(float64(outsideWidth) * m)
		g.h = int(float64(outsideHeight) * m)
	} else {
		g.w = outsideWidth
		g.h = outsideHeight
	}
	g.mgr.SetDisplaySize(float32(g.w), float32(g.h))
	g.cam.Resize(g.w, g.h)
	return g.w, g.h
}

func main() {
	ren := renderer.New(nil)

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Hans")

	cam := camera.NewCamera(1280, 720, 0, 0, 0, 1)
	cam.SetZoom(1.0)

	empty_picture := InitializePicture("", 1, 1, 0)
	empty_picture = DrawBitmap(empty_picture, image.NewRGBA(image.Rect(0, 0, 1, 1)))

	gg := &G{
		mgr:          ren,
		drawing_type: 0,
		duration:     "",
		filepath:     "No file selected",
		freq:         "",
		picture:      empty_picture,
		cam:          cam,
	}

	gg.inputSystem.Init(input.SystemConfig{
		DevicesEnabled: input.AnyDevice,
	})

	keymap := input.Keymap{
		ActionScrollVertical: {input.KeyWheelVertical},
		ActionMouseClick:     {input.KeyMouseLeft},
	}

	gg.inputHandler = gg.inputSystem.NewHandler(0, keymap)

	ebiten.RunGame(gg)
}
