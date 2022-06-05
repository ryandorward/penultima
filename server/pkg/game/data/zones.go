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
	"app/pkg/game/tiles"
	"app/pkg/game/entity"

)

type rawTiledMap struct {
	Width      int `json:"width"`
	Height     int `json:"height"`
	Type string `json:"type"`
	ParentZoneName string `json:"parentZoneName"`
	Properties []struct { 
		Name  string          `json:"name"`
		Type  string          `json:"type"`
		Value json.RawMessage `json:"value"`
	} `json:"properties"`
	GrowsFood bool `json:"growsFood"`
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
			Tile     string    `json:"tileName"`
			LightRadius int		`json:"lightRadius"`
			WarpTargetName string `json:"warpTargetName"`
			WarpTargetX int `json:"warpTargetX"`
			WarpTargetY int `json:"warpTargetY"`
			WarpTarget struct {
				Name string `json:"name"`
				X    int    `json:"x"`
				Y    int    `json:"y"`
			}
			Properties []struct {
				Name  string          `json:"name"`
				Type  string          `json:"type"`
				Value json.RawMessage `json:"value"`
			} `json:"properties"`
		} `json:"objects"`
	} `json:"layers"`
	NPCs []model.NPCProperties `json:"npcs"`
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

		// tie zones that are nested to their parent (e.g. towns)
		if (z.ParentZoneName != "") {
			for _, oz := range zones {
				if oz.Name == z.ParentZoneName {
					z.ParentZone = oz			
					// find x,y coords of this zone in parent zone
				}
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
		Type: mapData.Type,
		ParentZoneName: mapData.ParentZoneName,
		Torroidal: false,
		Entities:     map[string]model.Entity{},
		WorldObjects: map[uuid.UUID]*model.WorldObject{},
		NPCs:     map[string]model.Entity{},
		GrowsFood: mapData.GrowsFood,
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
				z.Tiles, z.Width, z.Height = tiles.NewWorldMap("./data/zones/" + layer.File) 
			} else { // tile map is embedded in file @todo nah don't bother
				/*
				for _, tileID := range layer.Data { 
					z.Tiles = append(z.Tiles, Tiles[tileID-1]) // -1 because of air tile (TODO: add air tile to -1 or something)
				}
				*/
			}
		}
		if layer.Name == "world_objects" { 
			for _, obj := range layer.Objects {
				
				var Name string
				UUID := uuid.New()

				// var hasWarpTarget bool
				var hasFullHeal bool
				// var warpTargetName string
				// var warpTargetX, warpTargetY int
			
				for _, prop := range obj.Properties {			
													
					if prop.Name == "Name" {						
						err := json.Unmarshal(prop.Value, &Name)
						if err != nil {
							log.Fatal(err)
						}
					}
				
					/*
					if prop.Name == "WarpTargetName" {
					// 	hasWarpTarget = true
						fmt.Println("has warp target")
						err := json.Unmarshal(prop.Value, &warpTargetName)
						if err != nil {
							log.Fatal(err)
						}
					}
				
					if prop.Name == "WarpTargetX" {
						err := json.Unmarshal(prop.Value, &warpTargetX)
						if err != nil {
							log.Fatal(err)
						}
					}
					if prop.Name == "WarpTargetY" {
						err := json.Unmarshal(prop.Value, &warpTargetY)
						if err != nil {
							log.Fatal(err)
						}
					}
					*/
					if prop.Name == "fullHeal" {
						err := json.Unmarshal(prop.Value, &hasFullHeal)
						if err != nil {
							log.Fatal(err)
						}
					}
				}
							
				//	z.WorldObjects[Name] = &model.WorldObject{
				z.WorldObjects[UUID] = &model.WorldObject{
					UUID: UUID,
					Name: obj.Name,
					// Tile: obj.TileID - 1,
					Tile: obj.Tile,
					X:    obj.X,
					Y:    obj.Y,
					Type: model.WorldObjectType(obj.Type),
					LightRadius: obj.LightRadius,
				}
			
				if (obj.WarpTarget.Name != "") {													
					z.WorldObjects[UUID].WarpTarget = &model.WarpTarget{
						ZoneName: obj.WarpTarget.Name,
						X:        obj.WarpTarget.X,
						Y:        obj.WarpTarget.Y,
					}											
				}
			
				if hasFullHeal {
					z.WorldObjects[UUID].HealZone = &model.HealZone{
						Full: true,
					}
				}
			}
		}
	}

	
	for _, npc := range mapData.NPCs {

		n := entity.NewNPC()


		/*
		n.SetName(npc.Name)
		n.SetPosition(npc.X, npc.Y)
		n.SetTile(tiles.Tiles[npc.Tile].ID)	
		*/
		n.SetProperties(npc)

		//	n.SetType(npc.Type)
		fmt.Println(npc.Name)
		z.Entities[npc.Name] = n
		z.AddEntity(n)

	}

	// TODO: worldObjects (either create from props, or object layer)

	return &z
}
