//lint:file-ignore ST1003 Ignore underscore naming in generated code

package pactl

type PactlAppsJSON struct {
	Balance           float64 `json:"balance"`
	BufferLatencyUsec float64 `json:"buffer_latency_usec"`
	ChannelMap        string  `json:"channel_map"`
	Client            string  `json:"client"`
	Corked            bool    `json:"corked"`
	Driver            string  `json:"driver"`
	Format            string  `json:"format"`
	Index             float64 `json:"index"`
	Mute              bool    `json:"mute"`
	OwnerModule       any     `json:"owner_module"`
	Properties        struct {
		Adapt_follower_spaNode        string `json:"adapt.follower.spa-node"`
		Application_iconName          string `json:"application.icon_name"`
		Application_Language          string `json:"application.language"`
		Application_Name              string `json:"application.name"`
		Application_Process_Binary    string `json:"application.process.binary"`
		Application_Process_Host      string `json:"application.process.host"`
		Application_Process_ID        string `json:"application.process.id"`
		Application_process_machineID string `json:"application.process.machine_id"`
		Application_process_sessionID string `json:"application.process.session_id"`
		Application_Process_User      string `json:"application.process.user"`
		Client_API                    string `json:"client.api"`
		Client_ID                     string `json:"client.id"`
		Clock_quantumLimit            string `json:"clock.quantum-limit"`
		Factory_ID                    string `json:"factory.id"`
		Library_Name                  string `json:"library.name"`
		Media_Class                   string `json:"media.class"`
		Media_Name                    string `json:"media.name"`
		ModuleStreamRestore_id        string `json:"module-stream-restore.id"`
		Node_Autoconnect              string `json:"node.autoconnect"`
		Node_driverID                 string `json:"node.driver-id"`
		Node_Latency                  string `json:"node.latency"`
		Node_Loop_Name                string `json:"node.loop.name"`
		Node_Name                     string `json:"node.name"`
		Node_Rate                     string `json:"node.rate"`
		Node_wantDriver               string `json:"node.want-driver"`
		Object_ID                     string `json:"object.id"`
		Object_Register               string `json:"object.register"`
		Object_Serial                 string `json:"object.serial"`
		Port_Group                    string `json:"port.group"`
		Pulse_Attr_Maxlength          string `json:"pulse.attr.maxlength"`
		Pulse_Attr_Minreq             string `json:"pulse.attr.minreq"`
		Pulse_Attr_Prebuf             string `json:"pulse.attr.prebuf"`
		Pulse_Attr_Tlength            string `json:"pulse.attr.tlength"`
		Pulse_Server_Type             string `json:"pulse.server.type"`
		Stream_isLive                 string `json:"stream.is-live"`
		Window_X11_Display            string `json:"window.x11.display"`
	} `json:"properties"`
	ResampleMethod      string  `json:"resample_method"`
	SampleSpecification string  `json:"sample_specification"`
	Sink                float64 `json:"sink"`
	SinkLatencyUsec     float64 `json:"sink_latency_usec"`
	Volume              struct {
		FrontLeft struct {
			DB           string  `json:"db"`
			Value        float64 `json:"value"`
			ValuePercent string  `json:"value_percent"`
		} `json:"front-left"`
		FrontRight struct {
			DB           string  `json:"db"`
			Value        float64 `json:"value"`
			ValuePercent string  `json:"value_percent"`
		} `json:"front-right"`
	} `json:"volume"`
}
