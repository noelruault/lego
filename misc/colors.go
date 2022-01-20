package colors

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
)

type Hex string

type RGB struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

type Color struct {
	Hex    string
	ID     int
	LegoID int
	Name   string
	R      int
	G      int
	B      int
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

func colors() {
	var rgb RGB
	var err error

	dbs := dbConnect("lego_artboard")
	rows, err := dbs.Query(`
		SELECT
			id,
			hex
		FROM
			colors
	`)
	if err != nil {
		log.Fatal(err)
	}

	var colors []Color
	for rows.Next() {
		var id int
		var hex string
		err := rows.Scan(&id, &hex)
		if err != nil {
			panic(err)
		}

		colors = append(colors, Color{
			ID:  id,
			Hex: hex,
		})
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	for i := range colors {
		rgb, err = Hex2RGB(Hex(colors[i].Hex))
		if err != nil {
			log.Printf("Couldn't convert hex to rgb: %s", colors[i].Hex)
		}

		colors[i].R = int(rgb.Red)
		colors[i].G = int(rgb.Green)
		colors[i].B = int(rgb.Blue)
	}

	for i := range colors {
		sqlStatement := fmt.Sprintf(`
		UPDATE colors SET
			R = %d,
			G = %d,
			B = %d
		WHERE id = %d
		`, colors[i].R, colors[i].G, colors[i].B, colors[i].ID)

		_, err := dbs.Exec(sqlStatement)
		if err != nil {
			panic(err)
		}
	}

}

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
