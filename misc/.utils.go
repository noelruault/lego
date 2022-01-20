package main

import (
	"context"
	"database/sql"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	pkgcolors "gopkg.in/go-playground/colors.v1"
)

func manualDecode() {
	ctx := context.Background()
	dbs := dbConnect("lego_artboard")
	imagename := "./image.png"

	// Process image to get its pixels
	pixels := processImage(ctx, dbs, "./image.png", "png")
	log.Println("Image processed gracefully.")

	x, y := getMatrixLength(pixels)
	log.Printf("Image matrix length: %dx%d", x, y)

	sqlStatement := `
		INSERT INTO artboards (imagename, pixels)
		VALUES ($1, $2)
	`
	_, err := dbs.Exec(sqlStatement, imagename, fmt.Sprintf("%v", pixels))
	if err != nil {
		panic(err)
	}

	// The image package is based on interfaces. Just define a new type with those methods.
	// Your type's ColorModel would return color.RGBAModel, Bounds - your rectangle's borders,
	// and At - the color at (x, y) that you can compute if you know the image's dimensions.

	imageRGBA := setImagePixels(x, y, pixels)

	// outputFile is a File type which satisfies Writer interface
	outputFile, err := os.Create("test.png")
	if err != nil {
		panic(err)
	}

	// Encode takes a writer interface and an image interface
	// We pass it the File and the RGBA
	png.Encode(outputFile, imageRGBA)

	// Don't forget to close files
	outputFile.Close()
}

// Get the bi-dimensional pixel array
// https://stackoverflow.com/a/41185404
func getPixels(file io.Reader) ([][]Pixel, error) {
	img, _, err := image.Decode(file)

	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels [][]Pixel
	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}

	return pixels, nil
}

// img.At(x, y).RGBA() returns four uint32 values; we want a Pixel
func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

func setImagePixels(x, y int, pixels [][]Pixel) *image.RGBA {
	rgba := image.NewRGBA(
		image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: x, Y: y},
		},
	)

	for i := 0; i < len(pixels); i++ {
		for j := 0; j < len(pixels[i]); j++ {
			rgba.SetRGBA(i, j, color.RGBA{
				R: uint8(pixels[i][j].R),
				G: uint8(pixels[i][j].G),
				B: uint8(pixels[i][j].B),
			})
		}
	}

	return rgba
}

type Hex string

type RGB struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

func (h Hex) toRGB() (RGB, error) {
	return Hex2RGB(h)
}

func Hex2RGB(hex Hex) (RGB, error) {
	var rgb RGB
	values, err := strconv.ParseUint(string(hex), 16, 32)

	if err != nil {
		return RGB{}, err
	}

	rgb = RGB{
		Red:   uint8(values >> 16),
		Green: uint8((values >> 8) & 0xFF),
		Blue:  uint8(values & 0xFF),
	}

	return rgb, nil
}

func processImage(ctx context.Context, dbs *sql.DB, imagepath, imagetype string) [][]Pixel {
	// You can register another format here
	image.RegisterFormat(imagetype, imagetype, png.Decode, png.DecodeConfig)
	file, err := os.Open(imagepath)
	if err != nil {
		fmt.Println("Error: File could not be opened")
		os.Exit(1)
	}
	defer file.Close()

	pixels, err := getPixels(file)

	if err != nil {
		fmt.Println("Error: Image could not be decoded")
		os.Exit(1)
	}

	legocolors := getLegoColors(dbs)

	var ipixels [][]Pixel
	var jpixels []Pixel
	for i := 0; i < len(pixels); i++ {
		jpixels = make([]Pixel, 0) // clean

		for j := 0; j < len(pixels[i]); j++ {
			var maxdistance, mindistance float64
			var maxhexd, minhexd string
			var maxhexid, minhexid int
			var minpixel Pixel

			cof, err := pkgcolors.RGB(uint8(pixels[i][j].R), uint8(pixels[i][j].G), uint8(pixels[i][j].B))
			if err != nil {
				panic(err)
			}
			hexfoundstring := strings.Trim(cof.ToHEX().String(), "#")

			// for every original pixel color, find nearest and farest lego color
			for k := range legocolors {
				distance := calculateDistance(pixels[i][j], Pixel{R: legocolors[k].R, G: legocolors[k].G, B: legocolors[k].B})
				if k == 0 {
					mindistance, maxdistance = math.MaxFloat64, 0
				}

				co, err := pkgcolors.RGB(uint8(legocolors[k].R), uint8(legocolors[k].G), uint8(legocolors[k].B))
				if err != nil {
					panic(err)
				}

				hexstring := strings.Trim(co.ToHEX().String(), "#")

				if distance > maxdistance {
					maxdistance = distance
					maxhexd = hexstring
					maxhexid = legocolors[i].ID
				}
				if distance < mindistance {
					mindistance = distance
					minhexd = hexstring
					minpixel = Pixel{R: legocolors[k].R, G: legocolors[k].G, B: legocolors[k].B}
					minhexid = legocolors[i].ID
				}
			}

			// for every original pixel color, update nearest and farest lego color
			sqlStatement := `
				INSERT INTO seen (
					hex, minlegoid, hexmindistance, maxlegoid, hexmaxdistance
				)
				VALUES ($1, $2, $3, $4, $5)
				ON CONFLICT DO NOTHING
			`

			_, err = dbs.Exec(sqlStatement, hexfoundstring, minhexid, minhexd, maxhexid, maxhexd)
			if err != nil {
				panic(err)
			}

			// print closest lego color, continue
			jpixels = append(jpixels, minpixel)
		}

		ipixels = append(ipixels, jpixels)
	}

	return ipixels
}
func translateColors() {
	dbs := dbConnect("lego_artboard")
	rows, err := dbs.Query(`
			SELECT
				id,
				r,
				g,
				b
			FROM
				colors
		`)
	if err != nil {
		log.Fatal(err)
	}

	var colors []Color
	for rows.Next() {
		var id int
		var r, g, b int
		err := rows.Scan(&id, &r, &g, &b)
		if err != nil {
			panic(err)
		}

		colors = append(colors, Color{
			ID: id,
			R:  r,
			G:  g,
			B:  b,
		})
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	for i := range colors {

		co, err := pkgcolors.RGB(uint8(colors[i].R), uint8(colors[i].G), uint8(colors[i].B))
		if err != nil {
			panic(err)
		}

		sqlStatement := fmt.Sprintf(`
				UPDATE colors SET
					hexmatch = %s
				WHERE id = %d
				`, strings.Trim(co.ToHEX().String(), "#"), colors[i].ID)

		_, err = dbs.Exec(sqlStatement)
		if err != nil {
			panic(err)
		}
	}
}

func getMatrixLength(matrix [][]Pixel) (y, x int) {
	var ielements, jelements int
	for i := 0; i < len(matrix); i++ {
		ielements += 1
		jelements = 0
		for j := 0; j < len(matrix[i]); j++ {
			jelements += 1
		}
	}
	return ielements, jelements
}
