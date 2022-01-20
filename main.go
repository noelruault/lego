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
	"math/rand"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/image/draw"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	pieces      int
	colors      int
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	outfolder   = "./assets/out/"
	outimages   = outfolder + "images/"
)

func dbConnect(dbname string) *sql.DB {
	const (
		dbhost = "localhost"
		dbport = 5432
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d dbname=%s sslmode=disable",
		dbhost, dbport, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type LegoColor struct {
	Hex    string
	ID     int
	LegoID int
	Name   string
	R      int
	G      int
	B      int
}

func getLegoColors(dbs *sql.DB) []LegoColor {
	rows, err := dbs.Query(`
		SELECT  id, hex, r, g, b, name
		FROM    colors
	`)
	if err != nil {
		log.Fatal(err)
	}

	var legocolors []LegoColor
	for rows.Next() {
		var id, r, g, b int
		var hex, name string
		err := rows.Scan(&id, &hex, &r, &g, &b, &name)
		if err != nil {
			panic(err)
		}

		legocolors = append(legocolors, LegoColor{
			ID:   id,
			Hex:  hex,
			R:    r,
			G:    g,
			B:    b,
			Name: name,
		})
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return legocolors
}

// Pixel struct example
type Pixel struct {
	R int
	G int
	B int
	A int
}

func resizeImage(inimage io.Reader, x, y int) io.Reader {
	// Decode the image (from PNG to image.Image):
	source, err := png.Decode(inimage)
	if err != nil {
		panic(err)
	}
	// Set the expected size that you want:
	// dst := image.NewRGBA(image.Rect(x, y, src.Bounds().Max.X/2, src.Bounds().Max.Y/2))
	m := image.NewRGBA(
		image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: x, Y: y},
		},
	)

	// Resize and encode:
	draw.NearestNeighbor.Scale(m, m.Rect, source, source.Bounds(), draw.Over, nil)

	output, err := os.Create(outimages + randStringRunes(4) + "_resize.png")
	if err != nil {
		panic(err)
	}
	png.Encode(output, m)
	output.Seek(0, 0)

	return output
}

// calculateDistance defines the distance between two pixels.
// Inspiration: https://stackoverflow.com/a/1847112
func calculateDistance(p1, p2 Pixel) (distance float64) {
	p1R, p1G, p1B := float64(p1.R), float64(p1.G), float64(p1.B)
	p2R, p2G, p2B := float64(p2.R), float64(p2.G), float64(p2.B)

	return math.Sqrt(math.Pow((p2R-p1R)*0.30, 2) + math.Pow((p2G-p1G)*0.59, 2) + math.Pow((p2B-p1B)*0.11, 2))
}

// lego converts an image inside an io.Reader into its lego version, an image in the "png" format is expected
func lego(ctx context.Context, dbs *sql.DB, inimage io.Reader) *image.RGBA {
	// Calling the generic image.Decode() will tell give us the data
	// and type of image it is as a string. We expect "png"
	imageData, _, err := image.Decode(inimage)
	if err != nil {
		panic(err)
	}

	x0len, y0len := imageData.Bounds().Min.X, imageData.Bounds().Min.Y
	xlen, ylen := imageData.Bounds().Max.X, imageData.Bounds().Max.Y
	legocolors := getLegoColors(dbs)
	legoimage := image.NewRGBA(
		image.Rectangle{
			Min: image.Point{X: x0len, Y: y0len},
			Max: image.Point{X: xlen, Y: ylen},
		},
	)

	buildingMap := make([][]string, xlen)
	uniqueColors := make(map[string]struct{}, xlen*ylen)

	var mindistance float64
	var bestmatch int
	var r, g, b uint8
	for x := 0; x < xlen; x++ {
		// mindistance = math.MaxFloat64 // Arbitrary high value to allow finding a lower number
		buildingMap[x] = make([]string, ylen)

		for y := 0; y < ylen; y++ {
			mindistance = math.MaxFloat64 // Arbitrary high value to allow finding a lower number
			r32, g32, b32, _ := imageData.At(x, y).RGBA()
			r, g, b = uint8(r32), uint8(g32), uint8(b32)

			// for each pixel, loop over all the lego colors to find the closest color
			for i := range legocolors {
				distance := calculateDistance(
					Pixel{R: legocolors[i].R, G: legocolors[i].G, B: legocolors[i].B},
					Pixel{R: int(r), G: int(g), B: int(b)},
				)

				// building a new image by replacing the real color for the most-close-lego-color
				// https://cs.opensource.google/go/go/+/refs/tags/go1.17.5:src/image/image.go;l=96
				if distance < mindistance {
					mindistance = distance
					bestmatch = i
					legoimage.SetRGBA(x, y, color.RGBA{
						R: uint8(legocolors[i].R), G: uint8(legocolors[i].G), B: uint8(legocolors[i].B), A: 255,
					})
				}
			}

			// Set RGB color on a specific pixel
			legoimage.SetRGBA(x, y, color.RGBA{
				R: uint8(legocolors[bestmatch].R), G: uint8(legocolors[bestmatch].G), B: uint8(legocolors[bestmatch].B), A: 255,
			})

			// Add lego color to the building map
			buildingMap[x][y] = fmt.Sprintf("[%d][%d] = R:%d, G:%d, B:%d\t-%s\n",
				x, y, legocolors[bestmatch].R, legocolors[bestmatch].G, legocolors[bestmatch].B, legocolors[bestmatch].Name)

			uniqueColors[legocolors[bestmatch].Name] = struct{}{}

		} // end y loop
	} // end x loop

	f, _ := os.Create(outfolder + "build_map.txt")
	for i := 0; i < len(buildingMap); i++ {
		for j := 0; j < len(buildingMap[i]); j++ {
			_, _ = f.WriteString(buildingMap[i][j])
			pieces++
		}
	}
	f.Close()

	// set num of colors
	colors = len(uniqueColors)

	return legoimage
}

func main() {
	run()
}

func run() {
	if len(os.Args) != 4 {
		log.Printf("Usage: %s file_name X Y", os.Args[0])
		os.Exit(1)
	}

	xlen, _ := strconv.Atoi(os.Args[2])
	ylen, _ := strconv.Atoi(os.Args[3])

	ctx := context.Background()
	dbs := dbConnect("lego_artboard")

	// open input image
	existingImageFile, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer existingImageFile.Close()

	// resize input image
	resizedimg := resizeImage(existingImageFile, xlen, ylen)
	// img, _ := resizedimg.(*os.File) // FIXME: will panic if file is invalid
	// img.Close()

	// parse pixels, find closest color based on lego pieces
	outimage := lego(ctx, dbs, resizedimg)

	// create and encode output image
	outname := outimages + randStringRunes(8) + "_out.png"
	outfile, err := os.Create(outname)
	if err != nil {
		panic(err)
	}
	png.Encode(outfile, outimage)

	log.Printf("input=%q, dimensions=%dx%d, output=%q\n", os.Args[1], xlen, ylen, outname)
	log.Printf("For this Lego conversion have been used %d pieces and %d colors\n", pieces, colors)
	log.Printf("The building map has been generated on %q file", "build_map.txt")
	// log.Printf("used: colors=%d, pieces=%d, build")
}
