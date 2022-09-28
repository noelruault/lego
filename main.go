package main

import (
	"context"
	"encoding/csv"
	"flag"
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

	"golang.org/x/image/draw"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	letterRunes   = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	defaultColors = []LegoColor{
		{Hex: "FFE371", LegoID: 1027, Name: "Modulex Light Yellow", R: 255, G: 227, B: 113},
		{Hex: "CA1F08", LegoID: 1010, Name: "Vintage Red", R: 202, G: 31, B: 8},
		{Hex: "F08F1C", LegoID: 182, Name: "Trans-Orange", R: 240, G: 143, B: 28},
		{Hex: "6C6E68", LegoID: 72, Name: "Dark Bluish Gray", R: 108, G: 110, B: 104},
		{Hex: "FFA70B", LegoID: 462, Name: "Medium Orange", R: 255, G: 167, B: 11},
		{Hex: "720E0F", LegoID: 320, Name: "Dark Red", R: 114, G: 14, B: 15},
		{Hex: "D67572", LegoID: 335, Name: "Sand Red", R: 214, G: 117, B: 114},
		{Hex: "6D6E5C", LegoID: 8, Name: "Dark Gray", R: 109, G: 110, B: 92},
		{Hex: "F2705E", LegoID: 12, Name: "Salmon", R: 242, G: 112, B: 94},
		{Hex: "FCFCFC", LegoID: 47, Name: "Trans-Clear", R: 252, G: 252, B: 252},
		{Hex: "F2F3F2", LegoID: 183, Name: "Pearl White", R: 242, G: 243, B: 242},
		{Hex: "FECCCF", LegoID: 77, Name: "Light Pink", R: 254, G: 204, B: 207},
		{Hex: "184632", LegoID: 288, Name: "Dark Green", R: 24, G: 70, B: 50},
		{Hex: "9FC3E9", LegoID: 212, Name: "Bright Light Blue", R: 159, G: 195, B: 233},
		{Hex: "68BCC5", LegoID: 1003, Name: "Glitter Trans-Light Blue", R: 104, G: 188, B: 197},
		{Hex: "7988A1", LegoID: 137, Name: "Metal Blue", R: 121, G: 136, B: 161},
		{Hex: "0020A0", LegoID: 33, Name: "Trans-Dark Blue", R: 0, G: 32, B: 160},
		{Hex: "C1DFF0", LegoID: 43, Name: "Trans-Very Lt Blue", R: 193, G: 223, B: 240},
		{Hex: "9B9A5A", LegoID: 326, Name: "Olive Green", R: 155, G: 154, B: 90},
		{Hex: "645A4C", LegoID: 60, Name: "Chrome Antique Brass", R: 100, G: 90, B: 76},
		{Hex: "F2CD37", LegoID: 14, Name: "Yellow", R: 242, G: 205, B: 55},
		{Hex: "F6D7B3", LegoID: 78, Name: "Light Flesh", R: 246, G: 215, B: 179},
		{Hex: "4B9F4A", LegoID: 10, Name: "Bright Green", R: 75, G: 159, B: 74},
		{Hex: "E1D5ED", LegoID: 31, Name: "Lavender", R: 225, G: 213, B: 237},
		{Hex: "84B68D", LegoID: 34, Name: "Trans-Green", R: 132, G: 182, B: 141},
		{Hex: "5A93DB", LegoID: 73, Name: "Medium Blue", R: 90, G: 147, B: 219},
		{Hex: "A5A5CB", LegoID: 129, Name: "Glitter Trans-Purple", R: 165, G: 165, B: 203},
		{Hex: "3CB371", LegoID: 62, Name: "Chrome Green", R: 60, G: 179, B: 113},
		{Hex: "05131D", LegoID: 132, Name: "Speckle Black-Silver", R: 5, G: 19, B: 29},
		{Hex: "FBE890", LegoID: 1005, Name: "Trans-Fire Yellow", R: 251, G: 232, B: 144},
		{Hex: "8E5597", LegoID: 1007, Name: "Reddish Lilac", R: 142, G: 85, B: 151},
		{Hex: "FFF03A", LegoID: 226, Name: "Bright Light Yellow", R: 255, G: 240, B: 58},
		{Hex: "05131D", LegoID: 0, Name: "Black", R: 5, G: 19, B: 29},
		{Hex: "0055BF", LegoID: 1, Name: "Blue", R: 0, G: 85, B: 191},
		{Hex: "C870A0", LegoID: 5, Name: "Dark Pink", R: 200, G: 112, B: 160},
		{Hex: "B4D2E3", LegoID: 9, Name: "Light Blue", R: 180, G: 210, B: 227},
		{Hex: "FC97AC", LegoID: 13, Name: "Pink", R: 252, G: 151, B: 172},
		{Hex: "C2DAB8", LegoID: 17, Name: "Light Green", R: 194, G: 218, B: 184},
		{Hex: "E4CD9E", LegoID: 19, Name: "Tan", R: 228, G: 205, B: 158},
		{Hex: "81007B", LegoID: 22, Name: "Purple", R: 129, G: 0, B: 123},
		{Hex: "2032B0", LegoID: 23, Name: "Dark Blue-Violet", R: 32, G: 50, B: 176},
		{Hex: "923978", LegoID: 26, Name: "Magenta", R: 146, G: 57, B: 120},
		{Hex: "635F52", LegoID: 32, Name: "Trans-Black IR Lens", R: 99, G: 95, B: 82},
		{Hex: "D9E4A7", LegoID: 35, Name: "Trans-Bright Green", R: 217, G: 228, B: 167},
		{Hex: "AEEFEC", LegoID: 41, Name: "Trans-Light Blue", R: 174, G: 239, B: 236},
		{Hex: "F8F184", LegoID: 42, Name: "Trans-Neon Green", R: 248, G: 241, B: 132},
		{Hex: "A5A5CB", LegoID: 52, Name: "Trans-Purple", R: 165, G: 165, B: 203},
		{Hex: "DAB000", LegoID: 54, Name: "Trans-Neon Yellow", R: 218, G: 176, B: 0},
		{Hex: "FF800D", LegoID: 57, Name: "Trans-Neon Orange", R: 255, G: 128, B: 13},
		{Hex: "582A12", LegoID: 70, Name: "Reddish Brown", R: 88, G: 42, B: 18},
		{Hex: "A0A5A9", LegoID: 71, Name: "Light Bluish Gray", R: 160, G: 165, B: 169},
		{Hex: "6C6E68", LegoID: 76, Name: "Speckle DBGray-Silver", R: 108, G: 110, B: 104},
		{Hex: "FFFFFF", LegoID: 79, Name: "Milky White", R: 255, G: 255, B: 255},
		{Hex: "CC702A", LegoID: 84, Name: "Medium Dark Flesh", R: 204, G: 112, B: 42},
		{Hex: "7C503A", LegoID: 86, Name: "Dark Flesh", R: 124, G: 80, B: 58},
		{Hex: "4C61DB", LegoID: 89, Name: "Royal Blue", R: 76, G: 97, B: 219},
		{Hex: "FEBABD", LegoID: 100, Name: "Light Salmon", R: 254, G: 186, B: 189},
		{Hex: "6874CA", LegoID: 112, Name: "Blue-Violet", R: 104, G: 116, B: 202},
		{Hex: "B3D7D1", LegoID: 118, Name: "Aqua", R: 179, G: 215, B: 209},
		{Hex: "D9E4A7", LegoID: 120, Name: "Light Lime", R: 217, G: 228, B: 167},
		{Hex: "AE7A59", LegoID: 134, Name: "Copper", R: 174, G: 122, B: 89},
		{Hex: "DCBC81", LegoID: 142, Name: "Pearl Light Gold", R: 220, G: 188, B: 129},
		{Hex: "DFEEA5", LegoID: 158, Name: "Yellowish Green", R: 223, G: 238, B: 165},
		{Hex: "898788", LegoID: 179, Name: "Flat Silver", R: 137, G: 135, B: 136},
		{Hex: "E4ADC8", LegoID: 230, Name: "Trans-Pink", R: 228, G: 173, B: 200},
		{Hex: "96709F", LegoID: 236, Name: "Trans-Light Purple", R: 150, G: 112, B: 159},
		{Hex: "BDC6AD", LegoID: 294, Name: "Glow In Dark Trans", R: 189, G: 198, B: 173},
		{Hex: "AA7F2E", LegoID: 297, Name: "Pearl Gold", R: 170, G: 127, B: 46},
		{Hex: "3592C3", LegoID: 313, Name: "Maersk Blue", R: 53, G: 146, B: 195},
		{Hex: "36AEBF", LegoID: 322, Name: "Medium Azure", R: 54, G: 174, B: 191},
		{Hex: "ADC3C0", LegoID: 323, Name: "Light Aqua", R: 173, G: 195, B: 192},
		{Hex: "A0BCAC", LegoID: 378, Name: "Sand Green", R: 160, G: 188, B: 172},
		{Hex: "E0E0E0", LegoID: 383, Name: "Chrome Silver", R: 224, G: 224, B: 224},
		{Hex: "B67B50", LegoID: 450, Name: "Fabuland Brown", R: 182, G: 123, B: 80},
		{Hex: "C0F500", LegoID: 1002, Name: "Glitter Trans-Neon Green", R: 192, G: 245, B: 0},
		{Hex: "039CBD", LegoID: 1008, Name: "Vintage Blue", R: 3, G: 156, B: 189},
		{Hex: "845E84", LegoID: 373, Name: "Sand Purple", R: 132, G: 94, B: 132},
		{Hex: "C91A09", LegoID: 4, Name: "Red", R: 201, G: 26, B: 9},
		{Hex: "FFFFFF", LegoID: 117, Name: "Glitter Trans-Clear", R: 255, G: 255, B: 255},
		{Hex: "595D60", LegoID: 1016, Name: "Modulex Charcoal Gray", R: 89, G: 93, B: 96},
		{Hex: "EF9121", LegoID: 1012, Name: "Fabuland Orange", R: 239, G: 145, B: 33},
		{Hex: "7DB538", LegoID: 1043, Name: "Modulex Foil Light Green", R: 125, G: 181, B: 56},
		{Hex: "899B5F", LegoID: 81, Name: "Metallic Green", R: 137, G: 155, B: 95},
		{Hex: "68BCC5", LegoID: 1053, Name: "Trans-Blue Opal", R: 104, G: 188, B: 197},
		{Hex: "F4F4F4", LegoID: 1013, Name: "Modulex White", R: 244, G: 244, B: 244},
		{Hex: "AfB5C7", LegoID: 1014, Name: "Modulex Light Bluish Gray", R: 175, G: 181, B: 199},
		{Hex: "FE8A18", LegoID: 25, Name: "Orange", R: 254, G: 138, B: 24},
		{Hex: "5C5030", LegoID: 1020, Name: "Modulex Terracotta", R: 92, G: 80, B: 48},
		{Hex: "AC78BA", LegoID: 30, Name: "Medium Lavender", R: 172, G: 120, B: 186},
		{Hex: "DEC69C", LegoID: 1022, Name: "Modulex Buff", R: 222, G: 198, B: 156},
		{Hex: "DBAC34", LegoID: 82, Name: "Metallic Gold", R: 219, G: 172, B: 52},
		{Hex: "5AC4DA", LegoID: 1051, Name: "Pastel Blue", R: 90, G: 196, B: 218},
		{Hex: "FFFFFF", LegoID: 1039, Name: "Modulex Clear", R: 255, G: 255, B: 255},
		{Hex: "05131D", LegoID: 75, Name: "Speckle Black-Copper", R: 5, G: 19, B: 29},
		{Hex: "467083", LegoID: 1033, Name: "Modulex Teal Blue", R: 70, G: 112, B: 131},
		{Hex: "1B2A34", LegoID: 64, Name: "Chrome Black", R: 27, G: 42, B: 52},
		{Hex: "BBA53D", LegoID: 334, Name: "Chrome Gold", R: 187, G: 165, B: 61},
		{Hex: "4D4C52", LegoID: 1018, Name: "Modulex Black", R: 77, G: 76, B: 82},
		{Hex: "7DB538", LegoID: 1030, Name: "Modulex Pastel Green", R: 125, G: 181, B: 56},
		{Hex: "FED557", LegoID: 1028, Name: "Modulex Ochre Yellow", R: 254, G: 213, B: 87},
		{Hex: "595D60", LegoID: 1040, Name: "Modulex Foil Dark Gray", R: 89, G: 93, B: 96},
		{Hex: "6B5A5A", LegoID: 1017, Name: "Modulex Tile Gray", R: 107, G: 90, B: 90},
		{Hex: "BDC618", LegoID: 1029, Name: "Modulex Lemon", R: 189, G: 198, B: 24},
		{Hex: "68AECE", LegoID: 1045, Name: "Modulex Foil Light Blue", R: 104, G: 174, B: 206},
		{Hex: "C9E788", LegoID: 1057, Name: "Trans-Light Bright Green", R: 201, G: 231, B: 136},
		{Hex: "9C9C9C", LegoID: 1015, Name: "Modulex Light Gray", R: 156, G: 156, B: 156},
		{Hex: "DF6695", LegoID: 114, Name: "Glitter Trans-Dark Pink", R: 223, G: 102, B: 149},
		{Hex: "6074A1", LegoID: 379, Name: "Sand Blue", R: 96, G: 116, B: 161},
		{Hex: "FBE696", LegoID: 18, Name: "Light Yellow", R: 251, G: 230, B: 150},
		{Hex: "6C96BF", LegoID: 61, Name: "Chrome Blue", R: 108, G: 150, B: 191},
		{Hex: "FFFFFF", LegoID: 15, Name: "White", R: 255, G: 255, B: 255},
		{Hex: "958A73", LegoID: 28, Name: "Dark Tan", R: 149, G: 138, B: 115},
		{Hex: "635F52", LegoID: 40, Name: "Trans-Black", R: 99, G: 95, B: 82},
		{Hex: "F5CD2F", LegoID: 46, Name: "Trans-Yellow", R: 245, G: 205, B: 47},
		{Hex: "AA4D8E", LegoID: 63, Name: "Chrome Pink", R: 170, G: 77, B: 142},
		{Hex: "73DCA1", LegoID: 74, Name: "Medium Green", R: 115, G: 220, B: 161},
		{Hex: "3F3691", LegoID: 85, Name: "Dark Purple", R: 63, G: 54, B: 145},
		{Hex: "F9BA61", LegoID: 125, Name: "Light Orange", R: 249, G: 186, B: 97},
		{Hex: "575857", LegoID: 148, Name: "Pearl Dark Gray", R: 87, G: 88, B: 87},
		{Hex: "7DBFDD", LegoID: 232, Name: "Sky Blue", R: 125, G: 191, B: 221},
		{Hex: "352100", LegoID: 308, Name: "Dark Brown", R: 53, G: 33, B: 0},
		{Hex: "F785B1", LegoID: 351, Name: "Medium Dark Pink", R: 247, G: 133, B: 177},
		{Hex: "A95500", LegoID: 484, Name: "Dark Orange", R: 169, G: 85, B: 0},
		{Hex: "1E601E", LegoID: 1009, Name: "Vintage Green", R: 30, G: 96, B: 30},
		{Hex: "F3C305", LegoID: 1011, Name: "Vintage Yellow", R: 243, G: 195, B: 5},
		{Hex: "F785B1", LegoID: 1038, Name: "Modulex Pink", R: 247, G: 133, B: 177},
		{Hex: "0057A6", LegoID: 1044, Name: "Modulex Foil Dark Blue", R: 0, G: 87, B: 166},
		{Hex: "FED557", LegoID: 1048, Name: "Modulex Foil Yellow", R: 254, G: 213, B: 87},
		{Hex: "F08F1C", LegoID: 1052, Name: "Glitter Trans-Orange", R: 240, G: 143, B: 28},
		{Hex: "FCFCFC", LegoID: 1055, Name: "Trans-Clear Opal", R: 252, G: 252, B: 252},
		{Hex: "583927", LegoID: 1056, Name: "Trans-Brown Opal", R: 88, G: 57, B: 39},
		{Hex: "8320B7", LegoID: 1059, Name: "Trans-Purple Opal", R: 131, G: 32, B: 183},
		{Hex: "0020A0", LegoID: 1061, Name: "Trans-Dark Blue Opal", R: 0, G: 32, B: 160},
		{Hex: "EBD800", LegoID: 1062, Name: "Lemon", R: 235, G: 216, B: 0},
		{Hex: "9391E4", LegoID: 1001, Name: "Medium Violet", R: 147, G: 145, B: 228},
		{Hex: "F47B30", LegoID: 1025, Name: "Modulex Orange", R: 244, G: 123, B: 48},
		{Hex: "B4D4F7", LegoID: 1006, Name: "Trans-Light Royal Blue", R: 180, G: 212, B: 247},
		{Hex: "7C9051", LegoID: 1031, Name: "Modulex Olive Green", R: 124, G: 144, B: 81},
		{Hex: "B48455", LegoID: 178, Name: "Flat Dark Gold", R: 180, G: 132, B: 85},
		{Hex: "F7AD63", LegoID: 1049, Name: "Modulex Foil Orange", R: 247, G: 173, B: 99},
		{Hex: "94E5AB", LegoID: 1058, Name: "Trans-Light Green", R: 148, G: 229, B: 171},
		{Hex: "F45C40", LegoID: 1024, Name: "Modulex Pink Red", R: 244, G: 92, B: 64},
		{Hex: "27867E", LegoID: 1032, Name: "Modulex Aqua Green", R: 39, G: 134, B: 126},
		{Hex: "FF698F", LegoID: 1050, Name: "Coral", R: 255, G: 105, B: 143},
		{Hex: "84B68D", LegoID: 1060, Name: "Trans-Green Opal", R: 132, G: 182, B: 141},
		{Hex: "330000", LegoID: 1019, Name: "Modulex Tile Brown", R: 51, G: 0, B: 0},
		{Hex: "BD7D85", LegoID: 1037, Name: "Modulex Violet", R: 189, G: 125, B: 133},
		{Hex: "B52C20", LegoID: 1023, Name: "Modulex Red", R: 181, G: 44, B: 32},
		{Hex: "9C9C9C", LegoID: 1041, Name: "Modulex Foil Light Gray", R: 156, G: 156, B: 156},
		{Hex: "4B0082", LegoID: 1046, Name: "Modulex Foil Violet", R: 75, G: 0, B: 130},
		{Hex: "61AFFF", LegoID: 1035, Name: "Modulex Medium Blue", R: 97, G: 175, B: 255},
		{Hex: "0057A6", LegoID: 1034, Name: "Modulex Tile Blue", R: 0, G: 87, B: 166},
		{Hex: "6400", LegoID: 1042, Name: "Modulex Foil Dark Green", R: 0, G: 100, B: 0},
		{Hex: "ABADAC", LegoID: 150, Name: "Pearl Very Light Gray", R: 171, G: 173, B: 172},
		{Hex: "CD6298", LegoID: 69, Name: "Light Purple", R: 205, G: 98, B: 152},
		{Hex: "008F9B", LegoID: 3, Name: "Dark Turquoise", R: 0, G: 143, B: 155},
		{Hex: "E4ADC8", LegoID: 29, Name: "Bright Pink", R: 228, G: 173, B: 200},
		{Hex: "9BA19D", LegoID: 7, Name: "Light Gray", R: 155, G: 161, B: 157},
		{Hex: "E6E3E0", LegoID: 151, Name: "Very Light Bluish Gray", R: 230, G: 227, B: 224},
		{Hex: "D9D9D9", LegoID: 1000, Name: "Glow in Dark White", R: 217, G: 217, B: 217},
		{Hex: "237841", LegoID: 2, Name: "Green", R: 35, G: 120, B: 65},
		{Hex: "55A5AF", LegoID: 11, Name: "Light Turquoise", R: 85, G: 165, B: 175},
		{Hex: "C9CAE2", LegoID: 20, Name: "Light Violet", R: 201, G: 202, B: 226},
		{Hex: "BBE90B", LegoID: 27, Name: "Lime", R: 187, G: 233, B: 11},
		{Hex: "0A3463", LegoID: 272, Name: "Dark Blue", R: 10, G: 52, B: 99},
		{Hex: "68AECE", LegoID: 1036, Name: "Modulex Pastel Blue", R: 104, G: 174, B: 206},
		{Hex: "05131D", LegoID: 133, Name: "Speckle Black-Gold", R: 5, G: 19, B: 29},
		{Hex: "C91A09", LegoID: 36, Name: "Trans-Red", R: 201, G: 26, B: 9},
		{Hex: "FCB76D", LegoID: 1004, Name: "Trans-Flame Yellowish Orange", R: 252, G: 183, B: 109},
		{Hex: "CFE2F7", LegoID: 143, Name: "Trans-Medium Blue", R: 207, G: 226, B: 247},
		{Hex: "A5A9B4", LegoID: 80, Name: "Metallic Silver", R: 165, G: 169, B: 180},
		{Hex: "9CA3A8", LegoID: 135, Name: "Pearl Light Gray", R: 156, G: 163, B: 168},
		{Hex: "F8BB3D", LegoID: 191, Name: "Bright Light Orange", R: 248, G: 187, B: 61},
		{Hex: "E6E3DA", LegoID: 503, Name: "Very Light Gray", R: 230, G: 227, B: 218},
		{Hex: "078BC9", LegoID: 321, Name: "Dark Azure", R: 7, G: 139, B: 201},
		{Hex: "DF6695", LegoID: 45, Name: "Trans-Dark Pink", R: 223, G: 102, B: 149},
		{Hex: "907450", LegoID: 1021, Name: "Modulex Brown", R: 144, G: 116, B: 80},
		{Hex: "4354A3", LegoID: 110, Name: "Violet", R: 67, G: 84, B: 163},
		{Hex: "D09168", LegoID: 92, Name: "Flesh", R: 208, G: 145, B: 104},
		{Hex: "8B0000", LegoID: 1047, Name: "Modulex Foil Red", R: 139, G: 0, B: 0},
		{Hex: "CE1D9B", LegoID: 1054, Name: "Trans-Medium Reddish Violet Opal", R: 206, G: 29, B: 155},
		{Hex: "B31004", LegoID: 216, Name: "Rust", R: 179, G: 16, B: 4},
		{Hex: "583927", LegoID: 6, Name: "Brown", R: 88, G: 57, B: 39},
		{Hex: "F7AD63", LegoID: 1026, Name: "Modulex Light Orange", R: 247, G: 173, B: 99},
		{Hex: "F3CF9B", LegoID: 68, Name: "Very Light Orange", R: 243, G: 207, B: 155},
		{Hex: "C7D23C", LegoID: 115, Name: "Medium Lime", R: 199, G: 210, B: 60},
		{Hex: "FA9C1C", LegoID: 366, Name: "Earth Orange", R: 250, G: 156, B: 28},
		{Hex: "D4D5C9", LegoID: 21, Name: "Glow In Dark Opaque", R: 212, G: 213, B: 201},
	}
)

type LegoColor struct {
	Hex    string
	LegoID int
	Name   string
	R      int
	G      int
	B      int
}

// https://rebrickable.com/downloads/
func colorsFromCSV(f io.Reader) ([]LegoColor, error) {
	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("unable to parse file as CSV: err=%v", err)
	}

	// the funny -1, 1: is to take care of the header
	var colors = make([]LegoColor, len(records)-1, len(records)-1)
	for i, record := range records[1:] {
		legoid, _ := strconv.Atoi(record[0])
		r, _ := strconv.Atoi(record[3])
		g, _ := strconv.Atoi(record[4])
		b, _ := strconv.Atoi(record[5])
		colors[i] = LegoColor{
			LegoID: legoid,
			Name:   record[1],
			Hex:    record[2],
			R:      r,
			G:      g,
			B:      b,
		}
	}

	return colors, nil
}

// Pixel struct example
type Pixel struct {
	R int
	G int
	B int
	A int
}

// calculateDistance defines the distance between two pixels.
// Inspiration: https://stackoverflow.com/a/1847112
func calculateDistance(p1, p2 Pixel) (distance float64) {
	p1R, p1G, p1B := float64(p1.R), float64(p1.G), float64(p1.B)
	p2R, p2G, p2B := float64(p2.R), float64(p2.G), float64(p2.B)

	return math.Sqrt(math.Pow((p2R-p1R)*0.30, 2) + math.Pow((p2G-p1G)*0.59, 2) + math.Pow((p2B-p1B)*0.11, 2))
}

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func resizePNGImage(inimage io.Reader, x, y int) (*image.RGBA, error) {
	// Decode the image (from PNG to image.Image):
	source, err := png.Decode(inimage)
	if err != nil {
		return nil, fmt.Errorf("decoding png: err=%v", err)
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

	return m, nil
}

type Lego struct {
	colors []LegoColor
}

// mapFromImage converts an image inside an io.Reader into its lego version, an image in the "png" format is expected
func (l *Lego) mapFromImage(ctx context.Context, imageData *image.RGBA) (*Conversion, error) {
	if len(l.colors) < 1 {
		return nil, fmt.Errorf("there aren't disponible colors to work with")
	}

	x0len, y0len := imageData.Bounds().Min.X, imageData.Bounds().Min.Y
	xlen, ylen := imageData.Bounds().Max.X, imageData.Bounds().Max.Y
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
			for i := range l.colors {
				distance := calculateDistance(
					Pixel{R: l.colors[i].R, G: l.colors[i].G, B: l.colors[i].B},
					Pixel{R: int(r), G: int(g), B: int(b)},
				)

				// building a new image by replacing the real color for the most-close-lego-color
				// https://cs.opensource.google/go/go/+/refs/tags/go1.17.5:src/image/image.go;l=96
				if distance < mindistance {
					mindistance = distance
					bestmatch = i
					legoimage.SetRGBA(x, y, color.RGBA{
						R: uint8(l.colors[i].R), G: uint8(l.colors[i].G), B: uint8(l.colors[i].B), A: 255,
					})
				}
			}

			// Set RGB color on a specific pixel
			legoimage.SetRGBA(x, y, color.RGBA{
				R: uint8(l.colors[bestmatch].R), G: uint8(l.colors[bestmatch].G), B: uint8(l.colors[bestmatch].B), A: 255,
			})

			// Add lego color to the building map
			buildingMap[x][y] = fmt.Sprintf("[%d][%d] = R:%d, G:%d, B:%d\t-%s\n",
				x, y, l.colors[bestmatch].R, l.colors[bestmatch].G, l.colors[bestmatch].B, l.colors[bestmatch].Name)

			uniqueColors[l.colors[bestmatch].Name] = struct{}{}

		} // end y loop
	} // end x loop

	return &Conversion{
		Image:      legoimage,
		ColorsUsed: len(uniqueColors),
		BuildMap:   buildingMap,
	}, nil
}

type Flags struct {
	ColorsCSVPath string
	ImagePath     string
	OutPath       string
	XLen          int
	YLen          int
}

func parseFlags() *Flags {
	colorsCSVPath := flag.String("colors", "", "(Required) CSV file that contains a list of colors, with the format: legoid,name,hex,r,g,b ")
	imagePath := flag.String("image", "", "(Required) Target image path")
	outPath := flag.String("out", "", "")
	xlen := flag.Int("xlen", 100, "")
	ylen := flag.Int("ylen", 100, "")

	flag.Parse()

	return &Flags{
		ColorsCSVPath: *colorsCSVPath,
		ImagePath:     *imagePath,
		OutPath:       *outPath,
		XLen:          *xlen,
		YLen:          *ylen,
	}
}

type Conversion struct {
	Image      *image.RGBA
	PiecesUsed int
	ColorsUsed int
	BuildMap   [][]string
}

func (c *Conversion) result(outPath string) error {
	randResultName := randStringRunes(8)

	var piecesUsed int
	buildMapFileName := outPath + randResultName + "_build_map.txt"
	f, _ := os.Create(buildMapFileName)
	for i := 0; i < len(c.BuildMap); i++ {
		for j := 0; j < len(c.BuildMap[i]); j++ {
			_, _ = f.WriteString(c.BuildMap[i][j])
			piecesUsed++
		}
	}
	f.Close()
	c.PiecesUsed = piecesUsed

	// create and encode output image
	ImageResultFileName := outPath + randResultName + "_out.png"
	outfile, err := os.Create(ImageResultFileName)
	if err != nil {
		return fmt.Errorf("creating output file: err=%v", err)
	}
	png.Encode(outfile, c.Image)

	log.Printf("For this Lego conversion have been used %d pieces and %d colors\n", c.PiecesUsed, c.ColorsUsed)
	log.Printf("The image preview has been generated at %q ", ImageResultFileName)
	log.Printf("The building map has been generated at %q", buildMapFileName)

	return nil
}

func run() {
	ctx := context.Background()
	flags := parseFlags()

	if flags.ImagePath == "" {
		flag.Usage()
		os.Exit(1)
	}

	if flags.OutPath != "" {
		if _, err := os.Stat(flags.OutPath); os.IsNotExist(err) {
			log.Printf("directory path %q does not exist", flags.OutPath)
			os.Exit(1)
		}
	}

	csvColors := defaultColors
	if flags.ColorsCSVPath != "" {
		colorsFile, err := os.Open(flags.ColorsCSVPath)
		if err != nil {
			log.Printf("opening colors CSV: file=%s, err=%v", flags.ColorsCSVPath, err)
			os.Exit(1)
		}
		defer colorsFile.Close()

		csvColors, err = colorsFromCSV(colorsFile)
		if err != nil {
			log.Printf("retrieving colors: file=%s, err=%v", flags.ColorsCSVPath, err)
			os.Exit(1)
		}
	}

	inputImage, err := os.Open(flags.ImagePath)
	if err != nil {
		log.Printf("opening PNG image: file=%s, err=%v", flags.ImagePath, err)
		os.Exit(1)
	}
	defer inputImage.Close()

	resizedimg, err := resizePNGImage(inputImage, flags.XLen, flags.YLen)
	if err != nil {
		log.Printf("resizing image: err=%v", err)
		os.Exit(1)
	}

	// parse pixels, find closest color based on the available lego pieces
	lego := Lego{colors: csvColors}
	conversion, err := lego.mapFromImage(ctx, resizedimg)
	if err != nil {
		log.Printf("mapping image to lego artboard: err=%v", err)
		os.Exit(1)
	}

	log.Printf("INFO: input=%q, dimensions=%dx%d", flags.ImagePath, flags.XLen, flags.YLen)
	err = conversion.result(flags.OutPath)
	if err != nil {
		log.Printf("processing result: err=%v", err)
		os.Exit(1)
	}

}

func main() {
	run()
}
