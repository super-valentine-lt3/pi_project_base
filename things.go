package main 
import (
     "github.com/elgopher/pi"
     "slices"
     //"fmt"
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