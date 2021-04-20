package main

import (
	_ "github.com/Oleham/simplevp/db" // To load the database
	"github.com/Oleham/simplevp/gui"
)

func main() {
	gui.StartGUI()
}
