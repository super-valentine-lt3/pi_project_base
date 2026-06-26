package main 
import (
     "github.com/elgopher/pi"
     "slices"
     "fmt"
     "strings"
     "strconv" 
     "log"
     _ "embed"
    "math/rand/v2"
) 

// Inclusive 
func Random(Min, Max int) int {
	return Min + rand.IntN((Max+1)-Min)
}

// Bomb ----

//go:embed "assets/bomb_explode.png"
var bombSpritesPNG []byte
const BombSpriteFile = "bomb_explode.json"
const BombSpriteDirectory = "./assets"
const BombSpriteStartAnim = "normal"

type Bomb struct {
	Sprite *SpriteAnim 
	GameObject GameObject 
	Detonated bool  
	Dead bool 
}

func NewBomb(obj GameObject, 
            sprite_file string, 
            sprite_directory string, 
            default_anim string ) *Bomb{
    bomb := &Bomb{}
    spriteAnim := NewSpriteAnim(bombSpritesPNG, sprite_file, sprite_directory, default_anim, 21, 21, 4.0)
    spriteAnim.AsePlayer.OnFrameChange = func() {
    	if spriteAnim.AsePlayer.CurrentTag.Name == "explode" {
    		if spriteAnim.AsePlayer.FrameIndex == 8 {
    			bomb.Dead = true 
    		} 
    	}
    }
    bomb.Sprite = spriteAnim
    bomb.GameObject = obj
    return bomb
}

func (b *Bomb) Draw() {
    b.Sprite.Draw(b.GameObject.Pos.X, b.GameObject.Pos.Y)    
}

func (b *Bomb) Update(w *World) {
	// if Intersects(w.Player.GetArea(), b.GetArea()) {
	// 	b.Sprite.Play("explode")
	// }
	b.Sprite.Update(float32(1.0 / 60.0))
}

func (b *Bomb) GetArea() pi.IntArea {
    return pi.IntArea{b.GameObject.Pos.X, b.GameObject.Pos.Y, 16, 16}
}

// Use a bomb system to check and remove bombs if they're touched
// use the callback to do something after the animation completes 

func BombSystem (w *World) {
    // w.Bombs = slices.DeleteFunc(w.Bombs, func(b *Bomb) bool {
    //     return Intersects(w.Player.GetArea(), b.GetArea())
    // })

    w.Bombs = slices.DeleteFunc(w.Bombs, func(b *Bomb) bool {
        return b.Dead
    })

    for _, bomb := range w.Bombs {
    	if Intersects(w.Player.GetArea(), bomb.GetArea()) && !bomb.Detonated {
    		bomb.Detonated = true 
    		bomb.Sprite.Play("explode")
    		w.Player.DecreaseHealth(10)
    		//fmt.Println("i'm here ")
    	}
    }

    for _, bomb := range w.Bombs {
       bomb.Update(w)
    }
 
}

// End Bomb --- 

// Door --

type Door struct {
	Sprite pi.Sprite 
	GameObject GameObject 
	Locked bool 
}

func NewDoor(obj GameObject, sprite pi.Sprite, locked bool) Door {
	return Door {
		sprite, obj, locked}
}

func (d *Door) Draw() {
	 pi.DrawSprite(d.Sprite, d.GameObject.Pos.X, d.GameObject.Pos.Y)
}

func (d *Door) GetArea() pi.IntArea {
    return pi.IntArea{d.GameObject.Pos.X, d.GameObject.Pos.Y, 16, 16}
}


func DoorSystem(world *World) {
	if IntersectsTouch(world.Door.GetArea(), world.Player.GetArea()) && 
		IsKeyPressed(world.Player.Actions["interact"]) {
		fmt.Println("Trying to enter")
	}
}

// Gem ---
/**
 *  Gems give points 
 *  Points shown in UI
 *  Points on Character
 *  Points increased after collecting 
 * */
type Gem struct {
	Sprite pi.Sprite 
	GameObject GameObject 
	paletteMap map[int]int 
	PointValue int 
}

func NewGem(obj GameObject, sprite pi.Sprite) *Gem {
	
	var paletteMap map[int]int = nil 
	from := obj.Config.Template.Object.Properties.GetString("paletteFrom")
	to := obj.Config.Properties.GetString("paletteTo")
	if to != "" {
		paletteMap = make(map[int]int)
		fromList := strings.Split(from, ",")
		toList := strings.Split(to, ",")
		for i := 0; i < len(fromList); i++ {
			fromVal, err := strconv.Atoi(fromList[i])
			if err != nil {
				log.Fatalf("Conversion failed: %v", err)
			}
			toVal, err := strconv.Atoi(toList[i])
			if err != nil {
				log.Fatalf("Conversion failed: %v", err)
			}
			paletteMap[fromVal] = toVal
		}
	}
	var usedPoints int 
	customPoints := obj.Config.Properties.GetInt("points")
	if customPoints == 0 {
		usedPoints = obj.Config.Template.Object.Properties.GetInt("points")
	} else {
		usedPoints = customPoints
	}
	return &Gem {
		sprite, obj, paletteMap, usedPoints}
}

func ResetPalette() {
   pi.ResetColorTables()
   pi.SetTransparency(0, false)
   pi.SetTransparency(32, true)
}

func (d *Gem) Draw() {
	if d.paletteMap == nil {
	 	pi.DrawSprite(d.Sprite, d.GameObject.Pos.X, d.GameObject.Pos.Y)
	} else {
		for fromColor, toColor := range d.paletteMap {
			pi.RemapColor(pi.Color(fromColor), pi.Color(toColor))
		}
		pi.DrawSprite(d.Sprite, d.GameObject.Pos.X, d.GameObject.Pos.Y)
		ResetPalette()
	}
}

func (d *Gem) GetArea() pi.IntArea {
    return pi.IntArea{d.GameObject.Pos.X, d.GameObject.Pos.Y, 16, 16}
}

func GemSystem (w *World) {
    w.Gems = slices.DeleteFunc(w.Gems, func(g *Gem) bool {
    	if Intersects(w.Player.GetArea(), g.GetArea()) {
    		w.Player.AddPoints(g.PointValue)
    		return true 
    	}
        return false 
    }) 
}

// --- bat

//go:embed "assets/mybat.png"
var batSpritesPNG []byte
const BatSpriteFile = "mybat.json"
const BatSpriteDirectory = "./assets"
const BatSpriteStartAnim = "normal"

type Bat struct {
	Sprite *SpriteAnim 
	GameObject GameObject 
}

func NewBat(obj GameObject) *Bat{
    bat := &Bat{}
    spriteAnim := NewSpriteAnim(batSpritesPNG, BatSpriteFile, BatSpriteDirectory, BatSpriteStartAnim, 16, 16, 4.0)
    bat.Sprite = spriteAnim
    bat.GameObject = obj
    return bat
}

func (b *Bat) Draw() {
    b.Sprite.Draw(b.GameObject.Pos.X, b.GameObject.Pos.Y)    
}

func (b *Bat) Update(w *World) {
	b.Sprite.Update(float32(1.0 / 60.0))
}

// --- spider
//go:embed "assets/crab_try.png"
var crabSpritesPNG []byte
const CrabSpriteFile = "crab_try.json"
const CrabSpriteDirectory = "./assets"
const CrabSpriteStartAnim = "idle"

type Crab struct {
	Sprite *SpriteAnim 
	GameObject GameObject 
	Dir Direction 
	Dir2 Direction 
	Timer float64 
	InitTime float64 
	MoveAnimSpeed float32
	IdleAnimSpeed float32 
	Idle bool 
}

func NewCrab(obj GameObject) *Crab{
    crab := &Crab{}
    spriteAnim := NewSpriteAnim(crabSpritesPNG, CrabSpriteFile,
    					       CrabSpriteDirectory, CrabSpriteStartAnim, 16, 16, -1.0)
    crab.Sprite = spriteAnim
    crab.GameObject = obj
    crab.InitTime = pi.Time 
    crab.Timer = 4.0 
    crab.MoveAnimSpeed = .75 
    crab.IdleAnimSpeed = -1.0 
    crab.Idle = true 
    return crab
}

func (c *Crab) Draw() {
    c.Sprite.Draw(c.GameObject.Pos.X, c.GameObject.Pos.Y)    
}

func (c *Crab) Move(w *World, dir Direction) {
	 if dir == Up {  
       tempY :=  c.GameObject.Pos.Y - speed 
       tempX := c.GameObject.Pos.X 
        if CanMove(w, tempX, tempY) {
            c.GameObject.Pos.Y = c.GameObject.Pos.Y - speed 
        }
    } else if dir == Down {
        tempY :=  c.GameObject.Pos.Y + speed 
        tempX := c.GameObject.Pos.X 
        if CanMove(w, tempX, tempY)  {
            c.GameObject.Pos.Y = c.GameObject.Pos.Y + speed 
        }
    } else if dir == Left {
        tempX := c.GameObject.Pos.X - speed 
        tempY := c.GameObject.Pos.Y 
        if CanMove(w, tempX, tempY)  {
            c.GameObject.Pos.X = c.GameObject.Pos.X - speed 
        }
    } else if dir == Right {
        tempX := c.GameObject.Pos.X + speed 
        tempY := c.GameObject.Pos.Y 
        if CanMove(w, tempX, tempY)  {
            c.GameObject.Pos.X = c.GameObject.Pos.X + speed 
        }
    } else if dir == UpRight {
       tempY :=  c.GameObject.Pos.Y - speed 
       tempX := c.GameObject.Pos.X + speed 
	   if CanMove(w, tempX, tempY)  {
	   		c.GameObject.Pos.X = tempX 
			c.GameObject.Pos.Y = tempY 
	   }
    } else if dir == UpLeft {
       tempY :=  c.GameObject.Pos.Y - speed 
       tempX := c.GameObject.Pos.X - speed 
	   if CanMove(w, tempX, tempY)  {
	   		c.GameObject.Pos.X = tempX 
			c.GameObject.Pos.Y = tempY 
	   }
    } else if dir == DownRight {
       tempY :=  c.GameObject.Pos.Y + speed 
       tempX := c.GameObject.Pos.X + speed 
	   if CanMove(w, tempX, tempY)  {
	   		c.GameObject.Pos.X = tempX 
			c.GameObject.Pos.Y = tempY 
	   }
    } else if dir == DownLeft {
       tempY :=  c.GameObject.Pos.Y + speed 
       tempX := c.GameObject.Pos.X - speed 
	   if CanMove(w, tempX, tempY)  {
	   		c.GameObject.Pos.X = tempX 
			c.GameObject.Pos.Y = tempY 
	   }
    }
}

func (c *Crab) Update(w *World) {
	if pi.Time - c.InitTime >= c.Timer {
		if c.Idle {
			c.Sprite.Play("move")
			c.Sprite.SetSpeed(c.MoveAnimSpeed)
			c.Dir = Direction(Random(0, 7))
			c.InitTime = pi.Time 
			c.Idle = false 
		} else {
			c.Sprite.Play("idle")
			c.Sprite.SetSpeed(c.IdleAnimSpeed)
			c.InitTime = pi.Time 
			c.Idle = true 
		}
	} else if !c.Idle {
		c.Move(w, c.Dir)
		//  if c.Dir == Up {  
	    //    tempY :=  c.GameObject.Pos.Y - speed 
	    //    tempX := c.GameObject.Pos.X 
	    //     if CanMove(w, tempX, tempY) {
	    //         c.GameObject.Pos.Y = c.GameObject.Pos.Y - speed 
	    //     }
	    // } else if c.Dir == Down {
	    //     tempY :=  c.GameObject.Pos.Y + speed 
	    //     tempX := c.GameObject.Pos.X 
	    //     if CanMove(w, tempX, tempY)  {
	    //         c.GameObject.Pos.Y = c.GameObject.Pos.Y + speed 
	    //     }
	    // } else if c.Dir == Left {
	    //     tempX := c.GameObject.Pos.X - speed 
	    //     tempY := c.GameObject.Pos.Y 
	    //     if CanMove(w, tempX, tempY)  {
	    //         c.GameObject.Pos.X = c.GameObject.Pos.X - speed 
	    //     }
	    // } else if c.Dir == Right {
	    //     tempX := c.GameObject.Pos.X + speed 
	    //     tempY := c.GameObject.Pos.Y 
	    //     if CanMove(w, tempX, tempY)  {
	    //         c.GameObject.Pos.X = c.GameObject.Pos.X + speed 
	    //     }
	    // } 
	}
	c.Sprite.Update(float32(1.0 / 30.0))
}

func (c *Crab) GetArea() pi.IntArea {
    return pi.IntArea{c.GameObject.Pos.X, c.GameObject.Pos.Y, 16, 16}
}


func CrabSystem (w *World) {
    for _, crab := range w.Crabs {
    	if Intersects(w.Player.GetArea(), crab.GetArea()) {
    		w.Player.DecreaseHealth(15)
    	}
    }

    for _, crab := range w.Crabs {
       crab.Update(w)
    }
}


