package main 

import (
  "github.com/elgopher/pi"          // import pi core package
   "github.com/elgopher/pi/picofont" // import very small pico-8 font
   "fmt"
   "github.com/elgopher/pi/pigui"
   "github.com/elgopher/pi/pikey"
      "github.com/elgopher/pi/pievent"

   
)

var pausePanelRoot *pigui.Element
// from Arne32 ):
const (
   ColorBlack = 0
   ColorWhite = 7
)

func init() {
   pausePanelRoot = pigui.New()
   // add a panel (container) at global coordinates
   attachPanel(pausePanelRoot, 100, 48, 63, 32)
   pikey.Target().Subscribe(pikey.Event{pikey.EventDown, pikey.Esc}, func(e pikey.Event, h pievent.Handler){
      Paused = !Paused 
   })

}
var Paused bool 

func PauseSystem()  {
   // if !Paused {
   //    Paused = IsKeyPressed(pikey.Esc) 
   // } else if Paused && IsKeyPressedDuration(pikey.Esc, 2) {
   //    Paused = false 
   // }
}

func UISystem(w *World) {
      picofont.Print(fmt.Sprintf("POINTS: %d  HEALTH: %d BOMBS: %d", w.Player.Points, w.Player.Health, w.Player.BombCount),
                                  70, 2)
}

func attachPanel(parent *pigui.Element, x, y, w, h int) *pigui.Element {
   panel := pigui.Attach(parent, x, y, w, h)
   panel.OnDraw = func(event pigui.DrawEvent) {
      pi.SetColor(ColorWhite)
      pi.Rect(0, 0, panel.W-1, panel.H-1)
      pi.SetColor(ColorBlack)
      pi.RectFill(1, 1, panel.W-2, panel.H-2)
      
      pi.SetColor(ColorWhite)
      picofont.Print("PAUSED" , 22, 8)

   }
   return panel
}

// func attachButton(parent *pigui.Element, x, y, w, h int, label string) *pigui.Element {
//    btn := pigui.Attach(parent, x, y, w, h)
//    btn.OnDraw = func(event pigui.DrawEvent) {
//       var frame, bg, text pi.Color = lightGray, blue, white
//       if event.HasPointer {
//          frame, bg, text = lightGray, lightBlue, white
//       }

//       if event.Pressed {
//          pi.Camera.Y -= 1 // the camera is automatically reset after drawing the element
//          bg = blue
//       }

//       pi.SetColor(frame)
//       pi.Rect(0, 0, w-2, h-2)

//       pi.SetColor(bg)
//       pi.RectFill(1, 1, w-3, h-3)

//       pi.SetColor(text)
//       picofont.Print(label, 6, 4)
//    }
//    return btn
// }
