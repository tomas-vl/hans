package main

import (
	"image"
	"image/png"
	"os"
	"sort"

	"git.sr.ht/~sbinet/gg"
	"github.com/golang/freetype/truetype"
)

type Picture struct {
	filepath    string
	canvas      *gg.Context
	orig_width  float64
	orig_height float64
	border_size float64
	iX0         float64
	iY0         float64
	iX1         float64
	iY1         float64
}

type Pair struct {
	letter string
	pos    float64
}

func OpenImage(filepath string) image.Image {
	// Load the image data
	bitmap_as_file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	// Decode the PNG image data to an image.Image
	bitmap_as_image, err := png.Decode(bitmap_as_file)
	if err != nil {
		panic(err)
	}
	return bitmap_as_image
}

func InitializePicture(filepath string, width float64, height float64, border_size float64) Picture {
	canvas_width := width + (2 * border_size)
	canvas_height := height + (2 * border_size)
	canvas := gg.NewContext(int(canvas_width), int(canvas_height))

	picture := Picture{
		filepath:    filepath,
		canvas:      canvas,
		orig_width:  width,
		orig_height: height,
		border_size: border_size,
		iX0:         border_size,
		iY0:         border_size,
		iX1:         width,
		iY1:         height,
	}
	// načtení písma
	font_filepath := "font/CharisSIL/CharisSIL-Regular.ttf"
	font_bytes, err := os.ReadFile(font_filepath)
	if err != nil {
		panic(err)
	}

	font, err := truetype.Parse(font_bytes)
	if err != nil {
		panic(err)
	}
	font_size := picture.border_size - 15
	face := truetype.NewFace(font, &truetype.Options{Size: font_size})
	picture.canvas.SetFontFace(face)
	return picture
}

func DrawBitmap(pic Picture, bitmap image.Image) Picture {
	pic.canvas.DrawImage(bitmap, int(pic.iX0), int(pic.iY0))
	return pic
}

func DrawRectangle(pic Picture) Picture {
	pic.canvas.SetRGB(0, 0, 0)
	pic.canvas.SetLineWidth(3)
	pic.canvas.DrawRectangle(pic.iX0, pic.iY0, pic.iX1, pic.iY1)
	pic.canvas.Stroke()
	return pic
}

func DrawLabels(pic Picture, frequency string, duration string) Picture {
	// // střed obrázku
	// x_center := float64(pic.canvas.Width()) / 2
	// y_center := float64(pic.canvas.Height()) / 2

	// vytvoření popisků
	label_freq := "Frekvence 0–" + frequency + " Hz"
	label_duration := "Trvání " + duration + " s"

	// délka a pozice popisků
	width_of_label_duration, _ := pic.canvas.MeasureString(label_duration)
	x_duration := pic.orig_width - width_of_label_duration + pic.border_size
	y_duration := pic.border_size - 10

	width_of_label_freq, _ := pic.canvas.MeasureString(label_freq)
	x_freq := -pic.orig_height + width_of_label_freq/2 - pic.border_size
	y_freq := pic.border_size / 2

	// vkreslení popisku "Trvání X s"
	pic.canvas.DrawStringAnchored(label_duration, x_duration, y_duration, 0, 0)

	// vkreslení popisku "Frekvence 0–X Hz", který je otočený o 90° proti směru hodinových ručiček
	pic.canvas.Rotate(gg.Radians(-90))
	pic.canvas.DrawStringAnchored(label_freq, x_freq, y_freq, 0.5, 0.5)
	pic.canvas.Rotate(gg.Radians(90))
	return pic
}

func DrawLetter(pic Picture, x_position float64, letter string) Picture {
	y_position := pic.orig_height + pic.border_size + 20
	letter_width, _ := pic.canvas.MeasureString(letter)
	pic.canvas.SetRGB(0, 0, 0)
	pic.canvas.DrawStringWrapped(letter, x_position, y_position, 0.5, 0.5, letter_width, 0, 1)
	return pic
}

func DrawLine(pic Picture, position float64) Picture {
	y0 := pic.border_size
	y1 := pic.orig_height + pic.border_size
	pic.canvas.SetRGB(1, 0, 0)
	pic.canvas.SetLineWidth(3)
	pic.canvas.DrawLine(position, y0, position, y1)
	pic.canvas.Stroke()
	return pic
}

func GetResolution(img image.Image) (float64, float64) {
	bounds := img.Bounds()
	pic_width := float64(bounds.Dx())
	pic_height := float64(bounds.Dy())
	return pic_width, pic_height
}

func RecalculateLetterPositions(line_positions []float64, letter_positions []Pair, start float64, end float64) []Pair {
	temp_line_positions := make([]float64, len(line_positions)+2)
	copy(temp_line_positions, line_positions)
	temp_line_positions = append(temp_line_positions, start, end)
	sort.Float64s(temp_line_positions)

	new_letter_positions := make([]Pair, len(letter_positions))
	copy(new_letter_positions, letter_positions)

	for i := 0; i < len(letter_positions); i++ {
		letter_pos := letter_positions[i].pos

		for j := 0; j < (len(temp_line_positions) - 1); j++ {
			if temp_line_positions[j] <= letter_pos && letter_pos <= temp_line_positions[j+1] {
				average := (temp_line_positions[j] + temp_line_positions[j+1]) / 2
				new_letter_positions[i] = Pair{letter_positions[i].letter, average}
				break
			}
		}
	}
	return new_letter_positions
}
