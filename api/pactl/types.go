package pactl

import "github.com/undg/go-prapi/api/buildinfo"

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
	ID        int    `json:"id" doc:"Unique numeric identifier of the source"`
	Name      string `json:"name" doc:"Unique string identifier of the source"`
	Label     string `json:"label" doc:"Human-readable label for the source"`
	Volume    int    `json:"volume" doc:"Current volume level of the source"`
	Muted     bool   `json:"muted" doc:"Whether the source is muted"`
	Monitor   string `json:"monitor" doc:"Name of monitor source capturing this source's output"`
	Monitored bool   `json:"monitored" doc:"Whether source is being monitored"`
}

type App struct {
	ID       int    `json:"id" doc:"The id of the sink. Same  as name"`
	OutputID int    `json:"outputId" doc:"Id of parrent device, same as output.id"`
	Label    string `json:"label" doc:"Human-readable label for the sink"`
	Volume   int    `json:"volume" doc:"Current volume level of the sink"`
	Muted    bool   `json:"muted" doc:"Whether the sink is muted"`
}
