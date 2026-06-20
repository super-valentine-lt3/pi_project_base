package main
   
import (
   "github.com/elgopher/pi"          // import pi core package
   "github.com/elgopher/pi/picofont" // import very small pico-8 font
   "github.com/elgopher/pi/piebiten" // import backend
    "github.com/lafriks/go-tiled"
    _ "embed"
    "fmt"
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
         spriteName := object.Template.Object.Properties.GetString("sprite")
         position := pi.Position{X: int(object.X), Y: int(object.Y)}
         gameObj := GameObject{object, position}
         posList, ok := objects[spriteName]
         if !ok {
            posList = make([]GameObject, 0)
         }
         posList = append(posList, gameObj)
         objects[spriteName ] = posList 
       }
   }

   return ObjectMap {
      Width: gameMap.Width, 
      Height: gameMap.Height, Objects: objects}
}


type TileMap struct {
    Width  int
    Height int
    //Solid  [][]bool
    Tiles [][]string
}
func NewTileMap() TileMap {
   tiles := make([][]string, gameMap.Height)

   for y := range tiles {
       tiles[y] = make([]string, gameMap.Width)
   }  

   for _, layer := range gameMap.Layers {
       for pos, tile := range layer.Tiles {
           if tile.Nil {
               continue
           }

            tt, err := tile.Tileset.GetTilesetTile(tile.ID)
           if err != nil { continue }
           //if tt.Properties.GetBool("solid") {
               x := pos % gameMap.Width
               y := pos / gameMap.Width

               tiles[y][x] = tt.Type 
           //}
       }
   }

   return TileMap {
      Width: gameMap.Width, 
      Height: gameMap.Height, Tiles: tiles}
}

//go:embed "assets/tiny_dungeon_tilesheet.png"
var spritesPNG []byte

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


func main() {
   pi.Palette = pi.DecodePalette(spritesPNG)
   sprites := pi.DecodeCanvas(spritesPNG)

   // getting the tiles 
   tileList := gameMap.Tilesets[0].Tiles
   for _, tile := range tileList {
      r  := gameMap.Tilesets[0].GetTileRect(tile.ID)
      x, y, width, height := r.Min.X, r.Min.Y, r.Dx(), r.Dy()
      tileSet.Tiles[tile.Type] = pi.SpriteFrom(sprites, x, y, width, height)
   }

   tileMap := NewTileMap()
   objectMap := NewObjectMap()

   fmt.Println(objectMap.Objects)
   pi.SetScreenSize(256, 144) // set custom screen size
   pi.Draw = func() {      // draw will be executed each frame
      picofont.Print("HELLO WORLD", 2, 2)

      // Drawing Tiles 
      for y := 0; y < len(tileMap.Tiles); y++ {
         for x := 0; x < len(tileMap.Tiles[y]); x++ {
               pi.DrawSprite(tileSet.Tiles[tileMap.Tiles[y][x]], x*gameMap.TileWidth, y*gameMap.TileHeight)
         }
      }

      // Drawing Objects 
      for name, objs := range objectMap.Objects{
         for _, obj := range objs {
            pi.DrawSprite(tileSet.Tiles[name], obj.Pos.X, obj.Pos.Y)
         }
      }

   }
   piebiten.Run() // run backend
}
