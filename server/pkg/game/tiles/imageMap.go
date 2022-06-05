package tiles

import (
	"os"
	"fmt"
	// "log"
	"image"	
	_ "image/png"
	"time"
	// "path/filepath"
	"app/pkg/game/model"
	"app/pkg/game/util"
)

var WorldWidth, WorldHeight int

type WorldMap struct {
	Grid [][] int
}

// @todo: using pix method would be 4-10x faster, see: https://stackoverflow.com/questions/33186783/get-a-pixel-array-from-from-golang-image-image
func image_to_array(img image.Image) [][][3]uint32 {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	out := make([][][3]uint32,height)

	for x := 0; x < width; x++ {	
		col := make([][3]uint32, width)
		for y := 0; y < height; y++ {
			// color := img.At(x, y)
			r,g,b,_ := img.At(x, y).RGBA() // discard transparency channel        
			col[y] = [3]uint32{r,g,b}
    }
		out[x] = col
	}
	return out
}

func imageArr_to_terrain_array(imageArr [][][3]uint32 ) [][]model.Tile {
	
	height := len(imageArr)
	width := len(imageArr[0])

	out := make([][]model.Tile,height)
	
	for x := 0; x < width; x++ {	
		newCol := make([]model.Tile, width)
		for y := 0; y < height; y++ {
			newCol[y] = png_pixel_to_tile(imageArr[x][y],x,y)
		}		
		out[x] = newCol					
	}	

	out = contextualize_tiles(out)
	
	return out
}

func contextualize_tiles(tiles [][]model.Tile) [][]model.Tile {
	
	height := len(tiles)
	width := len(tiles[0])

	out := make([][]model.Tile,height)
	
	for x := 0; x < width; x++ {	
		newCol := make([]model.Tile, width)
		for y := 0; y < height; y++ {
			// neightbouring tiles
			// nbrN := tiles[x][util.WrapMod(y-1,width)]
			//nbrE := tiles[util.WrapMod(x+1,height)][y]
			nbrS := tiles[x][util.WrapMod(y+1,width)]
			// nbrW := tiles[util.WrapMod(x-1,height)][y]

			tile := tiles[x][y]


			tileN := tiles[x][util.WrapMod(y-1,width)] // assume torroidal hope it doesn't get weird
			tileE := tiles[util.WrapMod(x+1,height)][y]
			tileS := tiles[x][util.WrapMod(y+1,width)]
			tileW := tiles[util.WrapMod(x-1,height)][y]

			if (tile == Tiles["land"] && nbrS == Tiles["shallow_water"]) {
				newCol[y] = tiles[x][y]
			}	else {
				newCol[y] = tiles[x][y]
			}			

			if (tile == Tiles["pathEW"]) {
				// 	if (tileE == Tiles["pathEW"] && tileW == Tiles["pathEW"]) { // do nothing

				pathCount := 	util.Bool2Int(tileN == Tiles["pathEW"]) + util.Bool2Int(tileE == Tiles["pathEW"]) + util.Bool2Int(tileS == Tiles["pathEW"]) + util.Bool2Int(tileW == Tiles["pathEW"])

				if (pathCount > 2) { 
					newCol[y] = Tiles["pathIntersection"]
				} else if (tileW == Tiles["pathEW"] && tileN == Tiles["pathEW"]) { 
					newCol[y] = Tiles["pathNW"]
				} else if (tileW == Tiles["pathEW"] && tileS == Tiles["pathEW"]) {
					newCol[y] = Tiles["pathSW"]
				} else if (tileE == Tiles["pathEW"] && tileN == Tiles["pathEW"]) {
					newCol[y] = Tiles["pathNE"]
				} else if (tileE == Tiles["pathEW"] && tileS == Tiles["pathEW"]) {
					newCol[y] = Tiles["pathSE"]
				} else if (tileE == Tiles["pathEW"] && tileS == Tiles["pathEW"]) {
					newCol[y] = Tiles["pathSE"]
				} else if (tileN == Tiles["pathEW"] && tileS == Tiles["pathEW"]) {
					newCol[y] = Tiles["pathNS"]
				} else if (tileS == Tiles["pathEW"] || tileN == Tiles["pathEW"]) {				
					newCol[y] = Tiles["pathNS"]
				}

			}

		}		
		out[x] = newCol					
	} 
	
	return out
}

/*
func get_tile_opacity (tile int) float32 {
	switch tile {
		case 8: // high mountain
		case 10: // heavy forest
			return 0
		default:
			return 1
	}
	return 1;
}
*/


func cantorPairing(k1, k2 int) int {
	return (k1+k2)*(k1+k2+1)/2 + k2
}

func png_pixel_to_tile( pixel[3]uint32, x, y int) model.Tile {

	switch pixel {
		case [3]uint32{4626, 32896, 39835}: // deep water
			return Tiles["deep_water"];	
		case [3]uint32{35980, 48316, 52171}: // medium water
			return Tiles["medium_water"];		
		case [3]uint32{50629, 60138, 62965}: // shallow water
			return Tiles["shallow_water"];	
		case [3]uint32{42148, 53456, 37779}: // land
			return Tiles["land"];	
		case [3]uint32{63479, 61423, 49344}: // foothills
			return Tiles["foothills"]; // 6
		case [3]uint32{46517, 44461, 32896}: // low mountain
			return Tiles["low_mountain"]; // 7;
		case [3]uint32{44204, 32382, 23644}: // high mountain
			return Tiles["high_mountain"]; //8;
		case [3]uint32{39578, 44204, 13878}: // light forest
			return Tiles["light_forest"]; //9;
		case [3]uint32{32125, 32382, 7196}: // heavy forest
			return Tiles["heavy_forest"]; //10;
		case [3]uint32{65535, 65535, 13107}: // beach
			return Tiles["beach"]; //11;
		case [3]uint32{52428, 65535, 13107}: // desert - cactus
			return Tiles["desert_cactus"]; //12;
		case [3]uint32{62708, 55769, 23130}: // desert - sand	
			return Tiles["desert_sand"]; //13;
		case [3]uint32{26214, 52428, 39321} : // marsh/wetland
			return Tiles["marsh"]; //14;		
		case [3]uint32{40863, 48830, 25957}	:
			cp := cantorPairing(x, y)
			m := cp % 4
			if m == 0 {
				return Tiles["brush"]; //15;
			} else if m == 1 {
				return Tiles["brush2"]; //16;
			} else if m == 2  {
				return Tiles["brush3"]; //17;
			} else { 
				return Tiles["brush4"]; //18;
			}
		case [3]uint32{43690, 21845, 0}	:
			return Tiles["pathEW"]; // @todo generate the context-relavant path tiles
		case [3]uint32{31097, 35980, 30069}	:
			return Tiles["stones"];
		case [3]uint32{21845, 21845, 21845}	:
			return Tiles["rockWall"]; // @todo generate the context-relavant path tiles
		case [3]uint32{43690, 43690, 43690} :
			return Tiles["stoneWall"];	
		case [3]uint32{43690, 0, 0}:
			return Tiles["tileFloor"];	
		case [3]uint32{25443, 13107, 0}:
			return Tiles["woodFloorHoriz"];			
		default: 
			fmt.Println("unknown pixel: ", pixel)
			return Tiles["nothing"];	// 0
	}

}



func NewWorldMap(img_path string) ([][]model.Tile, int, int) {

	if _, err := os.Stat(img_path); os.IsNotExist(err) {
		fmt.Println("Map file does not exist")
	}

	aa := time.Now()
	file, err := os.Open(img_path)
	if err != nil {
		return nil, 0, 0
	}
	defer file.Close()
	m, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Image decode error")
		return nil, 0, 0
	}	
	bounds := m.Bounds()
	WorldWidth = bounds.Max.X
	WorldHeight = bounds.Max.Y

	bb := time.Now()
  fmt.Println("Read file time: ", float64(bb.Nanosecond() - aa.Nanosecond()) / 1e9)

	aa = time.Now()
	imageArr := image_to_array(m)
	bb = time.Now()
	fmt.Println("Image to array Time: ", float64(bb.Nanosecond() - aa.Nanosecond()) / 1e9)
	
	terrainArr := imageArr_to_terrain_array(imageArr)

	for y := 0; y < 15; y++ { 
		for x := 0; x < 2; x++ {			  		
			// fmt.Printf("y:%d Terrain: %d | ", y, terrainArr[x][y])
		}
		fmt.Printf("\n")
	}	
		 	
	return terrainArr, WorldWidth, WorldHeight

}