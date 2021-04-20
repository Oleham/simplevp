package main

import (
	"github.com/Oleham/simplevp/db"
	_ "github.com/Oleham/simplevp/db" // To load the database
	"github.com/Oleham/simplevp/gui"
)

func main() {
	db.UpdateJobs() // Temp
	gui.StartGUI()
}
