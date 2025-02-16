package pactl

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/undg/go-prapi/logger"
)

// setVolume adjusts volume state for PulseAudio devices.
//
// Parameters:
//   - kind: device type ("sink", "sink-input", "source", "source-input")
//   - nameOrID: name for sinks/sources, numeric ID for inputs
//   - volume: volume level
func setVolume(kind string, nameOrID string, volume string) {
	volumeInPercent := fmt.Sprint(volume) + "%"

	cmd := exec.Command("pactl", "set-"+kind+"-volume", nameOrID, volumeInPercent)

	logger.Debug().Msgf("$> pactl set-"+kind+"-volume %s %s", nameOrID, volumeInPercent)
	logger.Info().Str("kind", kind).Str("nameOrID", nameOrID).Str("volumeInPercent", volumeInPercent).Msg("exec.Command(pactl ***) in setVolume()")

	_, err := cmd.Output()
	if err != nil {
		logger.Error().Err(err).Msg("exec.Command(pactl ***) FAIL in setVolume()")
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

	cmd := exec.Command("pactl", "set-"+kind+"-mute", nameOrID, mutedStr)

	logger.Debug().Msgf("$> pactl set-"+kind+"-mute %s %s", nameOrID, mutedStr)
	logger.Info().Str("kind", kind).Str("nameOrID", nameOrID).Str("mutedStr", mutedStr).Msg("exec.Command(pactl ***) in setMuted()")

	_, err := cmd.Output()
	if err != nil {
		logger.Error().Err(err).Msg("exec.Command(pactl ***) FAIL in setMuted()")
	}
}

// moveApp moves input or output app between sink/source devices
//
// Parameters:
//   - kind: device type ("sink-input", "source-output")
//   - appID: sink-input ID or source-output ID
//   - deviceName: sink name or source name
func moveApp(kind string, appID string, deviceName string) {
	cmd := exec.Command("pactl", "move-"+kind, appID, deviceName)

	logger.Debug().Msgf("$> pactl move-"+kind+"-mute %s %s", appID, deviceName)
	logger.Info().Str("kind", kind).Str("appID", appID).Str("deviceName", deviceName).Msg("exec.Command(pactl ***) in moveApp()")

	_, err := cmd.Output()
	if err != nil {
		logger.Error().Err(err).Msg("exec.Command(pactl ***) FAIL in moveApp()")
	}
}
