package main 

import (
     "github.com/elgopher/pi"
    "github.com/elgopher/pi/pikey"
    "github.com/solarlune/goaseprite"
    "os"
    "fmt"
)

func IsKeyPressed(key pikey.Key) bool {
    return pikey.Duration(key) > 0
}
var speed = 2

type Direction int

const (
    Up Direction = iota
    Down 
    Left 
    Right 
)

var CurrentDirection = Down 

type Character struct {
    Sprite *SpriteAnim 
    Actions map[string]pikey.Key 
    CurrentDirection Direction 
    GameObject GameObject 
}

func (c *Character) SetAction(name string, key pikey.Key) {
    c.Actions[name] = key 
}


type SpriteAnim struct {
    Sprite    *goaseprite.File
    AsePlayer *goaseprite.Player
    SpriteSheet   map[pi.IntArea]pi.Sprite  
}

func (sprite *SpriteAnim) Play(animation string) {
    fmt.Println(animation)
    sprite.AsePlayer.Play(animation)
}
func (sprite *SpriteAnim) Update(delta float32) {
    sprite.AsePlayer.Update(delta)
}

func NewSpriteAnim(
    file string, directory string, start_anim string) *SpriteAnim {
    //sprite, err := goaseprite.Open("character_base_16x16.json", os.DirFS("./assets"))
    
    sprite, err := goaseprite.Open(file, os.DirFS(directory))

    if err != nil {
        panic(err)
    }
    spriteAnim := &SpriteAnim {
        Sprite: sprite, 
    }

    spriteAnim.AsePlayer = spriteAnim.Sprite.CreatePlayer()

    sprites := pi.DecodeCanvas(characterSpritesPNG)

    spriteAnim.SpriteSheet = make(map[pi.IntArea]pi.Sprite)
    for _, frame := range sprite.Frames {
        spriteAnim.SpriteSheet[pi.IntArea{frame.X, frame.Y, 16, 16}] = pi.SpriteFrom(sprites, frame.X, frame.Y, 16, 16)
    }

    spriteAnim.AsePlayer.Play(start_anim)
    return spriteAnim
}

func (sprite *SpriteAnim) Draw(PosX int,  PosY int) {
    x1, y1, x2, y2 :=  sprite.AsePlayer.CurrentFrameCoords()
    area := pi.IntArea{x1, y1, x2-x1, y2-y1}
    pi.DrawSprite(sprite.SpriteSheet[area], PosX, PosY)
}

func NewCharacter(obj GameObject, 
            sprite_file string, 
            sprite_directory string, 
            default_anim string ) *Character{
    character := &Character{}
    spriteAnim := NewSpriteAnim(sprite_file, sprite_directory, default_anim)
    character.Sprite = spriteAnim
    character.GameObject = obj

    character.Actions = make(map[string]pikey.Key)
    character.CurrentDirection = Down 

    return character
}

func (c *Character) Draw() {
    c.Sprite.Draw(c.GameObject.Pos.X, c.GameObject.Pos.Y)    
}

func (c *Character) Update() {//Map *CollisionMap) {
    //if ebiten.IsKeyPressed(ebiten.KeyUp) {
    if IsKeyPressed(c.Actions["move_up"]) {  
       // tempY :=  c.PosY - speed 
       // tempX := c.PosX 
       //if Map.CanMove(tempX, tempY) {
            c.Sprite.Play("walk_up")
            c.GameObject.Pos.Y = c.GameObject.Pos.Y - speed 
            c.CurrentDirection = Up 
        //}
    } else if IsKeyPressed(c.Actions["move_down"]) {
       // tempY :=  c.PosY + speed 
       // tempX := c.PosX 
        //if Map.CanMove(tempX, tempY) {
            c.Sprite.Play("walk_down")
            c.GameObject.Pos.Y = c.GameObject.Pos.Y + speed 
            c.CurrentDirection = Down 
        //}
    } else if IsKeyPressed(c.Actions["move_left"]) {
        // tempX := c.PosX - speed 
        // tempY := c.PosY 
        //if Map.CanMove(tempX, tempY) {
            c.Sprite.Play("walk_left")
            c.GameObject.Pos.X = c.GameObject.Pos.X - speed 
            c.CurrentDirection = Left 
        //}
    } else if IsKeyPressed(c.Actions["move_right"]) {
        //tempX := c.PosX + speed 
        //tempY := c.PosY 
       // if Map.CanMove(tempX, tempY) {
            c.Sprite.Play("walk_right")
            c.GameObject.Pos.X = c.GameObject.Pos.X + speed 
            c.CurrentDirection = Right 
        //}
    } else {
        switch c.CurrentDirection {
        case Up:
           c.Sprite.Play("idle_up")
        case Down:
           c.Sprite.Play("idle_down")
        case Left:
           c.Sprite.Play("idle_left")
        case Right:
           c.Sprite.Play("idle_right")         
        default:
           c.Sprite.Play("idle_down")
        }
    }

    c.Sprite.Update(float32(1.0 / 60.0))
}
