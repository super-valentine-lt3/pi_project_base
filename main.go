package main
   
import (
   "github.com/elgopher/pi"          // import pi core package
   "github.com/elgopher/pi/picofont" // import very small pico-8 font
   "github.com/elgopher/pi/piebiten" // import backend
    "github.com/lafriks/go-tiled"
    _ "embed"
    "fmt"
   "github.com/elgopher/pi/pikey"
)

const mapPath = "assets/room_test_1.tmx" // Path to your Tiled Map.
var gameMap *tiled.Map 
var tileSet TileSet 

type GameObject struct {
   Config *tiled.Object 
   Pos pi.Position 
}

type ObjectMap struct {
   Width int 
   Height int 
   Objects map[string][]GameObject  
}

func NewObjectMap() ObjectMap {
   objects := make(map[string][]GameObject)

   for _, layer := range gameMap.ObjectGroups {
       for _, object := range layer.Objects  {
         
         var name string 
         if object.Type != "" {
            if object.Name == "Door" { fmt.Println("I'm here 1.")}
            name = object.Type 
         } else if object.Template.Object.Type != "" {
            name = object.Template.Object.Type
                        if object.Name == "Door" { fmt.Println("I'm here 2.")}
         } else {
            name =  object.Template.Object.Properties.GetString("sprite")
                        if object.Name == "Door" { fmt.Println("I'm here 3.")}
            fmt.Println(name)
         }
         position := pi.Position{X: int(object.X), Y: int(object.Y)- gameMap.TileHeight}
         gameObj := GameObject{object, position}
         posList, ok := objects[name]
         if !ok {
            posList = make([]GameObject, 0)
         }
         posList = append(posList, gameObj)
         objects[name ] = posList 
       }
   }

   return ObjectMap {
      Width: gameMap.Width, 
      Height: gameMap.Height, Objects: objects}
}

type Tile struct {
   SpriteName string 
   Solid bool 
   X, Y int 
   Side *Direction 
}

func (t *Tile) GetArea() pi.IntArea {
   if t.Side == nil {
      return pi.IntArea{t.X, t.Y, 16, 16}
   } else {
      if *t.Side == Left {
         return pi.IntArea{t.X, t.Y, 8, 16}
      } else if *t.Side == Right {
         return pi.IntArea{t.X+8, t.Y, 8, 16}
      } else if *t.Side == Down {
         return pi.IntArea{t.X, t.Y+8, 16, 8}
      } else {
         return pi.IntArea{t.X, t.Y, 16, 16}
      }
   }
}

type TileMap struct {
    Width  int
    Height int
    //Solid  [][]bool
    Tiles map[string][][]Tile
}
func NewTileMap() TileMap {
   tiles := make(map[string][][]Tile)


   for _, layer := range gameMap.Layers {

      tileLayer := make([][]Tile, gameMap.Height)

      for y := range tileLayer {
          tileLayer[y] = make([]Tile, gameMap.Width)
      }  

      tiles[layer.Name] = tileLayer
       for pos, tile := range layer.Tiles {
           if tile.Nil {
               continue
           }

            tt, err := tile.Tileset.GetTilesetTile(tile.ID)
           if err != nil { continue }
           //if tt.Properties.GetBool("solid") {
               x := pos % gameMap.Width
               y := pos / gameMap.Width
               fmt.Printf("X: %d Y: %d\n", x, y)


               hasSide := len(tt.Properties.Get("side")) > 0 

               var Side *Direction 
               if hasSide {
                  sideString := tt.Properties.GetString("side")
                  if sideString == "left" {
                     side := Left
                     Side = &side
                  } else if sideString == "right" {
                     side := Right
                     Side = &side
                  } else if sideString == "up" {
                     side := Up 
                     Side = &side
                  } else if sideString == "down" {
                     side := Down 
                     Side = &side
                  }
               }
               tileLayer[y][x] = Tile{tt.Type, tt.Properties.GetBool("solid"), x*gameMap.TileWidth, y*gameMap.TileHeight, Side} 
           //}
       }
   }

   return TileMap {
      Width: gameMap.Width, 
      Height: gameMap.Height, Tiles: tiles}
}

//go:embed "assets/tiny_dungeon_tilesheet.png"
var spritesPNG []byte

//go:embed "assets/character_try_16x16_indexed.png"
var characterSpritesPNG []byte
const CharacterSpriteFile = "character_try_16x16_indexed.json"
const CharacterSpriteDirectory = "./assets"
const CharacterSpriteStartAnim = "idle_down "

//go:embed "assets/bomb_explode.png"
var bombSpritesPNG []byte
const BombSpriteFile = "bomb_explode.json"
const BombSpriteDirectory = "./assets"
const BombSpriteStartAnim = "normal"

func init() {
   var err error

    // Parse .tmx file.
    gameMap, err = tiled.LoadFile(mapPath)
    if err != nil {
        panic(err)
    }
    tileSet = TileSet{}
    tileSet.Tiles = make(map[string]pi.Sprite)
}

type TileSet struct {
   Tiles map[string]pi.Sprite 
}

func GetObjectGroup(name string) *tiled.ObjectGroup {
   for _, objectGroup := range gameMap.ObjectGroups {
      if objectGroup.Name == name {
         return objectGroup 
      }
   }
   return nil 
}

func GetObjectFromObjectLayer(objectGroup *tiled.ObjectGroup, name string) *tiled.Object{
   for _, object := range objectGroup.Objects {
      if object.Name == name {
         return object 
      }
   }  
   return nil 
}

func DrawTileLayer(tileMap *TileMap, layerName string) {
     layer := tileMap.Tiles[layerName]
   // Drawing Tiles 
   for y := 0; y < len(layer); y++ {
      for x := 0; x < len(layer[y]); x++ {
            pi.DrawSprite(tileSet.Tiles[layer[y][x].SpriteName], x*gameMap.TileWidth, y*gameMap.TileHeight)
      }
   }
}

type World struct {
   Player *Character 
   Bombs []*Bomb 
   Door *Door 
   TileMap *TileMap 
}

func Intersects(a, b pi.IntArea) bool {
    return ((a.X < b.X + b.W) &&
        (a.X + a.W > b.X) &&
        (a.Y < b.Y + b.H) && (a.Y + a.H > b.Y))
    
}

func main() {
   
   pi.Palette = pi.DecodePalette(spritesPNG)
   pi.SetTransparency(0, false)
   pi.SetTransparency(32, true)
   sprites := pi.DecodeCanvas(spritesPNG)

   // getting the tiles 
   tileList := gameMap.Tilesets[0].Tiles
   for _, tile := range tileList {
      r  := gameMap.Tilesets[0].GetTileRect(tile.ID)
      x, y, width, height := r.Min.X, r.Min.Y, r.Dx(), r.Dy()
      tileSet.Tiles[tile.Type] = pi.SpriteFrom(sprites, x, y, width, height)
   }
   fmt.Println(tileSet.Tiles)
   tileMap := NewTileMap()
   objectMap := NewObjectMap()

   fmt.Println(objectMap.Objects)
   pi.SetScreenSize(256, 144) // set custom screen size

   player := objectMap.Objects["Player"][0]
   Char := NewCharacter(player, 
         CharacterSpriteFile, 
         CharacterSpriteDirectory, 
         CharacterSpriteStartAnim)
   Char.SetAction("move_up", pikey.Up)
   Char.SetAction("move_left", pikey.Left)
   Char.SetAction("move_right", pikey.Right)
   Char.SetAction("move_down", pikey.Down)
   Char.SetAction("shoot_projectile", pikey.Space)   

   bombs := make([]*Bomb, 0)
   for _, bomb := range objectMap.Objects["Bomb"] {
      bombs = append(bombs, NewBomb(bomb, BombSpriteFile, BombSpriteDirectory, BombSpriteStartAnim))
   }
   
   doorObj := objectMap.Objects["tile_door_1"][0]
   door := NewDoor(doorObj, tileSet.Tiles["tile_door_1"], true)

   world := World{Player: Char, Bombs: bombs, Door: &door,  TileMap: &tileMap}

   pi.Update = func() {
      //WallSystem(&tileMap, &world,  []string{"Tile Layer 1", "wallsides"} )
      BombSystem(&world)


      Char.Update(&world) 
   }

   pi.Draw = func() {      // draw will be executed each frame
      pi.Screen().Clear(32)

      picofont.Print("TEST GAME", 110, 2)

      // for _, layer := range tileMap.Tiles {
      //    // Drawing Tiles 
      //    for y := 0; y < len(layer); y++ {
      //       for x := 0; x < len(layer[y]); x++ {
      //             pi.DrawSprite(tileSet.Tiles[layer[y][x]], x*gameMap.TileWidth, y*gameMap.TileHeight)
      //       }
      //    }
      // }
      DrawTileLayer(&tileMap, "Tile Layer 1")
      DrawTileLayer(&tileMap, "wallsides")

      // Drawing Static Objects 
      // for name, objs := range objectMap.Objects{
      //    for _, obj := range objs {
      //       pi.DrawSprite(tileSet.Tiles[name], obj.Pos.X, obj.Pos.Y)
      //    }
      // }

      world.Door.Draw() 

      for _, bomb := range world.Bombs {
         bomb.Draw()
      }

      Char.Draw() 
   }
   piebiten.Run() // run backend
}
