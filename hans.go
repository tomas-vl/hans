package main

import (
	"fmt"
	"image"
	"image/color"

	imgui "github.com/gabstv/cimgui-go"
	ebimgui "github.com/gabstv/ebiten-imgui/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/jakecoffman/cp"
	camera "github.com/melonfunction/ebiten-camera"
	"github.com/sqweek/dialog"
)

type G struct {
	drawing_type bool
	duration     string
	freq         string
	filepath     string
	letter       string

	picture      Picture
	ebiten_image *ebiten.Image

	cam *camera.Camera

	retina bool
	w, h   int
}

var (
	lines   []float64 = []float64{}
	letters []Pair    = []Pair{}
	white             = ebiten.NewImage(1, 1)
	zoom              = 1.0
	font1   *imgui.Font
)

const (
	zoom_min   = 0.25
	zoom_max   = 10.0
	zoom_speed = 0.1
)

func init() {
	white.Fill(color.White)
}

func (g *G) Draw(screen *ebiten.Image) {
	// bílé pozadí
	g.cam.Surface.Clear()
	g.cam.Surface.Fill(color.White)

	g.cam.Surface.DrawImage(g.ebiten_image, nil)

	// vykreslení GUI
	g.cam.Blit(screen)
	ebimgui.Draw(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
}

func (g *G) Update() error {
	// Začátek práce s GUI
	ebimgui.Update(1.0 / 60.0)
	ebimgui.BeginFrame()
	{
		imgui.PushFont(font1)
		// Zde je první podokno. To se stará o vstup a výstup.
		imgui.Begin("INPUT & OUTPUT")
		{
			if imgui.Button("Select image") {
				var err error
				// načtení cesty k obrázku
				g.filepath, err = dialog.File().Filter("PNG image", "png").Load()
				if err != nil {
					fmt.Println("Chyba: Nebyl vybrán žádný soubor!")
				}
				fmt.Println("Načten soubor:", g.filepath)
				ebiten.SetWindowTitle(fmt.Sprintf("Hans — %s", g.filepath))

				// vytvoření struktury, kterou se bude kreslit
				input_pic := OpenImage(g.filepath)
				bounds := input_pic.Bounds()
				pic_width := float64(bounds.Dx())
				pic_height := float64(bounds.Dy())
				g.picture = InitializePicture(g.filepath, pic_width, pic_height, 45)
				// vkreslení obrázku a okraje
				g.picture = DrawBitmap(g.picture, input_pic)
				g.picture = DrawRectangle(g.picture)

				g.ebiten_image = ebiten.NewImageFromImage(g.picture.canvas.Image())

				// resetování zoomu
				g.cam.SetZoom(1.0)

				// resetování výstupních čar a písmen
				lines = nil
				letters = nil
			}
			imgui.SameLine()
			imgui.Text(g.filepath)

			imgui.Spacing()

			imgui.InputTextWithHint("Frequency", "", &g.freq, 0, nil)
			imgui.InputTextWithHint("Duration", "", &g.duration, 0, nil)
			imgui.InputTextWithHint("Letter", "", &g.letter, 0, nil)

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
		imgui.Begin("DRAWING SELECTOR")
		{
			imgui.Text("Select what you want to draw:")
			if imgui.RadioButtonBool("Line", !g.drawing_type) {
				g.drawing_type = false
			}
			if imgui.RadioButtonBool("Letter", g.drawing_type) {
				g.drawing_type = true
			}
		}
		imgui.End()
		// imgui.Begin("DEBUG")
		// {
		// 	imgui.Text(fmt.Sprintf("Frekvence:\t%s", g.freq))
		// 	imgui.Text(fmt.Sprintf("Trvání:\t%s", g.duration))
		// 	imgui.Text(fmt.Sprintf("Vybraný spetrogram:\t%s", g.filepath))
		// 	imgui.Text(fmt.Sprintf("Zoom:\t%f", g.cam.Scale))
		// 	imgui.Text(fmt.Sprintf("Letters:\tv", letters))
		// 	if imgui.Button("Print stuff") {
		// 		fmt.Println("Frekvence:", g.freq)
		// 		fmt.Println("Trvání:", g.duration)
		// 		fmt.Println("Vybraný spetrogram", g.filepath)
		// 		fmt.Println("Zoom", g.cam.Scale)

		// 		//letters = RecalculateLetterPositions(lines, letters, g.picture.border_size, float64(g.picture.canvas.Width())-g.picture.border_size)
		// 		// fmt.Println(letters)
		// 	}
		// }
		// imgui.End()
	}
	imgui.PopFont()
	ebimgui.EndFrame()
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

	// vykreslení čáry či písmena při kliknutí
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) && !imgui.CurrentIO().WantCaptureMouse() {
		cx, _ := ebiten.CursorPosition()
		ccx := float64(cx) / zoom

		if !g.drawing_type { // kreslení čar; g.drawing_type == false
			line_pos := ccx
			g.picture = DrawLine(g.picture, ccx)
			lines = append(lines, line_pos)

		} else if g.drawing_type { // kreslení písmen; g.drawing_type == true
			letter_pos := ccx
			g.picture = DrawLetter(g.picture, letter_pos, g.letter)
			letters = append(letters, Pair{g.letter, letter_pos})
		}
		g.ebiten_image = ebiten.NewImageFromImage(g.picture.canvas.Image())
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
	ebimgui.SetDisplaySize(float32(g.w), float32(g.h))
	g.cam.Resize(g.w, g.h)
	return g.w, g.h
}

func main() {

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Hans")

	cam := camera.NewCamera(1280, 720, 0, 0, 0, 1)
	cam.SetZoom(1.0)
	empty_picture := InitializePicture("", 1, 1, 0)
	empty_picture = DrawBitmap(empty_picture, image.NewRGBA(image.Rect(0, 0, 1, 1)))

	// generování vlastního "GlyphRange" se zdá být rozbité
	// var glyphs imgui.GlyphRange = imgui.NewGlyphRange()
	// ranges_builder := imgui.NewFontGlyphRangesBuilder()
	// ranges_builder.AddText("Frequency ")
	// ranges_builder.BuildRanges(glyphs)

	// cfg0 := imgui.NewFontConfig()

	// font1 = imgui.CurrentIO().Fonts().AddFontFromFileTTFV("font/NotoSans/NotoSans-Regular.ttf", 16, cfg0, glyphs.Data())

	font1 = imgui.CurrentIO().Fonts().AddFontFromFileTTF("font/NotoSans/NotoSans-Regular.ttf", 16.0)
	imgui.CurrentIO().Fonts().Build()

	gg := &G{
		drawing_type: false,
		duration:     "",
		filepath:     "No file selected",
		freq:         "",
		letter:       "",
		picture:      empty_picture,
		cam:          cam,
		ebiten_image: ebiten.NewImage(1, 1),
	}

	ebiten.RunGame(gg)
}
