package main

import (
	"fmt"
	"image/color"

	"github.com/gabstv/ebiten-imgui/renderer"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/inkyblackness/imgui-go/v4"
	"github.com/sqweek/dialog"
)

type G struct {
	mgr          *renderer.Manager
	drawing_type int
	duration     string
	filepath     string
	freq         string
	retina       bool
	w, h         int
}

func (g *G) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{40, 40, 40, 255})
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %.2f", ebiten.ActualTPS()))
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
					panic(err)
				}
				fmt.Println("Načten soubor:", g.filepath)
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
					panic(err)
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
	}
	g.mgr.EndFrame()
	// Konec práce s GUI

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
	return g.w, g.h
}

func main() {
	ren := renderer.New(nil)

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Hans")

	gg := &G{
		mgr:          ren,
		drawing_type: 0,
		duration:     "",
		filepath:     "No file selected",
		freq:         "",
	}

	ebiten.RunGame(gg)
}
