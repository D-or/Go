/*
 * Revision History:
 *     Initial: 2018/05/24      Lin Hao
 */

package utils

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"golang.org/x/image/math/fixed"

	"github.com/astaxie/beego"
	"github.com/golang/freetype"
	"github.com/nfnt/resize"
	"golang.org/x/image/font"
)

var (
	dpi      = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
	fontfile = flag.String("fontfile", "src/font/black.ttf", "filename of the ttf font")
	hinting  = flag.String("hinting", "none", "none | full")
	size     = flag.Float64("size", 70, "font size in points")
	spacing  = flag.Float64("spacing", 1.2, "line spacing (e.g. 2 means double spaced)")
	wonb     = flag.Bool("whiteonblack", false, "white text on a black background")
)

const (
	iw = 700  // Image width
	ih = 700  // Image width
	fw = 1000 // File width
	fh = 1000 // File height
)

// Save the file uploaded.
func Save(r *http.Request) (string, error) {
	image, head, err := r.FormFile("image")
	if err != nil {
		beego.Error("Get image Error: ", err)

		return "", err
	}
	defer image.Close()

	imageType := []string{".png", ".jpg", ".jpeg"}
	var suffix string
	for _, v := range imageType {
		if strings.HasSuffix(head.Filename, v) {
			suffix = v

			break
		}
	}

	// Make Dir
	if err := os.MkdirAll("/root/doublewoodh/generation/src/images/origin", os.ModePerm); err != nil {
		beego.Error("Mkdir 'origin' Error: ", err)

		return "", err
	}

	// Create File
	fileName := Now() + suffix
	filePath := fmt.Sprintf("src/images/origin/" + fileName)

	file, err := os.Create(filePath)
	if err != nil {
		beego.Error("Create Error: ", err)

		return "", err
	}
	defer file.Close()

	// image.Seek(0, 0)
	_, err = io.Copy(file, image)
	if err != nil {
		beego.Error("Save Error: ", err)

		return "", err
	}

	return fileName, nil
}

// Generate an image with text.
func Generate(fileName string, texts []string, wordPosition string) {
	srcFile, err := os.Open("src/images/origin/" + fileName)
	if err != nil {
		beego.Error("Open Error: ", err)
	}
	defer srcFile.Close()

	var srcImg image.Image
	if strings.HasSuffix(fileName, ".png") {
		srcImg, err = png.Decode(srcFile)
	} else {
		srcImg, err = jpeg.Decode(srcFile)
	}
	if err != nil {
		beego.Error("Decode Error: ", err)
	}

	srcImg = cutWhiteSpace(srcImg)

	srcImg = resize.Resize(0, ih, srcImg, resize.Lanczos3)

	// Initialize the context.
	rgba := image.NewRGBA(image.Rect(0, 0, fw, fh))

	for i := 0; i < fw; i++ {
		for j := 0; j < fh; j++ {
			rgba.Set(i, j, image.White.C)
		}
	}

	// Draw the white image.
	imageY := 0
	if wordPosition == "top" {
		imageY = -300
	}
	draw.Draw(rgba, rgba.Bounds(), srcImg, image.Pt((srcImg.Bounds().Size().X-fw)/2, imageY), draw.Over)

	// Read the font data.
	fontBytes, err := ioutil.ReadFile(*fontfile)
	if err != nil {
		beego.Error("ReadFile Error: ", err)
		return
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		beego.Error("ParseFont Error: ", err)
		return
	}

	// Create a new Context.
	c := freetype.NewContext()
	// Set the screen resolution in dots per inch.
	c.SetDPI(*dpi)
	// Set Font.
	c.SetFont(f)
	// Set FontSize.
	fontsize := fontSize(texts)
	c.SetFontSize(fontsize)
	// Set the clip rectangle for drawing.
	c.SetClip(rgba.Bounds())
	// Set the destination image for draw operations.
	c.SetDst(rgba)
	// Set the source image for draw operations.
	c.SetSrc(image.Black)

	switch *hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	// Draw the text.
	wordY := 720
	if wordPosition == "top" {
		wordY = 20
	}
	pt := freetype.Pt(50, wordY+int(c.PointToFixed(*size)>>6))
	pt = changeY(c, pt, texts, fontsize)
	for _, v := range texts {
		reg := regexp.MustCompile(`[a-z|0-9|!|?|,|.|;]`)
		letter := reg.FindAllString(v, -1)

		chineseLen := len([]rune(v)) - len(letter)

		width := float64((1000 - (int(fontsize)*len(letter)/2 + int(fontsize)*chineseLen)) / 2)
		pt.X = c.PointToFixed(width)

		for _, s := range []string{v} {
			_, err = c.DrawString(s, pt)
			if err != nil {
				beego.Error("Draw Error: ", err)
				return
			}
		}

		pt.Y += c.PointToFixed(fontsize * *spacing)
	}

	// Make Dir
	if err := os.MkdirAll("/root/doublewoodh/generation/src/images/generated", os.ModePerm); err != nil {
		beego.Error("Mkdir 'generated' Error: ", err)
		os.Exit(1)
	}

	// Save the RGBA image to disk.
	outFile, err := os.Create("src/images/generated/" + fileName)
	if err != nil {
		beego.Error("Create Error: ", err)
		os.Exit(1)
	}
	defer outFile.Close()

	b := bufio.NewWriter(outFile)

	err = png.Encode(b, rgba)
	if err != nil {
		beego.Error("Encode Error: ", err)
		os.Exit(1)
	}

	err = b.Flush()
	if err != nil {
		beego.Error("Flush Error: ", err)
		os.Exit(1)
	}

	beego.Debug("Wrote out.png OK.")
}

func fontSize(texts []string) float64 {
	switch len(texts) {
	case 1:
		return *size + 35.0
	case 2:
		return *size + 25.0
	default:
		return *size
	}
}

func changeY(c *freetype.Context, pt fixed.Point26_6, texts []string, fontSize float64) fixed.Point26_6 {
	switch len(texts) {
	case 1:
		pt.Y += c.PointToFixed(*size * *spacing)
	case 2:
		pt.Y += c.PointToFixed(*size**spacing) / 2
	default:
	}

	return pt
}

func cutWhiteSpace(src image.Image) image.Image {
	size := src.Bounds().Size()
	var (
		left   int
		right  int
		top    int
		bottom int
	)

	for i := 1; i < size.X-1; i++ {
		for j := 1; j < size.Y-1; j++ {
			r, g, b, _ := src.At(i, j).RGBA()
			if r < 65000 && g < 65000 && b < 65000 {
				left = i

				goto top
			}
		}
	}

top:
	for j := 1; j < size.Y-1; j++ {
		for i := 1; i < size.X-1; i++ {
			r, g, b, _ := src.At(i, j).RGBA()
			if r < 65000 && g < 65000 && b < 65000 {
				top = j

				goto right
			}
		}
	}

right:
	for i := size.X - 1; i > 1; i-- {
		for j := size.Y - 1; j > 1; j-- {
			r, g, b, _ := src.At(i, j).RGBA()
			if r < 65000 && g < 65000 && b < 65000 {
				right = i

				goto bottom
			}
		}
	}

bottom:
	for j := size.Y - 1; j > 1; j-- {
		for i := 1; i < size.X-1; i++ {
			r, g, b, _ := src.At(i, j).RGBA()
			if r < 65000 && g < 65000 && b < 65000 {
				bottom = j

				goto finish
			}
		}
	}

finish:
	img := src.(*image.YCbCr)
	subImg := img.SubImage(image.Rect(left, top, right, bottom)).(*image.YCbCr)

	return subImg
}
