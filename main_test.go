package main

import (
	"math"
	"testing"

	_ "github.com/lib/pq"
)

var legocolors = []Pixel{
	{R: 179, G: 215, B: 209}, // Aqua
	{R: 5, G: 19, B: 29},     // Black
	{R: 0, G: 85, B: 191},    // Blue
	{R: 104, G: 116, B: 202}, // Blue-Violet
	{R: 75, G: 159, B: 74},   // Bright Green
	{R: 159, G: 195, B: 233}, // Bright Light Blue
	{R: 248, G: 187, B: 61},  // Bright Light Orange
	{R: 255, G: 240, B: 58},  // Bright Light Yellow
	{R: 228, G: 173, B: 200}, // Bright Pink
	{R: 88, G: 57, B: 39},    // Brown
	{R: 100, G: 90, B: 76},   // Chrome Antique Brass
	{R: 27, G: 42, B: 52},    // Chrome Black
	{R: 108, G: 150, B: 191}, // Chrome Blue
	{R: 187, G: 165, B: 61},  // Chrome Gold
	{R: 60, G: 179, B: 113},  // Chrome Green
	{R: 170, G: 77, B: 142},  // Chrome Pink
	{R: 224, G: 224, B: 224}, // Chrome Silver
	{R: 174, G: 122, B: 89},  // Copper
	{R: 255, G: 105, B: 143}, // Coral
	{R: 7, G: 139, B: 201},   // Dark Azure
	{R: 10, G: 52, B: 99},    // Dark Blue
	{R: 32, G: 50, B: 176},   // Dark Blue-Violet
	{R: 108, G: 110, B: 104}, // Dark Bluish Gray
	{R: 53, G: 33, B: 0},     // Dark Brown
	{R: 124, G: 80, B: 58},   // Dark Flesh
	{R: 109, G: 110, B: 92},  // Dark Gray
	{R: 24, G: 70, B: 50},    // Dark Green
	{R: 169, G: 85, B: 0},    // Dark Orange
	{R: 200, G: 112, B: 160}, // Dark Pink
	{R: 63, G: 54, B: 145},   // Dark Purple
	{R: 114, G: 14, B: 15},   // Dark Red
	{R: 149, G: 138, B: 115}, // Dark Tan
	{R: 0, G: 143, B: 155},   // Dark Turquoise
	{R: 250, G: 156, B: 28},  // Earth Orange
	{R: 182, G: 123, B: 80},  // Fabuland Brown
	{R: 239, G: 145, B: 33},  // Fabuland Orange
	{R: 180, G: 132, B: 85},  // Flat Dark Gold
	{R: 137, G: 135, B: 136}, // Flat Silver
	{R: 208, G: 145, B: 104}, // Flesh
	{R: 255, G: 255, B: 255}, // Glitter Trans-Clear
	{R: 223, G: 102, B: 149}, // Glitter Trans-Dark Pink
	{R: 104, G: 188, B: 197}, // Glitter Trans-Light Blue
	{R: 192, G: 245, B: 0},   // Glitter Trans-Neon Green
	{R: 240, G: 143, B: 28},  // Glitter Trans-Orange
	{R: 165, G: 165, B: 203}, // Glitter Trans-Purple
	{R: 212, G: 213, B: 201}, // Glow In Dark Opaque
	{R: 189, G: 198, B: 173}, // Glow In Dark Trans
	{R: 217, G: 217, B: 217}, // Glow in Dark White
	{R: 35, G: 120, B: 65},   // Green
	{R: 225, G: 213, B: 237}, // Lavender
	{R: 235, G: 216, B: 0},   // Lemon
	{R: 173, G: 195, B: 192}, // Light Aqua
	{R: 180, G: 210, B: 227}, // Light Blue
	{R: 160, G: 165, B: 169}, // Light Bluish Gray
	{R: 246, G: 215, B: 179}, // Light Flesh
	{R: 155, G: 161, B: 157}, // Light Gray
	{R: 194, G: 218, B: 184}, // Light Green
	{R: 217, G: 228, B: 167}, // Light Lime
	{R: 249, G: 186, B: 97},  // Light Orange
	{R: 254, G: 204, B: 207}, // Light Pink
	{R: 205, G: 98, B: 152},  // Light Purple
	{R: 254, G: 186, B: 189}, // Light Salmon
	{R: 85, G: 165, B: 175},  // Light Turquoise
	{R: 201, G: 202, B: 226}, // Light Violet
	{R: 251, G: 230, B: 150}, // Light Yellow
	{R: 187, G: 233, B: 11},  // Lime
	{R: 53, G: 146, B: 195},  // Maersk Blue
	{R: 146, G: 57, B: 120},  // Magenta
	{R: 54, G: 174, B: 191},  // Medium Azure
	{R: 90, G: 147, B: 219},  // Medium Blue
	{R: 204, G: 112, B: 42},  // Medium Dark Flesh
	{R: 247, G: 133, B: 177}, // Medium Dark Pink
	{R: 115, G: 220, B: 161}, // Medium Green
	{R: 172, G: 120, B: 186}, // Medium Lavender
	{R: 199, G: 210, B: 60},  // Medium Lime
	{R: 255, G: 167, B: 11},  // Medium Orange
	{R: 147, G: 145, B: 228}, // Medium Violet
	{R: 121, G: 136, B: 161}, // Metal Blue
	{R: 219, G: 172, B: 52},  // Metallic Gold
	{R: 137, G: 155, B: 95},  // Metallic Green
	{R: 165, G: 169, B: 180}, // Metallic Silver
	{R: 255, G: 255, B: 255}, // Milky White
	{R: 39, G: 134, B: 126},  // Modulex Aqua Green
	{R: 77, G: 76, B: 82},    // Modulex Black
	{R: 144, G: 116, B: 80},  // Modulex Brown
	{R: 222, G: 198, B: 156}, // Modulex Buff
	{R: 89, G: 93, B: 96},    // Modulex Charcoal Gray
	{R: 255, G: 255, B: 255}, // Modulex Clear
	{R: 0, G: 87, B: 166},    // Modulex Foil Dark Blue
	{R: 89, G: 93, B: 96},    // Modulex Foil Dark Gray
	{R: 0, G: 100, B: 0},     // Modulex Foil Dark Green
	{R: 104, G: 174, B: 206}, // Modulex Foil Light Blue
	{R: 156, G: 156, B: 156}, // Modulex Foil Light Gray
	{R: 125, G: 181, B: 56},  // Modulex Foil Light Green
	{R: 247, G: 173, B: 99},  // Modulex Foil Orange
	{R: 139, G: 0, B: 0},     // Modulex Foil Red
	{R: 75, G: 0, B: 130},    // Modulex Foil Violet
	{R: 254, G: 213, B: 87},  // Modulex Foil Yellow
	{R: 189, G: 198, B: 24},  // Modulex Lemon
	{R: 175, G: 181, B: 199}, // Modulex Light Bluish Gray
	{R: 156, G: 156, B: 156}, // Modulex Light Gray
	{R: 247, G: 173, B: 99},  // Modulex Light Orange
	{R: 255, G: 227, B: 113}, // Modulex Light Yellow
	{R: 97, G: 175, B: 255},  // Modulex Medium Blue
	{R: 254, G: 213, B: 87},  // Modulex Ochre Yellow
	{R: 124, G: 144, B: 81},  // Modulex Olive Green
	{R: 244, G: 123, B: 48},  // Modulex Orange
	{R: 104, G: 174, B: 206}, // Modulex Pastel Blue
	{R: 125, G: 181, B: 56},  // Modulex Pastel Green
	{R: 247, G: 133, B: 177}, // Modulex Pink
	{R: 244, G: 92, B: 64},   // Modulex Pink Red
	{R: 181, G: 44, B: 32},   // Modulex Red
	{R: 70, G: 112, B: 131},  // Modulex Teal Blue
	{R: 92, G: 80, B: 48},    // Modulex Terracotta
	{R: 0, G: 87, B: 166},    // Modulex Tile Blue
	{R: 51, G: 0, B: 0},      // Modulex Tile Brown
	{R: 107, G: 90, B: 90},   // Modulex Tile Gray
	{R: 189, G: 125, B: 133}, // Modulex Violet
	{R: 244, G: 244, B: 244}, // Modulex White
	{R: 155, G: 154, B: 90},  // Olive Green
	{R: 254, G: 138, B: 24},  // Orange
	{R: 90, G: 196, B: 218},  // Pastel Blue
	{R: 87, G: 88, B: 87},    // Pearl Dark Gray
	{R: 170, G: 127, B: 46},  // Pearl Gold
	{R: 220, G: 188, B: 129}, // Pearl Light Gold
	{R: 156, G: 163, B: 168}, // Pearl Light Gray
	{R: 171, G: 173, B: 172}, // Pearl Very Light Gray
	{R: 242, G: 243, B: 242}, // Pearl White
	{R: 252, G: 151, B: 172}, // Pink
	{R: 129, G: 0, B: 123},   // Purple
	{R: 201, G: 26, B: 9},    // Red
	{R: 88, G: 42, B: 18},    // Reddish Brown
	{R: 142, G: 85, B: 151},  // Reddish Lilac
	{R: 76, G: 97, B: 219},   // Royal Blue
	{R: 179, G: 16, B: 4},    // Rust
	{R: 242, G: 112, B: 94},  // Salmon
	{R: 96, G: 116, B: 161},  // Sand Blue
	{R: 160, G: 188, B: 172}, // Sand Green
	{R: 132, G: 94, B: 132},  // Sand Purple
	{R: 214, G: 117, B: 114}, // Sand Red
	{R: 125, G: 191, B: 221}, // Sky Blue
	{R: 5, G: 19, B: 29},     // Speckle Black-Copper
	{R: 5, G: 19, B: 29},     // Speckle Black-Gold
	{R: 5, G: 19, B: 29},     // Speckle Black-Silver
	{R: 108, G: 110, B: 104}, // Speckle DBGray-Silver
	{R: 228, G: 205, B: 158}, // Tan
	{R: 99, G: 95, B: 82},    // Trans-Black
	{R: 99, G: 95, B: 82},    // Trans-Black IR Lens
	{R: 104, G: 188, B: 197}, // Trans-Blue Opal
	{R: 217, G: 228, B: 167}, // Trans-Bright Green
	{R: 88, G: 57, B: 39},    // Trans-Brown Opal
	{R: 252, G: 252, B: 252}, // Trans-Clear
	{R: 252, G: 252, B: 252}, // Trans-Clear Opal
	{R: 0, G: 32, B: 160},    // Trans-Dark Blue
	{R: 0, G: 32, B: 160},    // Trans-Dark Blue Opal
	{R: 223, G: 102, B: 149}, // Trans-Dark Pink
	{R: 251, G: 232, B: 144}, // Trans-Fire Yellow
	{R: 252, G: 183, B: 109}, // Trans-Flame Yellowish Orange
	{R: 132, G: 182, B: 141}, // Trans-Green
	{R: 132, G: 182, B: 141}, // Trans-Green Opal
	{R: 174, G: 239, B: 236}, // Trans-Light Blue
	{R: 201, G: 231, B: 136}, // Trans-Light Bright Green
	{R: 148, G: 229, B: 171}, // Trans-Light Green
	{R: 150, G: 112, B: 159}, // Trans-Light Purple
	{R: 180, G: 212, B: 247}, // Trans-Light Royal Blue
	{R: 207, G: 226, B: 247}, // Trans-Medium Blue
	{R: 206, G: 29, B: 155},  // Trans-Medium Reddish Violet Opal
	{R: 248, G: 241, B: 132}, // Trans-Neon Green
	{R: 255, G: 128, B: 13},  // Trans-Neon Orange
	{R: 218, G: 176, B: 0},   // Trans-Neon Yellow
	{R: 240, G: 143, B: 28},  // Trans-Orange
	{R: 228, G: 173, B: 200}, // Trans-Pink
	{R: 165, G: 165, B: 203}, // Trans-Purple
	{R: 131, G: 32, B: 183},  // Trans-Purple Opal
	{R: 201, G: 26, B: 9},    // Trans-Red
	{R: 193, G: 223, B: 240}, // Trans-Very Lt Blue
	{R: 245, G: 205, B: 47},  // Trans-Yellow
	{R: 230, G: 227, B: 224}, // Very Light Bluish Gray
	{R: 230, G: 227, B: 218}, // Very Light Gray
	{R: 243, G: 207, B: 155}, // Very Light Orange
	{R: 3, G: 156, B: 189},   // Vintage Blue
	{R: 30, G: 96, B: 30},    // Vintage Green
	{R: 202, G: 31, B: 8},    // Vintage Red
	{R: 243, G: 195, B: 5},   // Vintage Yellow
	{R: 67, G: 84, B: 163},   // Violet
	{R: 255, G: 255, B: 255}, // White
	{R: 242, G: 205, B: 55},  // Yellow
	{R: 223, G: 238, B: 165}, // Yellowish Green
}

func Test_calculateDistance(t *testing.T) {
	tests := []struct {
		name      string
		p         Pixel
		dataset   []Pixel
		wantPixel Pixel
	}{
		{
			name:      "color red",
			p:         Pixel{R: 255},
			dataset:   []Pixel{{R: 255}, {G: 255}, {B: 255}},
			wantPixel: Pixel{R: 255},
		},
		{
			name:      "color green",
			p:         Pixel{G: 255},
			dataset:   []Pixel{{R: 255}, {G: 255}, {B: 255}},
			wantPixel: Pixel{G: 255},
		},
		{
			name:      "color blue",
			p:         Pixel{B: 255},
			dataset:   []Pixel{{R: 255}, {G: 255}, {B: 255}},
			wantPixel: Pixel{B: 255},
		},
		{
			name: "from red to lego red",
			p:    Pixel{R: 255, G: 0, B: 0},
			dataset: []Pixel{
				{R: 254, G: 138, B: 24},  // Orange
				{R: 200, G: 112, B: 160}, // Dark Pink
				{R: 201, G: 26, B: 9},    // Red
			},
			wantPixel: Pixel{R: 201, G: 26, B: 9},
		},
		{
			name:      "from red to lego red using dataset",
			p:         Pixel{R: 255, G: 0, B: 0},
			dataset:   legocolors,
			wantPixel: Pixel{R: 201, G: 26, B: 9}, // Red
		},
		{
			name:      "from green to lego green using dataset",
			p:         Pixel{R: 0, G: 255, B: 0},
			dataset:   legocolors,
			wantPixel: Pixel{R: 115, G: 220, B: 161}, // Medium Green
		},
		{
			name:      "from blue to lego blue using dataset",
			p:         Pixel{R: 0, G: 0, B: 255},
			dataset:   legocolors,
			wantPixel: Pixel{R: 0, G: 32, B: 160}, // Trans-Dark Blue
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mindistance := math.MaxFloat64
			var found int

			for i := range tt.dataset {
				distance := calculateDistance(tt.p, tt.dataset[i])
				if distance < mindistance {
					mindistance = distance
					found = i
				}
			}

			if tt.dataset[found] != tt.wantPixel {
				t.Errorf("wrong pixel found as closest: found=%d, expected=%d", tt.dataset[found], tt.wantPixel)
			}

		})
	}
}
