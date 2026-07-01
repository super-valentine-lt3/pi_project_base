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
    "math"
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
	PickedUp bool 
}

func NewBombInGame(PosX, PosY int) *Bomb {
	obj := GameObject{nil, pi.Position{PosX, PosY}}
	return NewBomb(obj, BombSpriteFile, BombSpriteDirectory, BombSpriteStartAnim)
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
	// //  16, 16 -> 14, 14 to make bomb dropping not hurt player 
    return pi.IntArea{b.GameObject.Pos.X, b.GameObject.Pos.Y, 14, 14} 
}

// Use a bomb system to check and remove bombs if they're touched
// use the callback to do something after the animation completes 

func BombSystem (w *World) {
    // w.Bombs = slices.DeleteFunc(w.Bombs, func(b *Bomb) bool {
    //     return Intersects(w.Player.GetArea(), b.GetArea())
    // })

    w.Bombs = slices.DeleteFunc(w.Bombs, func(b *Bomb) bool {
        return b.Dead || b.PickedUp
    })

    for _, bomb := range w.Bombs {
    	// Bombs used to hurt player, now it adds to inventory 
    	// if Intersects(w.Player.GetArea(), bomb.GetArea()) && !bomb.Detonated {
    	// 	bomb.Detonated = true 
    	// 	bomb.Sprite.Play("explode")
    	// 	w.Player.DecreaseHealth(10)
    	// 	//fmt.Println("i'm here ")
    	// }
    	if Intersects(w.Player.GetArea(), bomb.GetArea()) && !bomb.Detonated {
    		bomb.PickedUp = true 
    		w.Player.PickUpBomb()
    		PlaySound(CollectSample)
    	}

    	for _, crab := range w.Crabs {
    		if Intersects(crab.GetArea(), bomb.GetArea()) && !bomb.Detonated {
	    		bomb.Detonated = true 
	    		bomb.Sprite.Play("explode")
	    		crab.Dead = true 
    		}    	
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
	NextRoom string 
}

func NewDoor(obj GameObject, sprite pi.Sprite) Door {
	locked := obj.Config.Template.Object.Properties.GetBool("locked")
	nextRoom := obj.Config.Properties.GetString("next_room")

	return Door {
		sprite, obj, locked, nextRoom}
}

func (d *Door) Draw() {
	 pi.DrawSprite(d.Sprite, d.GameObject.Pos.X, d.GameObject.Pos.Y)
}

func (d *Door) GetArea() pi.IntArea {
    return pi.IntArea{d.GameObject.Pos.X, d.GameObject.Pos.Y, 16, 16}
}


func DoorSystem(world *World) *string {
	// if IntersectsTouch(world.Door.GetArea(), world.Player.GetArea()) && 
	// 	IsKeyPressed(world.Player.Actions["interact"]) {
	// 	world.Player.TouchingDoor = true 
	// } else {
	// 	world.Player.TouchingDoor = false 
	// }
	if IntersectsTouch(world.Door.GetArea(), world.Player.GetArea()) {
		world.Player.TouchingDoor = true 
		if IsKeyPressed(world.Player.Actions["interact"]) {
			fmt.Println("trying to enter door")
			fmt.Println(world.Door.NextRoom)
			return &world.Door.NextRoom 
		}	
	} else {
		world.Player.TouchingDoor = false 
	}
	return nil 
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
    		PlaySound(CollectSample)
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
	HoverSpeed float64
	HoverHeight float64 
	HoverTimer float64
	BaseY int

	Dead bool 
	Health int 
    DamageFlash bool 
    DamageFlashTimer int 

}
func (b *Bat) DecreaseHealth(value int ) {
	b.Health -= value 
	if b.Health <= 0 {
		b.Dead = true 
	}
	b.DamageFlash = true 
}

func (b *Bat) GetArea() pi.IntArea {
    return pi.IntArea{b.GameObject.Pos.X, b.GameObject.Pos.Y, 16, 16}
}

func NewBat(obj GameObject) *Bat{
    bat := &Bat{}
    spriteAnim := NewSpriteAnim(batSpritesPNG, BatSpriteFile, BatSpriteDirectory, BatSpriteStartAnim, 16, 16, 4.0)
    bat.Sprite = spriteAnim
    bat.GameObject = obj
    bat.HoverSpeed = 2.0 
    bat.HoverHeight = 4.0
    bat.BaseY = obj.Pos.Y 
    bat.Health = 30
    bat.DamageFlashTimer = DamageFlashTime    
    return bat
}
var BatColors []int = []int{9, 10, 17}

func (b *Bat) Draw() {
    if b.DamageFlash {
        var colorToUse int 
        if Between(b.DamageFlashTimer, 1, 3) {
            colorToUse = 2
        } else {
            colorToUse = 7     
        }
        for _, color := range BatColors {
            pi.RemapColor(pi.Color(color), pi.Color(colorToUse))
        }
        b.Sprite.Draw(b.GameObject.Pos.X, b.GameObject.Pos.Y)    
        ResetPalette()
    } else {
   		b.Sprite.Draw(b.GameObject.Pos.X, b.GameObject.Pos.Y)    
	}
}

func (b *Bat) Update(w *World) {
	b.HoverTimer += 1.0 / 30.0 
	b.GameObject.Pos.Y = int(float64(b.BaseY) + math.Sin(b.HoverTimer * b.HoverSpeed) * b.HoverHeight);

	if b.DamageFlash {
        b.DamageFlashTimer -= 1 
        if b.DamageFlashTimer <= 0 {
            b.DamageFlash = false 
            b.DamageFlashTimer = DamageFlashTime
        }
    }

	b.Sprite.Update(float32(1.0 / 60.0))
}

func BatSystem (w *World) {
    for _, bat := range w.Bats {
    	if Intersects(w.Player.GetArea(), bat.GetArea()) {
    		w.Player.DecreaseHealth(10)
    	}
    }
    w.Bats = slices.DeleteFunc(w.Bats, func(b *Bat) bool {
        return b.Dead
    })

    for _, bat := range w.Bats {
       bat.Update(w)
    }
}


// --- crab
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
	Dead bool 
	Health int 
    DamageFlash bool 
    DamageFlashTimer int 

    WalkTimer float64
    WalkTime float64
    WalkSound bool 

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
    crab.Health = 20
    crab.DamageFlashTimer = DamageFlashTime
    crab.WalkTimer = .25

    return crab
}

func (c *Crab) DecreaseHealth(value int ) {
	c.Health -= value 
	if c.Health <= 0 {
		c.Dead = true 
	}
	c.DamageFlash = true 
}
var CrabColors []int = []int{0, 2, 19}

func (c *Crab) Draw() {

    if c.DamageFlash {
        var colorToUse int 
        if Between(c.DamageFlashTimer, 1, 3) {
            colorToUse = 2
        } else {
            colorToUse = 7     
        }
        for _, color := range CrabColors {
            pi.RemapColor(pi.Color(color), pi.Color(colorToUse))
        }
        c.Sprite.Draw(c.GameObject.Pos.X, c.GameObject.Pos.Y)    
        ResetPalette()
    } else {

	    c.Sprite.Draw(c.GameObject.Pos.X, c.GameObject.Pos.Y)    
	}
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
		if !c.WalkSound {
            c.WalkTime = pi.Time 
            c.WalkSound = true 
            PlaySound(CrabWalkSample)
        }

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
	if c.DamageFlash {
        c.DamageFlashTimer -= 1 
        if c.DamageFlashTimer <= 0 {
            c.DamageFlash = false 
            c.DamageFlashTimer = DamageFlashTime
        }
    }
    if (pi.Time - c.WalkTime) >= c.WalkTimer {
        c.WalkSound = false 
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
    w.Crabs = slices.DeleteFunc(w.Crabs, func(c *Crab) bool {
        return c.Dead
    })

    for _, crab := range w.Crabs {
       crab.Update(w)
    }
}

// projectile 
type Projectile struct{
	X int 
	Y int 
	Dir Direction 
	Speed int
	Color pi.Color
	Radius int 
	Dead bool 
}
func NewProjectile(
		BaseX int, 
		BaseY int, 
		Dir Direction, 
		Speed int, color pi.Color) *Projectile{
	proj := &Projectile{
		X: BaseX, 
		Y: BaseY,
		Speed: Speed,  
	}
	proj.Color = color 

	switch Dir {
    case Up:
       proj.Y = proj.Y - 10 
       proj.X = proj.X + 8 

    case Down:
       proj.Y = proj.Y + 10 
       proj.X = proj.X + 8 

    case Left:
       proj.X = proj.X - 10 
       proj.Y = proj.Y + 8 

    case Right:
       proj.X = proj.X + 10   
       proj.Y = proj.Y + 8     
    default:
       //proj.Y = proj.Y + 10 
    }
    proj.Dir = Dir 
    proj.Radius = 2 
    return proj 
}



func (g *Projectile) Update()  { 
	switch g.Dir {
    case Up:
       g.Y = g.Y - g.Speed 
    case Down:
       g.Y = g.Y + g.Speed 
    case Left:
       g.X = g.X - g.Speed 
    case Right:
       g.X = g.X + g.Speed       
    default:
       g.Y = g.Y + g.Speed 
    }

}

func (g *Projectile) Draw() {
	
	before := pi.SetColor(0)
	pi.CircFill(g.X, g.Y, g.Radius )
	pi.SetColor(g.Color)
	pi.CircFill(g.X, g.Y, g.Radius - 1)
	pi.SetColor(before)
}

func (g *Projectile) NextMove() (int, int) {
	var tempX, tempY int = g.X, g.Y

	switch g.Dir {
    case Up:
       tempY = g.Y - g.Speed 
    case Down:
       tempY = g.Y + g.Speed 
    case Left:
       tempX = g.X - g.Speed 
    case Right:
       tempX = g.X + g.Speed       
    default:
    }

    return tempX, tempY 
}

func ProjectileSystem (w *World) {
    w.Projectiles = slices.DeleteFunc(w.Projectiles, func(p *Projectile) bool {
    	tempX, tempY := p.NextMove() 
        return !CanMove(w, tempX, tempY) || p.Dead 
    })


    for _, proj := range w.Projectiles {
       proj.Update()
    	for _, crab := range w.Crabs {
    		if CircleIntersectsRect(crab.GetArea(), proj.X, proj.Y, proj.Radius)  {
	    		crab.DecreaseHealth(5)
	    		PlaySound(CrabHurtSample)
	    		proj.Dead = true 
    		}    	
    	}    
    	for _, bat := range w.Bats {
    		if CircleIntersectsRect(bat.GetArea(), proj.X, proj.Y, proj.Radius)  {
	    		bat.DecreaseHealth(10)
	    		PlaySound(BatHurtSample)
	    		proj.Dead = true 
    		}    	
    	}       	   
    }
}
