package main 
import (
     "github.com/elgopher/pi"
     "slices"
     "fmt"
     "strings"
     "strconv" 
     "log"
) 

// Bomb ----


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
    spriteAnim := NewSpriteAnim(bombSpritesPNG, sprite_file, sprite_directory, default_anim, 21, 21, 4)
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
