package main

import (
	"embed"
	"midicontrols/internal/midicontrols/gui"
)

//go:embed all:web/dist
var assets embed.FS

const bundleID = "de.kaes3kuch3n.MIDIControls"

func main() {
	// Create an instance of the app structure
	app := gui.NewApp(bundleID)
	gui.Launch(app, assets)
}
