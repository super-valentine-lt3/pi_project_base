package main 

import (
   "github.com/elgopher/pi"
   "github.com/elgopher/pi/picofont" 
   "github.com/elgopher/pi/pikey"
   "github.com/elgopher/pi/pievent"
  // "fmt"
)

type Screen interface {
	Init()
	Draw()
	Update() 
}

var LogoX, LogoY int = 100, 48
var TitleFlashTime int = 20

type TitleScreen struct {
	Logo pi.Sprite 
	FlashTimer int
	ShowMessage bool 
	EnterEventHandle pievent.Handler
}

func (ts *TitleScreen) StartGame(e pikey.Event, h pievent.Handler) {
	pikey.Target().Unsubscribe(ts.EnterEventHandle)
}

func (ts *TitleScreen) Init() {
	pi.Palette = pi.DecodePalette(spritesPNG)
   	MainSprites = pi.DecodeCanvas(spritesPNG)
   	ts.Logo = pi.SpriteFrom(MainSprites, 192, 0, 64, 32)
   	ts.ShowMessage = true 
    ts.EnterEventHandle = pikey.Target().Subscribe(pikey.Event{pikey.EventDown, pikey.Enter}, ts.StartGame)

}

func (ts *TitleScreen) Update() {
	ts.FlashTimer += 1 
	if ts.FlashTimer >= TitleFlashTime && ts.ShowMessage{
		ts.ShowMessage = false 
		ts.FlashTimer = 0
	} 
	if ts.FlashTimer >= TitleFlashTime && !ts.ShowMessage {
		ts.ShowMessage = true 
		ts.FlashTimer = 0 
	}
}

func (ts *TitleScreen) Draw() {
	pi.Screen().Clear(32)
	pi.DrawSprite(ts.Logo, LogoX, LogoY)
	if ts.ShowMessage {
		picofont.Print("PRESS ENTER TO START", 95, 84)
	} 
}