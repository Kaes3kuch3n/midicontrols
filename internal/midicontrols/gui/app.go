package gui

import (
	"context"
	"encoding/json"
	"fmt"
	"midicontrols/internal/midicontrols/midihandler"
	"os"
	"path"
)

const settingsFileName = "settings.json"

// App struct
type App struct {
	ctx                context.Context
	bundleID           string
	appDir             string
	settings           Settings
	selectedMIDIDevice midihandler.Device
	listenerStopCh     chan struct{}
}

// NewApp creates a new app application struct
func NewApp(bundleID string) *App {
	return &App{bundleID: bundleID}
}

// Startup is called when the app starts. The context is saved,
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.createAppDirs()
}

func (a *App) Shutdown(_ context.Context) {
	settings, err := json.Marshal(a.settings)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(a.getAppFilePath(settingsFileName), settings, 0644)
	if err != nil {
		panic(err)
	}
}

func (a *App) LoadSettings() Settings {
	file, err := os.ReadFile(a.getAppFilePath(settingsFileName))
	if err != nil {
		a.settings = loadDefaultSettings()
		midihandler.UpdateActions(a.settings.Actions)
		return a.settings
	}

	err = json.Unmarshal(file, &a.settings)
	if err != nil {
		panic(err)
	}

	midihandler.UpdateActions(a.settings.Actions)
	return a.settings
}

func (a *App) GetMIDIDevices() []string {
	return midihandler.GetInputDevices()
}

func (a *App) SelectDevice(deviceName string) string {
	if a.selectedMIDIDevice != nil {
		midihandler.DisconnectDevice(a.selectedMIDIDevice)
	}
	a.selectedMIDIDevice = midihandler.ConnectDevice(deviceName)
	a.listenerStopCh = midihandler.ProcessInput(a.selectedMIDIDevice)

	// Store selected device in settings
	a.settings.SelectedDevice = deviceName

	return deviceName
}

func (a *App) ListenForInput() midihandler.Event {
	a.listenerStopCh <- struct{}{}
	<-a.listenerStopCh
	println("listening for input...")
	var event *midihandler.Event
	for event == nil {
		event = midihandler.GetInputEvent(a.selectedMIDIDevice)
	}
	println("received input: ", event.Type, event.Value)
	a.listenerStopCh = midihandler.ProcessInput(a.selectedMIDIDevice)
	return *event
}

func (a *App) SetCommand(event midihandler.Event, release bool, command string) {
	a.settings.Actions.Set(event, release, command)
}

func (a *App) ClearCommand(event midihandler.Event, release bool) {
	a.settings.Actions.Clear(event, release)
}

func (a *App) GetActions() midihandler.Actions {
	return a.settings.Actions
}

func (a *App) getAppFilePath(filePath string) string {
	return path.Join(a.appDir, filePath)
}

func (a *App) createAppDirs() {
	dir, err := os.UserConfigDir()
	if err != nil {
		panic(fmt.Errorf("unable to access app directory [%w]", err))
	}
	// Create app directory
	a.appDir = path.Join(dir, a.bundleID)
	err = os.MkdirAll(a.appDir, 0755)
	if err != nil {
		panic(fmt.Errorf("unable to create app directory [%w]", err))
	}
}

func loadDefaultSettings() Settings {
	settings := Settings{Actions: midihandler.EmptyActions()}
	settings.Actions.Set(midihandler.Event{Type: midihandler.Note, Value: 36}, false, "open /Users/luishankel/Downloads/")
	return settings
}
