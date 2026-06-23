package main 

import (
 //  "github.com/elgopher/pi"          // import pi core package
   "github.com/elgopher/pi/picofont" // import very small pico-8 font
   "fmt"
)

func UISystem(w *World) {

      picofont.Print(fmt.Sprintf("POINTS: %d", w.Player.Points), 110, 2)

}