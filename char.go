package main 

import (
     "github.com/elgopher/pi"
    "github.com/elgopher/pi/pikey"
    "github.com/solarlune/goaseprite"
    "os"
    _ "embed"
    //"fmt"
)

func IsKeyPressedDuration(key pikey.Key, duration int) bool {
    return pikey.Duration(key) > duration
}

 
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

//go:embed "assets/character_try_16x16_indexed.png"
var characterSpritesPNG []byte
const CharacterSpriteFile = "character_try_16x16_indexed.json"
const CharacterSpriteDirectory = "./assets"
const CharacterSpriteStartAnim = "idle_down "

type Character struct {
    Sprite *SpriteAnim 
    Actions map[string]pikey.Key 
    CurrentDirection Direction 
    GameObject GameObject 
   // CanMove bool 
    Points int 
    Health int 
    DamageCooldown int 
    DamageCooldownActive bool 
}
func (c *Character) AddPoints(points int) {
    c.Points += points 
}

func (c *Character) DecreaseHealth(damage int) {
    if c.DamageCooldownActive == false {
        c.Health -= damage 
        c.DamageCooldownActive = true 
    }
}

func (c *Character) SetAction(name string, key pikey.Key) {
    c.Actions[name] = key 
}

func (c *Character) GetArea() pi.IntArea {
    return pi.IntArea{c.GameObject.Pos.X, c.GameObject.Pos.Y, 16, 16}
}

type SpriteAnim struct {
    Sprite    *goaseprite.File
    AsePlayer *goaseprite.Player
    SpriteSheet   map[pi.IntArea]pi.Sprite  
    DefaultSpeed float32
}

func (sprite *SpriteAnim) SetSpeed(speed float32) {
    if speed == -1 {
        sprite.AsePlayer.PlaySpeed = sprite.DefaultSpeed
    } else {
        sprite.AsePlayer.PlaySpeed = float32(speed)
    }
}

func (sprite *SpriteAnim) Play(animation string) {
    sprite.AsePlayer.Play(animation)
}
func (sprite *SpriteAnim) Update(delta float32) {
    sprite.AsePlayer.Update(delta)
}

func NewSpriteAnim(data []byte, 
    file string, directory string, start_anim string, width int, height int, playSpeed float32) *SpriteAnim {
    //sprite, err := goaseprite.Open("character_base_16x16.json", os.DirFS("./assets"))
    
    sprite, err := goaseprite.Open(file, os.DirFS(directory))

    if err != nil {
        panic(err)
    }
    spriteAnim := &SpriteAnim {
        Sprite: sprite, 
    }

    spriteAnim.AsePlayer = spriteAnim.Sprite.CreatePlayer()

    spriteAnim.DefaultSpeed = spriteAnim.AsePlayer.PlaySpeed
    if playSpeed > 0 {
        spriteAnim.AsePlayer.PlaySpeed = float32(playSpeed)
    }

    sprites := pi.DecodeCanvas(data)

    spriteAnim.SpriteSheet = make(map[pi.IntArea]pi.Sprite)
    for _, frame := range sprite.Frames {
        spriteAnim.SpriteSheet[pi.IntArea{frame.X, frame.Y, width, height}] = pi.SpriteFrom(sprites, frame.X, frame.Y, width, height)
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
    spriteAnim := NewSpriteAnim(characterSpritesPNG, sprite_file, sprite_directory, default_anim, 16, 16, -1.0)
    character.Sprite = spriteAnim
    character.GameObject = obj

    character.Actions = make(map[string]pikey.Key)
    character.CurrentDirection = Down 
    character.Health = 100 
    character.DamageCooldown = 25 // frames 
    return character
}

func (c *Character) Draw() {
    c.Sprite.Draw(c.GameObject.Pos.X, c.GameObject.Pos.Y)    
}

func (c *Character) Update(w *World) {//Map *CollisionMap) {
    //if ebiten.IsKeyPressed(ebiten.KeyUp) {
    if IsKeyPressed(c.Actions["move_up"]) {  
       tempY :=  c.GameObject.Pos.Y - speed 
       tempX := c.GameObject.Pos.X 
       //if Map.CanMove(tempX, tempY) {
        if CanMove(w, tempX, tempY) {
            c.Sprite.Play("walk_up")
            c.GameObject.Pos.Y = c.GameObject.Pos.Y - speed 
            c.CurrentDirection = Up 
        }
        //}
    } else if IsKeyPressed(c.Actions["move_down"]) {
        tempY :=  c.GameObject.Pos.Y + speed 
        tempX := c.GameObject.Pos.X 
        //if Map.CanMove(tempX, tempY) {
        if CanMove(w, tempX, tempY)  {
            c.Sprite.Play("walk_down")
            c.GameObject.Pos.Y = c.GameObject.Pos.Y + speed 
            c.CurrentDirection = Down 
        }
        //}
    } else if IsKeyPressed(c.Actions["move_left"]) {
        tempX := c.GameObject.Pos.X - speed 
        tempY := c.GameObject.Pos.Y 
        //if Map.CanMove(tempX, tempY) {
        if CanMove(w, tempX, tempY)  {
            c.Sprite.Play("walk_left")
            c.GameObject.Pos.X = c.GameObject.Pos.X - speed 
            c.CurrentDirection = Left 
        }
        //}
    } else if IsKeyPressed(c.Actions["move_right"]) {
        tempX := c.GameObject.Pos.X + speed 
        tempY := c.GameObject.Pos.Y 
       // if Map.CanMove(tempX, tempY) {
        if CanMove(w, tempX, tempY)  {
            c.Sprite.Play("walk_right")
            c.GameObject.Pos.X = c.GameObject.Pos.X + speed 
            c.CurrentDirection = Right 
        }
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

    if c.DamageCooldownActive {
        c.DamageCooldown -= 1 
        if c.DamageCooldown <= 0 {
            c.DamageCooldownActive = false 
            c.DamageCooldown = 25 
        }
    }
    c.Sprite.Update(float32(1.0 / 60.0))
}
