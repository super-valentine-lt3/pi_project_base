package main 

type Bomb struct {
	Sprite *SpriteAnim 
	GameObject GameObject 
}

func NewBomb(obj GameObject, 
            sprite_file string, 
            sprite_directory string, 
            default_anim string ) *Bomb{
    bomb := &Bomb{}
    spriteAnim := NewSpriteAnim(bombSpritesPNG, sprite_file, sprite_directory, default_anim)
    bomb.Sprite = spriteAnim
    bomb.GameObject = obj
    return bomb
}

func (b *Bomb) Draw() {
    b.Sprite.Draw(b.GameObject.Pos.X, b.GameObject.Pos.Y)    
}

