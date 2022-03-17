package action
import (
	"app/pkg/game/model"
	// "fmt"
	"app/pkg/game/event/network"
	"app/pkg/game/util"
)

type PeerGemAction struct {
	Peerer model.Entity
	Id  int
}

const peerGemWidth = 160
const peerGemHeight = 160

func (a *PeerGemAction) Execute() bool {

	viewWidth := peerGemWidth
	viewHeight := peerGemHeight																				
	fov := make([][]int8, viewHeight) // initialize a slice of viewHeight slices		
	for i:=0; i < viewWidth; i++ {					
		fov[i] = make([]int8, viewWidth) // initialize a slice of viewWidth in in each of viewHeight slices
	}
																									
	entityX, entityY := a.Peerer.GetPosition()
	zoneWidth, zoneHeight := a.Peerer.GetZone().GetDimensions()				
	halfViewWidth := viewWidth / 2
	halfViewHeight := viewHeight / 2
	
	for x := 0; x < viewWidth; x++ {	
		for y := 0; y < viewHeight; y++ {	
			nX := util.WrapMod(entityX+x - halfViewWidth, zoneWidth)
			nY := util.WrapMod(entityY+y - halfViewHeight, zoneHeight)										
			fov[x][y] = int8(a.Peerer.GetZone().GetTile(nX, nY).ID)
		} 
	}

	a.Peerer.GetClient().In <- network.NewServerMessageEvent("Peer Gem!")
	a.Peerer.GetClient().In <- network.NewPeerGemEvent(&fov)
	a.Peerer.GetClient().In <- network.NewServerMessageEvent("Press any key to exit.")
	
	return true
}