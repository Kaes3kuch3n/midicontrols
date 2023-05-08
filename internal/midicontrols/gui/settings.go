package gui

import "midicontrols/internal/midicontrols/midihandler"

type Settings struct {
	Actions        midihandler.Actions `json:"actions"`
	SelectedDevice string              `json:"selectedDevice"`
}
