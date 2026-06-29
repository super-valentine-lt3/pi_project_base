package main

import (
	 "github.com/lafriks/go-tiled"
	 "github.com/elgopher/pi"

)

func NewTileSet(gameMap *tiled.Map) TileSet {
   tileSet = TileSet{}
   tileSet.Tiles = make(map[string]pi.Sprite)

   tileList := gameMap.Tilesets[0].Tiles
   for _, tile := range tileList {
       r  := gameMap.Tilesets[0].GetTileRect(tile.ID)
       x, y, width, height := r.Min.X, r.Min.Y, r.Dx(), r.Dy()
       tileSet.Tiles[tile.Type] = pi.SpriteFrom(MainSprites, x, y, width, height)
   }

   return tileSet
}
