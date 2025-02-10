package pactl

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
)

// setVolume adjusts volume state for PulseAudio devices.
//
// Parameters:
//   - kind: device type ("sink", "sink-input", "source", "source-input")
//   - nameOrID: name for sinks/sources, numeric ID for inputs
//   - volume: volume level
func setVolume(kind string, nameOrID string, volume string) {
	errPrefix := "ERROR [setVolume(" + kind + ", " + nameOrID + ", " + volume + ")]"
	volumeInPercent := fmt.Sprint(volume) + "%"

	cmd := exec.Command("pactl", "set-"+kind+"-volume", nameOrID, volumeInPercent)
	_, err := cmd.Output()
	if err != nil {
		log.Printf("%s $> pactl set-"+kind+"-volume"+nameOrID+" "+volumeInPercent+": %s\n", errPrefix, err)
	}
}

// setMuted adjusts mute state for PulseAudio devices.
//
// Parameters:
//   - kind: device type ("sink", "sink-input", "source", "source-input")
//   - nameOrID: name for sinks/sources, numeric ID for inputs
//   - muted: muted state
func setMuted(kind string, nameOrID string, muted bool) {
	mutedStr := strconv.FormatBool(muted)
	errPrefix := "ERROR [setMuted(" + kind + ", " + nameOrID + ", " + mutedStr + ")]"

	cmd := exec.Command("pactl", "set-"+kind+"-mute", nameOrID, mutedStr)

	_, err := cmd.Output()
	if err != nil {
		log.Printf("%s $> pactl set-"+kind+"-mute "+nameOrID+" "+mutedStr+": %s\n", errPrefix, err)
	}
}

// moveApp moves input or output app between sink/source devices
//
// Parameters:
//   - kind: device type ("sink-input", "source-output")
//   - appID: sink-input ID or source-output ID
//   - deviceName: sink name or source name
func moveApp(kind string, appID string, deviceName string) {
	errPrefix := "ERROR [moveApp(" + kind + ", " + appID + ", " + deviceName + ")]"

	cmd := exec.Command("pactl", "move-"+kind, appID, deviceName)
	_, err := cmd.Output()
	if err != nil {
		log.Printf("%s $> pactl move-"+kind+" "+appID+" "+deviceName+": %s\n", errPrefix, err)
	}
}
