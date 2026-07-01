package main 

import (
     "github.com/elgopher/pi"
    "github.com/elgopher/pi/pikey"
    "github.com/solarlune/goaseprite"
    //"os"
    _ "embed"
    // "fmt"
    "github.com/elgopher/pi/pievent"

)
func PlayerCollisionArea(x, y int) pi.IntArea {
    return pi.IntArea{
        X: x + 4,
        Y: y + 10,
        W: 8,
        H: 6,
    }
}
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
    
    UpRight 
    DownRight
    UpLeft 
    DownLeft 
)


var CurrentDirection = Down 

//go:embed "assets/character_try_16x16_indexed.png"
var characterSpritesPNG []byte
const CharacterSpriteFile = "character_try_16x16_indexed.json"
const CharacterSpriteDirectory = "./assets"
const CharacterSpriteStartAnim = "idle_down "

var CharColors []int = []int{0, 2, 3, 31, 10, 21}
var DamageFlashTime = 6

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
    DamageFlash bool 
    DamageFlashTimer int 
    BombCount int 
    DroppingBomb bool 
    ShootingProjectile bool 
    TouchingDoor bool 
}

func (c *Character) PickUpBomb() {
    c.BombCount += 1 
}

func (c *Character) AddPoints(points int) {
    c.Points += points 
}

func (c *Character) DecreaseHealth(damage int) {
    if c.DamageCooldownActive == false {
        c.Health -= damage 
        c.DamageCooldownActive = true 
        c.DamageFlash = true 
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
    
   // sprite, err := goaseprite.Open(file, os.DirFS(directory))
    sprite, err := goaseprite.Open("assets/" + file, assetsFS)

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
    character.DamageFlashTimer = DamageFlashTime
    character.BombCount = 3 

   // pikey.Target().Subscribe(pikey.Event{pikey.EventDown, pikey.Space}, func(e pikey.Event, h pievent.Handler){
   //      character.DroppingBomb = true 
   // })
   pikey.Target().SubscribeAll(func(e pikey.Event, h pievent.Handler){
        if e.Type == pikey.EventDown  {
            if e.Key == pikey.CtrlLeft || e.Key == pikey.CtrlRight {
                character.DroppingBomb = true 
            }
            if e.Key == pikey.Space && !character.TouchingDoor{
                character.ShootingProjectile = true 
            }
        }
   })

    return character
}

func Between(Value, Start, End int) bool {
    return Value >= Start && Value <= End 
}

func (c *Character) Draw() {
    if c.DamageFlash {
        var colorToUse int 
        if Between(c.DamageFlashTimer, 1, 3) {
            colorToUse = 2
        } else {
            colorToUse = 7     
        }
        for _, color := range CharColors {
            pi.RemapColor(pi.Color(color), pi.Color(colorToUse))
        }
        c.Sprite.Draw(c.GameObject.Pos.X, c.GameObject.Pos.Y)    
        ResetPalette()
    } else {
        c.Sprite.Draw(c.GameObject.Pos.X, c.GameObject.Pos.Y)    
    } 
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

    if c.DroppingBomb && c.BombCount > 0 { //IsKeyPressed(c.Actions["drop_bomb"]) && c.BombCount > 0 {
        var bombX, bombY int  = c.GameObject.Pos.X, c.GameObject.Pos.Y 
        switch c.CurrentDirection {
        case Up: 
            bombY -= 16 
        case Down: 
            bombY += 16
        case Left: 
            bombX -= 16
        case Right: 
            bombX += 16
        }
        bomb := NewBombInGame(bombX, bombY) 
        w.Bombs = append(w.Bombs, bomb)
        c.BombCount -= 1 
        c.DroppingBomb = false 
    }
    
    if c.ShootingProjectile {
        w.Projectiles = append(w.Projectiles, NewProjectile(c.GameObject.Pos.X, c.GameObject.Pos.Y, c.CurrentDirection, 6, pi.Color(13)))
        c.ShootingProjectile = false 
    }

    if c.DamageCooldownActive {
        c.DamageCooldown -= 1 
        if c.DamageCooldown <= 0 {
            c.DamageCooldownActive = false 
            c.DamageCooldown = 25 
        }
    }
    if c.DamageFlash {
        c.DamageFlashTimer -= 1 
        if c.DamageFlashTimer <= 0 {
            c.DamageFlash = false 
            c.DamageFlashTimer = DamageFlashTime
        }
    }
    c.Sprite.Update(float32(1.0 / 60.0))
}
