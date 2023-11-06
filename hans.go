package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/gabstv/ebiten-imgui/renderer"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/inkyblackness/imgui-go/v4"
	"github.com/jakecoffman/cp"
	camera "github.com/melonfunction/ebiten-camera"
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
}

var (
	lines   []float64 = []float64{}
	letters []Pair    = []Pair{}
	white             = ebiten.NewImage(1, 1)
	zoom              = 1.0
)

const (
	zoom_min   = 0.5
	zoom_max   = 5.0
	zoom_speed = 0.1
)

func init() {
	white.Fill(color.White)
}

func (g *G) Draw(screen *ebiten.Image) {
	//	ebitenutil.DebugPrint(screen, msg)
	// bílé pozadí
	g.cam.Surface.Clear()
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

				// resetování zoomu
				g.cam.SetZoom(1.0)
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

				sav_input_pic := OpenImage(g.filepath)
				sav_bounds := sav_input_pic.Bounds()
				sav_pic_width := float64(sav_bounds.Dx())
				sav_pic_height := float64(sav_bounds.Dy())
				sav_output_picture := InitializePicture(g.filepath, sav_pic_width, sav_pic_height, 45)
				sav_output_picture = DrawBitmap(sav_output_picture, sav_input_pic)

				for i := 0; i < len(lines); i++ {
					sav_output_picture = DrawLine(sav_output_picture, lines[i])
				}

				letters = RecalculateLetterPositions(lines, letters, g.picture.border_size, float64(g.picture.canvas.Width())-g.picture.border_size)
				for i := 0; i < len(letters); i++ {
					sav_output_picture = DrawLetter(sav_output_picture, letters[i].pos, letters[i].letter)
				}

				sav_output_picture = DrawRectangle(sav_output_picture)
				sav_output_picture = DrawLabels(sav_output_picture, g.freq, g.duration)

				sav_output_picture.canvas.SavePNG(output_filepath)
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
			imgui.Text(fmt.Sprintf("Zoom:\t%f", g.cam.Scale))
			if imgui.Button("Print stuff") {
				fmt.Println("Frekvence:", g.freq)
				fmt.Println("Trvání:", g.duration)
				fmt.Println("Vybraný spetrogram", g.filepath)
				fmt.Println("Zoom", g.cam.Scale)

				letters = RecalculateLetterPositions(lines, letters, g.picture.border_size, float64(g.picture.canvas.Width())-g.picture.border_size)
				fmt.Println(letters)
			}
		}
		imgui.End()
	}
	g.mgr.EndFrame()
	// Konec práce s GUI

	// Ovládání
	// zoom při skrolování
	_, scroll_amount := ebiten.Wheel()
	if scroll_amount > 0 {
		zoom += zoom_speed
	} else if scroll_amount < 0 {
		zoom -= zoom_speed
	}
	zoom = cp.Clamp(zoom, zoom_min, zoom_max)
	g.cam.SetZoom(zoom)

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
		cx, cy := ebiten.CursorPosition()
		ccx, ccy := float64(cx)/zoom, float64(cy)/zoom
		wx, wy := g.cam.GetCursorCoords()
		fmt.Printf("X %d\tY %d\n", cx, cy)
		fmt.Printf("Cursor coords:\tX %f\tY %f\n", ccx, ccy)
		fmt.Printf("World coords:\tX %f\tY %f\n", wx, wy)

		if g.drawing_type == 0 {
			line_pos := ccx
			g.picture = DrawLine(g.picture, ccx)
			lines = append(lines, line_pos)
			fmt.Println(lines)
		} else if g.drawing_type == 1 {
			letter_pos := ccx
			letter := "a"
			g.picture = DrawLetter(g.picture, letter_pos, letter)
			letters = append(letters, Pair{"a", letter_pos})
			//letters = RecalculateLetterPositions(lines, letters, g.picture.border_size, float64(g.picture.canvas.Width())-g.picture.border_size)
			fmt.Println(letters)
		}
	}

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

	ebiten.RunGame(gg)
}
