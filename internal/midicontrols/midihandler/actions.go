package midihandler

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Actions map[EventType]map[uint8]Action

func EmptyActions() Actions {
	return map[EventType]map[uint8]Action{
		Note: make(map[uint8]Action),
		CC:   make(map[uint8]Action),
		Prog: make(map[uint8]Action),
	}
}

func (a Actions) Set(event Event, release bool, command string) {
	act, exists := a.getAction(event)
	if !exists {
		if release {
			act = Action{Release: &command}
		} else {
			act = Action{Press: &command}
		}
		a[event.Type][event.Value] = act
		return
	}

	if release {
		act.Release = &command
	} else {
		act.Press = &command
	}
	a[event.Type][event.Value] = act
}

func (a Actions) Clear(event Event, release bool) {
	act, exists := a.getAction(event)
	if !exists {
		return
	}

	if release {
		act.Release = nil
	} else {
		act.Press = nil
	}
}

func (a Actions) getAction(event Event) (act Action, exists bool) {
	actions, exists := a[event.Type]
	if !exists {
		a[event.Type] = make(map[uint8]Action)
	}
	act, exists = actions[event.Value]
	return act, exists
}

type Action struct {
	Press   *string `json:"press,omitempty"`
	Release *string `json:"release,omitempty"`
}

func (a Action) executePress() {
	if a.Press == nil {
		return
	}
	execute(*a.Press)
}

func (a Action) executeRelease() {
	if a.Release == nil {
		return
	}
	execute(*a.Release)
}

func execute(command string) {
	parts := strings.Fields(command)
	cmd := exec.Command(parts[0], parts[1:]...)

	fmt.Printf("Executing command '%s'\n", cmd.String())

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		println(err)
	}
}
