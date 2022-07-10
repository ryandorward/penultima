package entity

import (
	gameModel "app/pkg/game/model"	
	_ "github.com/lib/pq"
	"math/rand"
	// "time"
	"app/pkg/game/action"
	"app/pkg/game/tiles"
	"github.com/google/uuid"
	"sort"
	"fmt"
)

type moveMemory struct {
	memory []int
}

const moveMemoryLength = 8

func newMoveMemory() *moveMemory {
	m := &moveMemory{
		memory: make([]int, moveMemoryLength),
	}
	for i := range m.memory { // initialize movement memory, all "blocked". Should it be all "still" aka 0?
    m.memory[i] = -1
	}
	return m
}

func (m *moveMemory) enqueue(d int) {	
	m.memory = append(m.memory, d)
	m.memory = m.memory[1:]
}	

// most recent "memory" is at index 0
func (m *moveMemory) read(i int) int { 
	if i < 0 || i >= len(m.memory) {
		return -1
	} else {	
		return m.memory[len(m.memory)-i-1] 	// read memory backwards
	}
}	

// convenience
func (m *moveMemory) last() int { 
	return m.read(0)	
}	

// NPC is an extension of entityData
type NPC struct {
	entityData // inheritance, basically. entityData comes from server/pkg/game/entity/entity.go
	startingZone gameModel.Zone	
	properties gameModel.NPCProperties
	moveMemory *moveMemory // Memory goes back 8 moves. This is used for various movement algos. 0-4 for cardinal directions, 0 for stay still, -1 for blocked/need initialization	
	movementTick int // a simple counter that increments every time the NPC moves. Used for movement algorithms to throttle speed and give NPCs a pace
	currentMovementMod int // the speed mod can vary depending on jitter
}

func NewNPC() *NPC {
	// initialize
	n := &NPC{
		entityData: entityData{					
			UUID: uuid.New(),
			EnergyThreshold: playerEnergyThreshold,
			slowThresh: slowThresh,
			Type: gameModel.EntityTypeNPC,		
		},	
		moveMemory: newMoveMemory(),
		currentMovementMod: 11, // init. This gets calculated more bespoke on every move cycle
	}	
	return n
}

func int2XYvec(i int) (int, int) {
	switch i {
		case 0:
			return 0, 0
		case 1:
			return 0, -1
		case 2:
			return 1, 0
		case 3:
			return 0, 1
		case 4:
			return -1, 0
	}
	return 0, 0
}

func getOppositeDirection(d int) int {
	switch d {
		case 1:
			return 3
		case 2:
			return 4
		case 3:
			return 1
		case 4:
			return 2
	}
	return 0
}


func isAdjacent (x1, y1, x2, y2 int) bool {
	return (x1 != x2 && (y1 == y2-1 || y1 == y2+1)) || (y1 != y2 && (x1 == x2-1 || x1 == x2+1))
}

func (n *NPC) moveDrunken() {	
	//randomly pick a move to make in one of the 4 cardinal directions, or stay still
	walkDir := rand.Intn(6) 
	x,y := int2XYvec(walkDir+1) // because 0 means stay still, and we want to start at 1 so that we only queue up moves in cardinal directions

	if walkDir < 4 {
		n.QueuedAction = &action.MoveAction{
			Mover: n,
			X: x,
			Y: y, 
		} 
	}
}	

// move in random straight lines until blocked and then pick a new direction
// see ReceiveResult for how blocked is handled, and how it affects this movement behaviour
func (n *NPC) moveBumpercar() {		
	rando := rand.Intn(100)
	// there is a small chance the NPC will change direction
	// @todo: make the threshold for direction change configurable in NPC settings
	if n.moveMemory.last() == -1 || rando < n.properties.Movement.DirectionChangeProbability { // initialize movement direction. This happens when NPC is spawned, or when it's just been blocked	
		n.moveMemory.enqueue(rand.Intn(4) + 1)
	}
	
	x,y := int2XYvec(n.moveMemory.last())
	n.QueuedAction = &action.MoveAction{
		Mover: n,
		X: x,
		Y: y, 
	}
	
}	


// walk randomly, but biasd towards coming back to anchor point
func (n *NPC) moveAnchoredWander() {	
	/*
	x, y := n.GetPosition()
	ax := n.properties.Movement.Anchor.X
	ay := n.properties.Movement.Anchor.Y
	*/
}

func chickenWalkTileWeight(tile *gameModel.Tile) float64 {
	if tile.ID >= 15 && tile.ID <= 18 { // bush
		return 2.0
	}
	if tile.ID == 9 {// light forest
		return 1.5
	}
	if tile.ID == 10 {// heavy forest
		return 1.2
	}
	if tile.ID ==  5 {// land
		return 1
	}
	if tile.ID ==  6 {// foothills
		return 0.8
	}
	if tile.Solid == false {
		return tile.Speed
	}
	return 0.0
}

// Wander around preferring bush, light forest, and grass in that order
func (n *NPC) moveChickenWalker() {	
	
	x, y := n.GetPosition()
	zone := n.GetZone() 
	lastMove := n.moveMemory.last()	
	oppositeOfLastMove := getOppositeDirection(lastMove)

	directions := []int{1,2,3,4}
	// sort the directions using own compare function
	sort.Slice(directions, func(i, j int) bool {
		dix, diy := int2XYvec(directions[i])
		djx, djy := int2XYvec(directions[j])		
		ix, iy := zone.GetNewLocation(x,y,dix,diy) 	// translate delta into a zone coord
		jx, jy := zone.GetNewLocation(x,y,djx,djy) 	// translate delta into a zone coord

		iOccupied := false
		jOccupied := false

		// if one of the directions is occupied and the other is not, always prefer unoccupied space
		for _, e := range zone.GetEntities() {
			eX, eY := e.GetPosition()
			if eX == ix && eY == iy {
				iOccupied = true
			}
			if eX == jx && eY == jy {
				jOccupied = true
			}
		}
		if iOccupied && !jOccupied {
			return false
		}
		if !iOccupied && jOccupied {
			return true
		}
		// prefer not reversing direction
		if directions[i] == oppositeOfLastMove {
			return false
		}
		if directions[j] == oppositeOfLastMove {
			return true
		}
		// now give tile a weight based on its type
		return chickenWalkTileWeight(zone.GetTile(ix, iy)) > chickenWalkTileWeight(zone.GetTile(jx, jy))
	})

	fmt.Println(directions)

	// now pop the top item off directions and use it
	for _, d := range directions {
		x, y := int2XYvec(d)
		n.QueuedAction = &action.MoveAction{
			Mover: n,
			X: x,
			Y: y, 
		}
		n.moveMemory.enqueue(d)
		return
	}

}


// Generally wander around following the path of least resistance. Some obstactle avoidance
func (n *NPC) moveEasyWalker() {		

	x, y := n.GetPosition()
	zone := n.GetZone() 
	lastMove := n.moveMemory.last()	
	oppositeOfLastMove := getOppositeDirection(lastMove)
		
	// find the direction with the fastest tile
	var fastestTileSpeed float64
	var fastestDirs []int

	// determine the set of directions that are fast enough to consider moving to
	fastestTileSpeed = -1
	DIRECTIONS:
	for i := 1; i <= 4; i++ {
		if i == oppositeOfLastMove { // never turn around backwards
			continue DIRECTIONS
		}
		dx,dy := int2XYvec(i)	
		nx, ny := zone.GetNewLocation(x,y,dx,dy) 	// translate delta into a zone coord

		// check for other entities in this location
		for _, otherE := range zone.GetEntities() {
			otherX, otherY := otherE.GetPosition()
			if n != otherE && otherX == nx && otherY == ny {		
				continue DIRECTIONS				
			}
		}

		tile := zone.GetTile(nx,ny)
		if tile.Speed > fastestTileSpeed {
			fastestDirs = []int{i}
			fastestTileSpeed = tile.Speed			
		}	else if tile.Speed == fastestTileSpeed {
			fastestDirs = append(fastestDirs, i)
		}
	}

	rando := rand.Intn(100)
 
	dir := 0 
	// if not at a crossroads, and can continue in the direction we're moving in, do so.
	// Unless we're changing directions randomly
	if len(fastestDirs) < 2 || rando > n.properties.Movement.DirectionChangeProbability {
		for _,d := range fastestDirs {
			if d == lastMove {
				dir = d
				break
			}
		}
	}

	// if blocked on last move, consider the one before that "last" so we can move perpendicular to it
	if lastMove == -1 {
		lastMove = n.moveMemory.read(1)
		// fmt.Println("lastMove: ", lastMove)
	}
 
	// Next priority is to pick a dir that is perpendicular to last move
	var adjacentDirs []int
	if dir == 0 {
		dmx, dmy := int2XYvec(lastMove)
		for _,d := range fastestDirs {
			dx,dy := int2XYvec(d)
			if isAdjacent ( dmx, dmy, dx, dy ) {
				adjacentDirs = append(adjacentDirs, d)
			}
		}
		// pick one of the adjacent directions randomly, if there are any
		if len(adjacentDirs) > 0 {			
			dir = adjacentDirs[rand.Intn(len(adjacentDirs))]					
		}
	}

	// otherwise pick at random one of the items from fastestDirs
	if 	len(fastestDirs) == 0 { // but only if there are fastestDirs
		return
	}
	if dir == 0 {		
		dir = fastestDirs[rand.Intn(len(fastestDirs))]		
	}
	
	x,y = int2XYvec(dir)
	n.moveMemory.enqueue(dir)

	n.QueuedAction = &action.MoveAction{
		Mover: n,
		X: x, 
		Y: y,
	}	
}	

func (n *NPC) Tick() bool {	
	n.Move()
	return n.entityData.Tick()
}

func (n *NPC) Move() {	
	
	// when NPC initialized, n.currentMovementMod is set to 11. After the first move it gets set based on the SpeedMod property + jitter
	n.movementTick = (n.movementTick + 1) % n.currentMovementMod

	if n.movementTick == 0 {		
		// vary the speed of the movement by a random amount +/- n.properties.Movement.Jitter
		if n.properties.Movement.SpeedMod != 0 {
			n.currentMovementMod = n.properties.Movement.SpeedMod
		}
		if n.properties.Movement.Jitter != 0 {
			jitterModifier := rand.Intn(n.properties.Movement.Jitter*2) - n.properties.Movement.Jitter
			if jitterModifier > 0 {
				jitterModifier += jitterModifier // because speer is 1/n, it has more impact when jitter is negative. So we double +ve jitter
			}
			n.currentMovementMod += jitterModifier
		}	
		switch n.properties.Movement.Algorithm {
			case "drunken":
				n.moveDrunken()
			case "bumpercar":
				n.moveBumpercar()	
			case "easyWalker":
				n.moveEasyWalker()
			case "chickenWalker":
				n.moveChickenWalker()	
		}
	}
}

func (n *NPC) Act() gameModel.Action {
	a := n.QueuedAction
	n.QueuedAction = nil
	return a
}

func (n *NPC) ReceiveMessage(m string) string{
	return "hi"
}

func (n *NPC) SetProperties(p gameModel.NPCProperties) {
	n.SetName(p.Name)
	n.SetPosition(p.X, p.Y)
	n.SetTile(tiles.Tiles[p.Tile].ID)
	if (p.Health != 0) {
		n.Stats.HP = float64(p.Health)
	} else {
		n.Stats.HP = 10.0
	}	
	if (p.Type != "") {
		n.Type = p.Type
	}
	n.properties = p
}

func (n *NPC) ReceiveResult(msg string, code string) {
	if code == "blocked" {
		switch n.properties.Movement.Algorithm {		
			case "bumpercar":
			case "easyWalker":			
				n.moveMemory.enqueue(-1)
		}
	}
}

// TakeDamage returns if they would die so XP can be dished out
func (n *NPC) TakeDamage(damage float64) bool {	
	if (n.properties.IsMortal == true) {
		return n.entityData.TakeDamage(damage)
	} 
	return false
}