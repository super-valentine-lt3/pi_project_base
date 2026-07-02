package main 

import (
	"github.com/elgopher/pi" 
	"github.com/elgopher/pi/pikey"
	"github.com/elgopher/pi/pievent"
	"github.com/lafriks/go-tiled"
	"embed"
)

//go:embed assets
var assetsFS embed.FS

var RoomPath string = "assets/rooms/"

var FirstRoom string = "room_0.tmx"

// GameScreen handles multiple Rooms 
type GameScreen struct {
	Room *Room 
	Paused bool 
	TileSet TileSet 
}

func LoadMap(MapPath string) *tiled.Map {
   Map, err := tiled.LoadFile(
        MapPath,
        tiled.WithFileSystem(assetsFS),
   )

   if err != nil {
       panic(err)
   }
   return Map 
}

func (g *GameScreen) Init() {
   firstMap := LoadMap(RoomPath + FirstRoom)
   g.TileSet = NewTileSet(firstMap) 
   g.Room = NewRoom(g.TileSet, nil,  firstMap)
   PlayTheme()
   pikey.Target().Subscribe(pikey.Event{pikey.EventDown, pikey.Esc}, g.TogglePause)
}

func (g *GameScreen) TogglePause(e pikey.Event, h pievent.Handler) {
	g.Paused = !g.Paused 	
}

func (g *GameScreen) Draw() {
	pi.Screen().Clear(32)
	g.Room.Draw()
	if g.Paused {
		pausePanelRoot.Draw()
	}
}

func (g *GameScreen) Update() {
	if !g.Paused  {
		newRoom := g.Room.Update() 
		if newRoom != nil && *newRoom != "title_screen" {
			Map := LoadMap(RoomPath + *newRoom)
			g.Room = NewRoom(g.TileSet, g.Room.World.Player, Map)
		} else if newRoom != nil && *newRoom == "title_screen" {
			SetScreen(&TitleScreen{}, false)
			StopTheme()
		}
	} else {
		pausePanelRoot.Update()
	}
}

type Room struct {
	World *World 
}

func NewRoom(tileSet TileSet, character *Character, GameMap *tiled.Map) *Room {
   tileMap := NewTileMap(GameMap)
   objectMap := NewObjectMap(GameMap)

   var Char *Character 
   player := objectMap.Objects["Player"][0]

   if character != nil {
   		Char = character
   		Char.GameObject = player
   } else {
		Char = NewCharacter(player, 
		     CharacterSpriteFile, 
		     CharacterSpriteDirectory, 
		     CharacterSpriteStartAnim)
		Char.SetAction("move_up", pikey.Up)
		Char.SetAction("move_left", pikey.Left)
		Char.SetAction("move_right", pikey.Right)
		Char.SetAction("move_down", pikey.Down)
		Char.SetAction("interact", pikey.Space)   
   }
   //Char.SetAction("drop_bomb", pikey.CtrlLeft)   

   bombs := make([]*Bomb, 0)
   for _, bomb := range objectMap.Objects["Bomb"] {
      bombs = append(bombs, NewBomb(bomb, BombSpriteFile, BombSpriteDirectory, BombSpriteStartAnim))
   }
   
   // multiple doors in next game 
   doorObj := objectMap.Objects["tile_door_3"][0]
   door := NewDoor(doorObj, tileSet.Tiles["tile_door_3"])

   gems := make([]*Gem, 0)
   for _, gem := range objectMap.Objects["Gem"] {
      gems = append(gems, NewGem(gem, tileSet.Tiles["tile_gem"]))
   }
   
   bats := make([]*Bat, 0)
   for _, bat := range objectMap.Objects["Bat"] {
      bats = append(bats, NewBat(bat))
   }

   crabs := make([]*Crab, 0)
   for _, crab := range objectMap.Objects["Crab"] {
      crabs = append(crabs, NewCrab(crab))
   }

   world := &World{Player: Char, 
   				   Bombs: bombs, 
   				   Door: &door, 
   				   Gems: gems, 
   				   Bats: bats, 
   				   Crabs: crabs, 
   				   TileMap: &tileMap}


	return &Room{World: world}
}

func (r *Room) Draw() {
	UISystem(r.World)

	DrawTileLayer(r.World.TileMap, "Tile Layer 1")
	DrawTileLayer(r.World.TileMap, "wallsides")

	r.World.Door.Draw() 

	for _, bomb := range r.World.Bombs {
		bomb.Draw()
	}
	for _, gem := range r.World.Gems {
		gem.Draw()
	}
	for _, bat := range r.World.Bats {
		bat.Draw()
	}
	for _, crab := range r.World.Crabs {
		crab.Draw()
	}

	for _, proj := range r.World.Projectiles {
		proj.Draw()
	}

	r.World.Player.Draw() 
}

func (r *Room) Update() (NewRoom *string) {
	BombSystem(r.World)
	NewRoom = DoorSystem(r.World)
	GemSystem(r.World)
	BatSystem(r.World)
	CrabSystem(r.World)
	ProjectileSystem(r.World)
	r.World.Player.Update(r.World) 
	return 
}