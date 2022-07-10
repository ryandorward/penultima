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

	if a.Peerer.GetGemCount() < 1 {
		a.Peerer.GetClient().In <- network.NewServerResultEvent("Peer at a Gem. You don't have any!", "fail")
	}

	viewWidth := peerGemWidth
	viewHeight := peerGemHeight																				
	fov := make([][]int, viewHeight) // initialize a slice of viewHeight slices		
	for i:=0; i < viewWidth; i++ {					
		fov[i] = make([]int, viewWidth) // initialize a slice of viewWidth in in each of viewHeight slices
	}
																									
	entityX, entityY := a.Peerer.GetPosition()
	zone := a.Peerer.GetZone()
	zoneWidth, zoneHeight := zone.GetDimensions()				
	halfViewWidth := viewWidth / 2
	halfViewHeight := viewHeight / 2
	
	for x := 0; x < viewWidth; x++ {	
		for y := 0; y < viewHeight; y++ {		
			nX := entityX+x - halfViewWidth
			nY := entityY+y - halfViewHeight
			if zone.GetTorroidal() {
				nX = util.WrapMod(nX, zoneWidth)
				nY = util.WrapMod(nY, zoneHeight)														
			}		
			if nX < 0 || nY < 0 || nX >= zoneWidth || nY >= zoneHeight {												
				fov[x][y] = 0
			} else {
				fov[x][y] = a.Peerer.GetZone().GetTile(nX, nY).ID					
			}
		} 
	}
	
	// a.Peerer.GetClient().In <- network.NewServerMessageEvent("> Peer at a Gem.")	
	a.Peerer.GetClient().In <- network.NewServerResultEvent("Peer at a Gem.", "success")
	a.Peerer.GetClient().In <- network.NewPeerGemEvent(&fov)
	a.Peerer.AddGems(-1) // one less gem
	a.Peerer.GetClient().In <- network.NewServerMessageEvent("Press any key to exit.")
	
	return true
}