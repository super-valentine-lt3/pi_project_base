package main 

import (
   "github.com/elgopher/pi/pikey"
   "github.com/elgopher/pi/pievent"
   "github.com/lafriks/go-tiled"
)

var RoomPath string = "assets/rooms/"

var FirstRoom string = "room_0.tmx"

// GameScreen handles multiple Rooms 
type GameScreen struct {
	Room *Room 
	Paused bool 
	TileSet TileSet 
}


func (g *GameScreen) Init() {
   firstMap, err := tiled.LoadFile(mapPath)
   if err != nil {
       panic(err)
   }

   g.TileSet = NewTileSet(firstMap) 
   pikey.Target().Subscribe(pikey.Event{pikey.EventDown, pikey.Esc}, g.TogglePause)
}

func (g *GameScreen) TogglePause(e pikey.Event, h pievent.Handler) {
	g.Paused = !g.Paused 	
}

func (g *GameScreen) Draw() {
	g.Room.Draw()
}

func (g *GameScreen) Update() {
	if !g.Paused  {
		g.Room.Update() 
	}
}

type Room struct {
	World *World 
}

func NewRoom() *Room {
	return &Room{}
}

func (r *Room) Draw() {}
func (r *Room) Update() {}