package midihandler

import (
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

type Device *drivers.In

var actions Actions

func UpdateActions(a Actions) {
	actions = a
}

func GetInputDevices() (ports []string) {
	inPorts := midi.GetInPorts()
	ports = make([]string, len(inPorts))
	for i, port := range inPorts {
		ports[i] = port.String()
	}
	return ports
}

func ConnectDevice(deviceName string) (device Device) {
	port, err := midi.FindInPort(deviceName)
	if err != nil {
		println(deviceName)
		panic(err)
	}
	err = port.Open()
	if err != nil {
		panic(err)
	}
	return &port
}

func DisconnectDevice(device Device) {
	err := (*device).Close()
	if err != nil {
		panic(err)
	}
}

func ProcessInput(device Device) (stopCh chan struct{}) {
	stopCh = make(chan struct{})
	go func() {
		println("starting listener...")
		stopFunc, err := midi.ListenTo(*device, runActions)
		if err != nil {
			panic(err)
		}
		<-stopCh
		stopListener(stopFunc, stopCh)
		println("stopped listener")
	}()
	return stopCh
}

func GetInputEvent(device Device) *Event {
	msg := getMIDIMsg(device)
	var ch, value, vel, ctrl uint8
	switch {
	case msg.GetNoteOn(&ch, &value, &vel):
	case msg.GetNoteOff(&ch, &value, &vel):
		return &Event{
			Type:  Note,
			Value: value,
		}
	case msg.GetControlChange(&ch, &ctrl, &value):
		return &Event{
			Type:  CC,
			Value: ctrl,
		}
	case msg.GetProgramChange(&ch, &value):
		return &Event{
			Type:  Prog,
			Value: value,
		}
	}
	return nil
}

func getMIDIMsg(device Device) midi.Message {
	msgCh := make(chan midi.Message)
	stopListen, err := midi.ListenTo(*device, func(msg midi.Message, timestamp int32) {
		msgCh <- msg
	})
	if err != nil {
		panic(err)
	}
	defer stopListen()
	msg := <-msgCh
	close(msgCh)
	return msg
}

func runActions(msg midi.Message, _ int32) {
	println(msg.String())
	// Types: Note On/Off, Control Change Trigger/Release, Program Change
	var ch, value, vel, ctrl uint8
	switch {
	case msg.GetNoteOn(&ch, &value, &vel):
		action, exists := actions.getAction(Event{Type: Note, Value: value})
		if !exists {
			return
		}
		action.executePress()
		break
	case msg.GetNoteOff(&ch, &value, &vel):
		action, exists := actions.getAction(Event{Type: Note, Value: value})
		if !exists {
			return
		}
		action.executeRelease()
		break
	case msg.GetControlChange(&ch, &ctrl, &value):
		action, exists := actions.getAction(Event{Type: CC, Value: ctrl})
		if !exists {
			return
		}
		if value != 0 {
			action.executePress()
		} else {
			action.executeRelease()
		}
		break
	case msg.GetProgramChange(&ch, &value):
		action, exists := actions.getAction(Event{Type: Prog, Value: value})
		if !exists {
			return
		}
		action.executePress()
		break
	}
}

func stopListener(stopFunc func(), stopCh chan struct{}) {
	stopFunc()
	stopCh <- struct{}{}
	close(stopCh)
}
