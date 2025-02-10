package pactl

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/undg/go-prapi/buildinfo"
)

type Status = struct {
	Outputs   []Output            `json:"outputs" doc:"List of output devices"`
	Apps      []App               `json:"apps" doc:"List of applications"`
	Sources   []Source            `json:"sources" doc:"List of microphones and other sources"`
	BuildInfo buildinfo.BuildInfo `json:"buildInfo" doc:"Build information"`
}

type Output struct {
	ID     int    `json:"id" doc:"The id of the sink. Same  as name"`
	Name   string `json:"name" doc:"The name of the sink. Same as id"`
	Label  string `json:"label" doc:"Human-readable label for the sink"`
	Volume int    `json:"volume" doc:"Current volume level of the sink"`
	Muted  bool   `json:"muted" doc:"Whether the sink is muted"`
}

type Source struct {
	ID     int    `json:"id" doc:"The id of the source. Same  as name"`
	Name   string `json:"name" doc:"The name of the source. Same as id"`
	Label  string `json:"label" doc:"Human-readable label for the source"`
	Volume int    `json:"volume" doc:"Current volume level of the source"`
	Muted  bool   `json:"muted" doc:"Whether the source is muted"`
}

type App struct {
	ID       int    `json:"id" doc:"The id of the sink. Same  as name"`
	OutputID int    `json:"outputId" doc:"Id of parrent device, same as output.id"`
	Label    string `json:"label" doc:"Human-readable label for the sink"`
	Volume   int    `json:"volume" doc:"Current volume level of the sink"`
	Muted    bool   `json:"muted" doc:"Whether the sink is muted"`
}

func SetSinkVolume(sinkName string, volume string) {
	errPrefix := "ERROR [SetSinkVolume()]"
	volumeInPercent := fmt.Sprint(volume) + "%"

	cmd := exec.Command("pactl", "set-sink-volume", sinkName, volumeInPercent)
	_, err := cmd.Output()
	if err != nil {
		log.Printf("%s pactl set-sink-volume: %s\n", errPrefix, err)
		log.Printf("%s pactl set-sink-volume: {SINK_NAME: %s ; VOLUME: %s}\n", errPrefix, sinkName, volumeInPercent)
	}
}

func SetSinkMuted(sinkName string, muted bool) {
	errPrefix := "ERROR [SetSinkMuted()]"

	mutedCmd := "false"
	if muted {
		mutedCmd = "true"
	}

	cmd := exec.Command("pactl", "set-sink-mute", sinkName, mutedCmd)
	_, err := cmd.Output()
	if err != nil {
		log.Printf("%s pactl set-sink-mute: %s\n", errPrefix, err)
		log.Printf("%s pactl set-sink-mute: {SINK_NAME: %s ; MUTED: %s}\n", errPrefix, sinkName, mutedCmd)
	}
}

func SetSinkInputVolume(sinkInputID string, volume string) {
	errPrefix := "ERROR [SetSinkInputVolume()]"
	volumeInPercent := volume + "%"

	cmd := exec.Command("pactl", "set-sink-input-volume", sinkInputID, volumeInPercent)
	_, err := cmd.Output()
	if err != nil {
		log.Printf("%s pactl set-sink-input-volume: %s\n", errPrefix, err)
		log.Printf("%s pactl set-sink-input-volume: {SINK_INPUT_ID: %s ; VOLUME: %s}\n", errPrefix, sinkInputID, volumeInPercent)
	}
}

func SetSinkInputMuted(sinkInputID string, muted bool) {
	errPrefix := "ERROR [SetSinkInputMuted()]"

	mutedCmd := "false"
	if muted {
		mutedCmd = "true"
	}

	cmd := exec.Command("pactl", "set-sink-input-mute", sinkInputID, mutedCmd)
	_, err := cmd.Output()
	if err != nil {
		log.Printf("%s pactl set-sink-mute: %s\n", errPrefix, err)
		log.Printf("%s pactl set-sink-mute: {SINK_INPUT_ID: %s ; MUTED: %s}\n", errPrefix, sinkInputID, mutedCmd)
	}
}

func MoveSinkInput(sinkInputID string, sinkName string) {
	errPrefix := "Error [MoveSinkInput()]"

	cmd := exec.Command("pactl", "move-sink-input", sinkInputID, sinkName)
	_, err := cmd.Output()
	if err != nil {
		log.Printf("%s pactl move-sink-input: %s\n", errPrefix, err)
		log.Printf("%s pactl move-sink-input: {SINK_INPUT_ID: %s ; SINK_NAME: %s}\n", errPrefix, sinkInputID, sinkName)
	}
}

func SetSourceVolume(sourceName string, volume string) {
	errPrefix := "ERROR [SetSourceVolume()]"
	volumeInPercent := fmt.Sprint(volume) + "%"

	cmd := exec.Command("pactl", "set-s-volume", sourceName, volumeInPercent)
	_, err := cmd.Output()
	if err != nil {
		log.Printf("%s pactl set-source-volume: %s\n", errPrefix, err)
		log.Printf("%s pactl set-source-volume: {SOURCE_NAME: %s ; VOLUME: %s}\n", errPrefix, sourceName, volumeInPercent)
	}
}

func MoveSourceOutput(sourceOutputID string, sourceName string) {
	errPrefix := "Error [MoveSourceOutput()]"

	cmd := exec.Command("pactl", "move-source-output", sourceOutputID, sourceName)
	_, err := cmd.Output()
	if err != nil {
		log.Printf("%s pactl move-source-output: %s\n", errPrefix, err)
		log.Printf("%s pactl move-source-output: {SOURCE_OUTPUT_ID: %s ; SOURCE_NAME: %s}\n", errPrefix, sourceOutputID, sourceName)
	}
}

func parseOutput(output string) Output {
	idRe, _ := regexp.Compile(`Sink #(\d+)`)
	nameRe, _ := regexp.Compile(`Name: (.+)`)
	descRe, _ := regexp.Compile(`Description: (.+)`)
	volumeRe, _ := regexp.Compile(`Volume: .+?(\d+)%`)
	muteRe, _ := regexp.Compile(`Mute: (yes|no)`)

	id, _ := strconv.Atoi(idRe.FindStringSubmatch(output)[1])
	name := nameRe.FindStringSubmatch(output)[1]
	desc := descRe.FindStringSubmatch(output)[1]
	volume, _ := strconv.Atoi(volumeRe.FindStringSubmatch(output)[1])
	mute := muteRe.FindStringSubmatch(output)[1] == "yes"

	return Output{
		ID:     id,
		Name:   name,
		Label:  desc,
		Volume: volume,
		Muted:  mute,
	}
}

func GetOutputs() ([]Output, error) {
	cmd := exec.Command("pactl", "list", "sinks")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	sinks := strings.Split(string(out), "Sink #")
	outputs := make([]Output, 0, len(sinks)-1)

	for _, sink := range sinks[1:] {
		outputs = append(outputs, parseOutput("Sink #"+sink))
	}

	return outputs, nil
}

func GetApps() []App {
	cmd := exec.Command("pactl", "list", "sink-inputs")
	out, _ := cmd.Output()

	re, _ := regexp.Compile(`Sink Input #(\d+)[\s\S]*?Sink: (\d+)[\s\S]*?Mute: (yes|no)[\s\S]*?Volume:.*?(\d+)%[\s\S]*?application\.name = "(.*?)"`)
	matches := re.FindAllStringSubmatch(string(out), -1)

	apps := make([]App, len(matches))
	for i, m := range matches {
		id, _ := strconv.Atoi(m[1])
		outputID, _ := strconv.Atoi(m[2])
		volume, _ := strconv.Atoi(m[4])
		apps[i] = App{
			ID:       id,
			OutputID: outputID,
			Label:    m[5],
			Volume:   volume,
			Muted:    m[3] == "yes",
		}
	}

	return apps
}

func parseSources(output string) Source {
	idRe, _ := regexp.Compile(`Source #(\d+)`)
	nameRe, _ := regexp.Compile(`Name: (.+)`)
	descRe, _ := regexp.Compile(`Description: (.+)`)
	volumeRe, _ := regexp.Compile(`Volume: .+?(\d+)%`)
	muteRe, _ := regexp.Compile(`Mute: (yes|no)`)

	id, _ := strconv.Atoi(idRe.FindStringSubmatch(output)[1])
	name := nameRe.FindStringSubmatch(output)[1]
	desc := descRe.FindStringSubmatch(output)[1]
	volume, _ := strconv.Atoi(volumeRe.FindStringSubmatch(output)[1])
	mute := muteRe.FindStringSubmatch(output)[1] == "yes"

	return Source{
		ID:     id,
		Name:   name,
		Label:  desc,
		Volume: volume,
		Muted:  mute,
	}
}

func GetSources() ([]Source, error) {
	cmd := exec.Command("pactl", "list", "sources")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	source := strings.Split(string(out), "Source #")
	sources := make([]Source, 0, len(source)-1)

	for _, sink := range source[1:] {
		sources = append(sources, parseSources("Source #"+sink))
	}

	return sources, nil
}

func ListenForChanges(callback func()) {
	cmd := exec.Command("pactl", "subscribe")
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "sink") || strings.Contains(line, "server") {
			callback()
		}
	}
}

func GetStatus() Status {
	errPrefix := "ERROR [GetStatus()]"

	outputs, err := GetOutputs()
	if err != nil {
		log.Printf("%s GetOutputs(): %s", errPrefix, err)
	}

	sources, err := GetSources()
	if err != nil {
		log.Printf("%s GetSources(): %s", errPrefix, err)
	}

	apps := GetApps()

	bi := buildinfo.Get()

	return Status{
		Outputs:   outputs,
		Apps:      apps,
		Sources:   sources,
		BuildInfo: *bi,
	}
}
