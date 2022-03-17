package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"app/pkg/game/model"
	"app/pkg/game/zone"
	"github.com/google/uuid"

)

/* 

@todo: change so that Laers.Data is shaped like 2-d tile map, as in old imageMap

*/
type rawTiledMap struct {
	Width      int `json:"width"`
	Height     int `json:"height"`
	Properties []struct {
		Name  string          `json:"name"`
		Type  string          `json:"type"`
		Value json.RawMessage `json:"value"`
	} `json:"properties"`
	Layers []struct {
		ID      int    `json:"id"`
		Width   int    `json:"width"`
		Height  int    `json:"height"`
		Data    []int  `json:"data"`
		Name    string `json:"name"`
		Type    string `json:"type"`
		Subtype string `json:"subtype"`
		File 		string `json:"file"`
		X       int    `json:"x"`
		Y       int    `json:"y"`
		Opacity int    `json:"opacity"`
		Visible bool   `json:"visible"`
		Objects []struct {
			Name       string `json:"name"`
			X          int    `json:"x"`
			Y          int    `json:"y"`
			Width      int    `json:"width"`
			Height     int    `json:"height"`
			Type       string `json:"type"`
			TileID     int    `json:"gid"`
			Properties []struct {
				Name  string          `json:"name"`
				Type  string          `json:"type"`
				Value json.RawMessage `json:"value"`
			} `json:"properties"`
		} `json:"objects"`
	} `json:"layers"`
}

func LoadZones() map[string]*zone.Zone {
	zones := map[string]*zone.Zone{}

	files, err := ioutil.ReadDir("./data/zones")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		split := strings.Split(file.Name(), ".")
		name, ext := split[0], split[1]
		if ext == "json" {
			// zoneUUID := uuid.MustParse(name)
			zoneName := name
			zones[zoneName] = loadTiledMap(zoneName) // = loadTiledMap(zoneUUID)
		}
		
	}

	for _, z := range zones {
		for _, obj := range z.WorldObjects {
			if obj.WarpTarget != nil {
				obj.WarpTarget.Zone = zones[obj.WarpTarget.ZoneName] // tie warp targets to zones via names
			}
		}
	}

	return zones
}

/*
func loadTiledImageMap(file string) {

		// load the tile map from layer.File - use imageMap.go
		imageMapFile, err := os.Open(fmt.Sprintf("../data/zones/%s", file))
		if err != nil {
			log.Fatal(err)
		}
		defer imageMapFile.Close() 

}
*/

// func loadTiledMap(mapUUID uuid.UUID) *zone.Zone {
func loadTiledMap(mapName string) *zone.Zone {
	mapFile, err := os.Open(fmt.Sprintf("./data/zones/%s.json", mapName))
	if err != nil {
		log.Fatal(err)
	}
	defer mapFile.Close()

	rawData, err := ioutil.ReadAll(mapFile)
	if err != nil {
		log.Fatal(err)
	}

	var mapData rawTiledMap
	err = json.Unmarshal(rawData, &mapData)
	if err != nil {
		log.Fatal(err)
	}

	z := zone.Zone{
		Name:   mapName,
		Width:  mapData.Width,
		Height: mapData.Height,
		Tiles:  [][]model.Tile{},

		Torroidal: false,

		Entities:     map[string]model.Entity{},
		WorldObjects: map[string]*model.WorldObject{},
	} 

	for _, property := range mapData.Properties {
		if property.Name == "name" {
			err := json.Unmarshal(property.Value, &z.Name)
			if err != nil {
				log.Fatal(err)
			}
		}
		if property.Name == "torroidal" {
			err := json.Unmarshal(property.Value, &z.Torroidal)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, layer := range mapData.Layers {
		if layer.Name == "ground" {
			if layer.Type == "tilelayer" && layer.Subtype == "imageMap"{				
				z.Tiles, z.Width, z.Height = NewWorldMap("./data/zones/" + layer.File) 
			} else { // tile map is embedded in file @todo 
				/*
				for _, tileID := range layer.Data { 
					z.Tiles = append(z.Tiles, Tiles[tileID-1]) // -1 because of air tile (TODO: add air tile to -1 or something)
				}
				*/
			}
		}
		if layer.Name == "world_objects" { 
			for _, obj := range layer.Objects {
				var UUID uuid.UUID // try to get rid of UUID
				var Name string

				var hasWarpTarget, hasFullHeal bool
				var warpTargetName string
				var warpTargetX, warpTargetY int

				fmt.Println(obj.Name)

				for _, prop := range obj.Properties {			
					
					// fmt.Println("Prop name:", prop.Name)
				
					if prop.Name == "Name" {						
						err := json.Unmarshal(prop.Value, &Name)
						if err != nil {
							log.Fatal(err)
						}
					}
				
					if prop.Name == "warpTargetName" {
						hasWarpTarget = true
						err := json.Unmarshal(prop.Value, &warpTargetName)
						if err != nil {
							log.Fatal(err)
						}
					}
					if prop.Name == "warpTargetX" {
						err := json.Unmarshal(prop.Value, &warpTargetX)
						if err != nil {
							log.Fatal(err)
						}
					}
					if prop.Name == "warpTargetY" {
						err := json.Unmarshal(prop.Value, &warpTargetY)
						if err != nil {
							log.Fatal(err)
						}
					}
					if prop.Name == "fullHeal" {
						err := json.Unmarshal(prop.Value, &hasFullHeal)
						if err != nil {
							log.Fatal(err)
						}
					}
				}

				// fmt.Println("Name:", Name)

				z.WorldObjects[Name] = &model.WorldObject{
					UUID: UUID,
					Name: obj.Name,
					Tile: obj.TileID - 1,
					X:    obj.X / obj.Width,
					Y:    (obj.Y / obj.Height) - 1, // minus 1 because tiled objects start at the bottom left, tiles are top level (why the hell)
					Type: model.WorldObjectType(obj.Type),
				}
				if hasWarpTarget {
					z.WorldObjects[Name].WarpTarget = &model.WarpTarget{
						ZoneName: warpTargetName,
						X:        warpTargetX,
						Y:        warpTargetY,
					}
				}
				if hasFullHeal {
					z.WorldObjects[Name].HealZone = &model.HealZone{
						Full: true,
					}
				}
			}
		}
	}

	// TODO: worldObjects (either create from props, or object layer)

	return &z
}
