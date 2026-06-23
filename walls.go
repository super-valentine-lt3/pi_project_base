package main 

import (
    "github.com/elgopher/pi"
    //"github.com/lafriks/go-tiled"
   // "fmt"
)

type GameRect interface {
	GetAreaJ() pi.IntArea
}

func CanMove( world *World,  tempX int, tempY int) bool {
	for _, layer := range world.TileMap.Tiles {	   
	   for y := 0; y < len(layer); y++ {
	      for x := 0; x < len(layer[y]); x++ {
	      	tile := layer[y][x]
	      	//fmt.Println(tile)
	      	tempArea := pi.IntArea{tempX, tempY, 16, 16}
	      	if Intersects(tempArea, tile.GetArea()) && tile.Solid {
	      		return false 
	      	}
	      	//str := fmt.Sprintf("Name: %s IsSolid: %t \n", layer[y][x].SpriteName, layer[y][x].Solid)
	      	//fmt.Println(str)
	      }
	   }
	}
	return true 
}