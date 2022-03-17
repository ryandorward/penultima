package data

import (
	"os"
	"fmt"
	// "log"
	"image"	
	_ "image/png"
	"time"
	// "path/filepath"
	"app/pkg/game/model"
)

/* 
@todo: convert this to work with type rawTiledMap from zones.go
*/

var WorldWidth, WorldHeight int

type WorldMap struct {
	Grid [][] int8
}

/*
func (w WorldMap) getTerrain(location Coord) int {
	return w.Grid[location.X][location.Y]
}
*/

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
			//	newCol[y] = png_pixel_to_terrain_code(imageArr[x][y]) // @todo switch this to use png_pixel_to_tile - and function will return [][]model.Tile
			newCol[y] = png_pixel_to_tile(imageArr[x][y])
		}		
		out[x] = newCol					
	}			

	return out
}

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


func png_pixel_to_terrain_code( pixel[3]uint32) int {

	switch pixel {
		case [3]uint32{4626, 32896, 39835}: // deep water
			return 1;	
		case [3]uint32{35980, 48316, 52171}: // medium water
			return 2;	
		case [3]uint32{50629, 60138, 62965}: // shallow water
			return 3;	
		case [3]uint32{42148, 53456, 37779}: // land
			return 5;
		case [3]uint32{63479, 61423, 49344}: // foothills
			return 6;
		case [3]uint32{46517, 44461, 32896}: // low mountain
			return 7;
		case [3]uint32{44204, 32382, 23644}: // high mountain
			return 8;
		case [3]uint32{39578, 44204, 13878}: // light forest
			return 9;
		case [3]uint32{32125, 32382, 7196}: // heavy forest
			return 10;
		case [3]uint32{65535, 65535, 13107}: // beach
			return 11;
		case [3]uint32{52428, 65535, 13107}: // desert - cactus
			return 12;
		case [3]uint32{62708, 55769, 23130}: // desert - sand	
			return 13;
		case [3]uint32{26214, 52428, 39321} : // marsh/wetland
			return 14;								
		default: 
			return 0; 
	}

}

func png_pixel_to_tile( pixel[3]uint32) model.Tile {

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
		default: 
			return Tiles["nothing"];	// 0
	}

}

func NewWorldMap(img_path string) ([][]model.Tile, int, int) {

	/*
	dir, err := filepath.Abs("./")
	if err != nil {
		log.Fatal(err)
	}

	img_path:= dir + "/data/zones/xoxaria.png" 
	*/

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
	
	/*
	for y := 0; y < 15; y++ { 
		for x := 0; x < 2; x++ {		   
			r := imageArr[x][y][0]
			g := imageArr[x][y][1]
			b := imageArr[x][y][2]
			fmt.Printf("y:%d R %d  G %d  B %d | ", y, r,g,b)
		}
		fmt.Printf("\n")
	}	
	*/

	terrainArr := imageArr_to_terrain_array(imageArr)

	for y := 0; y < 15; y++ { 
		for x := 0; x < 2; x++ {			  		
			fmt.Printf("y:%d Terrain: %d | ", y, terrainArr[x][y])
		}
		fmt.Printf("\n")
	}	
		 
	/*
	myWorldMap := WorldMap{}
	myWorldMap.Grid = terrainArr
	*/
	
	return terrainArr, WorldWidth, WorldHeight

}