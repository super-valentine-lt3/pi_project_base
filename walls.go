package main 

import (
    "github.com/elgopher/pi"
    //"github.com/lafriks/go-tiled"
   "fmt"
)

func (s Direction) String() string {
	switch (s) {
		case Up:{ return "up" }
		case Down: {return "down"}
		case Left: {return "left"}
		case Right: {return "right"}
		case UpRight: {return "up right"}
		case UpLeft: {return "up left"}
		case DownRight: {return "down right"}
		case DownLeft: {return "down left"}
		default: { return "invalid"}
	}
}
type GameRect interface {
	GetArea() pi.IntArea
}

func CanMove( world *World,  tempX int, tempY int) bool {
	//tempArea := pi.IntArea{tempX, tempY, 16, 16}
	tempArea := PlayerCollisionArea(tempX, tempY)
	
	for _, layer := range world.TileMap.Tiles {	   
	   for y := 0; y < len(layer); y++ {
	      for x := 0; x < len(layer[y]); x++ {
	      	tile := layer[y][x]
	      	tileArea1, tileArea2, _ := tile.GetArea() 
	      	if Intersects(tempArea, tileArea1) && tile.Solid {
	      		return false 
	      	}
	      	if tileArea2 != nil && Intersects(tempArea, *tileArea2) && tile.Solid {
	      		return false 
	      	}
	      	//str := fmt.Sprintf("Name: %s IsSolid: %t \n", layer[y][x].SpriteName, layer[y][x].Solid)
	      	//fmt.Println(str)
	      }
	   }
	}

	if world.Door != nil && Intersects(tempArea, world.Door.GetArea()) {
		return false 
	}
	return true 
}