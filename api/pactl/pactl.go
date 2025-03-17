package pactl

import (
	"bufio"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/undg/go-prapi/api/buildinfo"
	"github.com/undg/go-prapi/api/logger"
)

func SetSinkVolume(sinkName string, volume string) {
	setVolume("sink", sinkName, volume)
}

func SetSinkMuted(sinkName string, muted bool) {
	setMuted("sink", sinkName, muted)
}

func SetSinkInputVolume(sinkInputID string, volume string) {
	setVolume("sink-input", sinkInputID, volume)
}

func SetSinkInputMuted(sinkInputID string, muted bool) {
	setMuted("sink-input", sinkInputID, muted)
}

func MoveSinkInput(sinkInputID string, sinkName string) {
	moveApp("sink-input", sinkInputID, sinkName)

}

func SetSourceVolume(sourceName string, volume string) {
	setVolume("source", sourceName, volume)
}

func SetSourceMuted(sourceName string, muted bool) {
	setMuted("source", sourceName, muted)

}

func SetSourceInputVolume(sourceInputID string, volume string) {
	setVolume("source-input", sourceInputID, volume)
}

func SetSourceInputMuted(sourceInputID string, muted bool) {
	setMuted("source-input", sourceInputID, muted)
}

func MoveSourceOutput(sourceOutputID string, sourceName string) {
	moveApp("source-output", sourceOutputID, sourceName)
}

func parseSink(sinkName string) Sink {
	idRe, _ := regexp.Compile(`Sink #(\d+)`)
	nameRe, _ := regexp.Compile(`Name: (.+)`)
	descRe, _ := regexp.Compile(`Description: (.+)`)
	volumeRe, _ := regexp.Compile(`Volume: .+?(\d+)%`)
	muteRe, _ := regexp.Compile(`Mute: (yes|no)`)

	id, _ := strconv.Atoi(idRe.FindStringSubmatch(sinkName)[1])
	name := nameRe.FindStringSubmatch(sinkName)[1]
	desc := descRe.FindStringSubmatch(sinkName)[1]
	volume, _ := strconv.Atoi(volumeRe.FindStringSubmatch(sinkName)[1])
	mute := muteRe.FindStringSubmatch(sinkName)[1] == "yes"

	return Sink{
		ID:     id,
		Name:   name,
		Label:  desc,
		Volume: volume,
		Muted:  mute,
	}
}

func GetSinks() ([]Sink, error) {
	cmd := exec.Command("pactl", "list", "sinks")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	sinksNames := strings.Split(string(out), "Sink #")
	sinks := make([]Sink, 0, len(sinksNames)-1)

	for _, sink := range sinksNames[1:] {
		sinks = append(sinks, parseSink("Sink #"+sink))
	}

	return sinks, nil
}

func GetSinkInputs() ([]SinkInput, error) {
	cmd := exec.Command("pactl", "list", "sink-inputs")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	re, err := regexp.Compile(`Sink Input #(\d+)[\s\S]*?Sink: (\d+)[\s\S]*?Mute: (yes|no)[\s\S]*?Volume:.*?(\d+)%[\s\S]*?application\.name = "(.*?)"`)
	if err != nil {
		return nil, err
	}
	matches := re.FindAllStringSubmatch(string(out), -1)

	sinkInputs := make([]SinkInput, len(matches))
	for i, m := range matches {
		id, err := strconv.Atoi(m[1])
		if err != nil {
			return nil, err
		}
		sinkID, err := strconv.Atoi(m[2])
		if err != nil {
			return nil, err
		}
		volume, err := strconv.Atoi(m[4])
		if err != nil {
			return nil, err
		}
		sinkInputs[i] = SinkInput{
			ID:     id,
			SinkID: sinkID,
			Label:  m[5],
			Volume: volume,
			Muted:  m[3] == "yes",
		}
	}

	return sinkInputs, nil
}

func parseSources(sourceName string) Source {
	idRe, _ := regexp.Compile(`Source #(\d+)`)
	nameRe, _ := regexp.Compile(`Name: (.+)`)
	descRe, _ := regexp.Compile(`Description: (.+)`)
	volumeRe, _ := regexp.Compile(`Volume: .+?(\d+)%`)
	muteRe, _ := regexp.Compile(`Mute: (yes|no)`)
	monitorRe, _ := regexp.Compile(`Monitor of Sink: (.+)`) // n/a or name of the Sink

	id, _ := strconv.Atoi(idRe.FindStringSubmatch(sourceName)[1])
	name := nameRe.FindStringSubmatch(sourceName)[1]
	desc := descRe.FindStringSubmatch(sourceName)[1]
	volume, _ := strconv.Atoi(volumeRe.FindStringSubmatch(sourceName)[1])
	muted := muteRe.FindStringSubmatch(sourceName)[1] == "yes"
	monitored := monitorRe.FindStringSubmatch(sourceName)[1] != "n/a"
	monitor := monitorRe.FindStringSubmatch(sourceName)[1]

	return Source{
		ID:        id,
		Name:      name,
		Label:     desc,
		Volume:    volume,
		Muted:     muted,
		Monitor:   monitor,
		Monitored: monitored,
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
	errPrefix := "ERROR [GetStatus()] -> "

	sinks, err := GetSinks()
	if err != nil {
		logger.Error().Err(err).Msgf("%s GetSinks()", errPrefix)
	}

	sources, err := GetSources()
	if err != nil {
		logger.Error().Err(err).Msgf("%s GetSources()", errPrefix)
	}

	sinkInputs, err := GetSinkInputs()
	if err != nil {
		logger.Error().Err(err).Msgf("%s GetSinkInputs()", errPrefix)
	}

	bi := buildinfo.Get()

	return Status{
		Sinks:      sinks,
		SinkInputs: sinkInputs,
		Sources:    sources,
		BuildInfo:  *bi,
	}
}
