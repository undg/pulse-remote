//lint:file-ignore ST1003 Ignore underscore naming in generated code

package pactl

type PactlSourceJSON struct {
	ActivePort any     `json:"active_port"`
	Balance    float64 `json:"balance"`
	BaseVolume struct {
		DB           string  `json:"db"`
		Value        float64 `json:"value"`
		ValuePercent string  `json:"value_percent"`
	} `json:"base_volume"`
	ChannelMap  string   `json:"channel_map"`
	Description string   `json:"description"`
	Driver      string   `json:"driver"`
	Flags       []string `json:"flags"`
	Formats     []string `json:"formats"`
	Index       float64  `json:"index"`
	Latency     struct {
		Actual     float64 `json:"actual"`
		Configured float64 `json:"configured"`
	} `json:"latency"`
	MonitorSource string  `json:"monitor_source"`
	Mute          bool    `json:"mute"`
	Name          string  `json:"name"`
	OwnerModule   float64 `json:"owner_module"`
	Ports         []any   `json:"ports"`
	Properties    struct {
		Alsa_Card                        string `json:"alsa.card"`
		Alsa_cardName                    string `json:"alsa.card_name"`
		Alsa_Class                       string `json:"alsa.class"`
		Alsa_Components                  string `json:"alsa.components"`
		Alsa_Device                      string `json:"alsa.device"`
		Alsa_driverName                  string `json:"alsa.driver_name"`
		Alsa_ID                          string `json:"alsa.id"`
		Alsa_longCardName                string `json:"alsa.long_card_name"`
		Alsa_mixerName                   string `json:"alsa.mixer_name"`
		Alsa_Name                        string `json:"alsa.name"`
		Alsa_resolutionBits              string `json:"alsa.resolution_bits"`
		Alsa_Subclass                    string `json:"alsa.subclass"`
		Alsa_Subdevice                   string `json:"alsa.subdevice"`
		Alsa_subdeviceName               string `json:"alsa.subdevice_name"`
		Alsa_Sync_ID                     string `json:"alsa.sync.id"`
		Api_acp_autoPort                 string `json:"api.acp.auto-port"`
		API_Alsa_Card                    string `json:"api.alsa.card"`
		API_Alsa_Card_Longname           string `json:"api.alsa.card.longname"`
		API_Alsa_Card_Name               string `json:"api.alsa.card.name"`
		API_Alsa_Path                    string `json:"api.alsa.path"`
		API_Alsa_Pcm_Card                string `json:"api.alsa.pcm.card"`
		API_Alsa_Pcm_Stream              string `json:"api.alsa.pcm.stream"`
		Api_alsa_useAcp                  string `json:"api.alsa.use-acp"`
		API_Dbus_ReserveDevice1          string `json:"api.dbus.ReserveDevice1"`
		API_Dbus_ReserveDevice1_Priority string `json:"api.dbus.ReserveDevice1.Priority"`
		Audio_Channels                   string `json:"audio.channels"`
		Audio_Position                   string `json:"audio.position"`
		Card_Profile_Device              string `json:"card.profile.device"`
		Client_ID                        string `json:"client.id"`
		Clock_Name                       string `json:"clock.name"`
		Clock_quantumLimit               string `json:"clock.quantum-limit"`
		Device_API                       string `json:"device.api"`
		Device_Bus                       string `json:"device.bus"`
		Device_busID                     string `json:"device.bus-id"`
		Device_busPath                   string `json:"device.bus_path"`
		Device_Class                     string `json:"device.class"`
		Device_Description               string `json:"device.description"`
		Device_Enum_API                  string `json:"device.enum.api"`
		Device_iconName                  string `json:"device.icon_name"`
		Device_ID                        string `json:"device.id"`
		Device_Name                      string `json:"device.name"`
		Device_Nick                      string `json:"device.nick"`
		Device_Plugged_Usec              string `json:"device.plugged.usec"`
		Device_Product_ID                string `json:"device.product.id"`
		Device_Product_Name              string `json:"device.product.name"`
		Device_Profile_Description       string `json:"device.profile.description"`
		Device_Profile_Name              string `json:"device.profile.name"`
		Device_Profile_Pro               string `json:"device.profile.pro"`
		Device_Routes                    string `json:"device.routes"`
		Device_Serial                    string `json:"device.serial"`
		Device_String                    string `json:"device.string"`
		Device_Subsystem                 string `json:"device.subsystem"`
		Device_Vendor_ID                 string `json:"device.vendor.id"`
		Device_Vendor_Name               string `json:"device.vendor.name"`
		Factory_ID                       string `json:"factory.id"`
		Factory_Name                     string `json:"factory.name"`
		Library_Name                     string `json:"library.name"`
		Media_Class                      string `json:"media.class"`
		Node_Driver                      string `json:"node.driver"`
		Node_Loop_Name                   string `json:"node.loop.name"`
		Node_Name                        string `json:"node.name"`
		Node_Nick                        string `json:"node.nick"`
		Node_pauseOnIdle                 string `json:"node.pause-on-idle"`
		Object_ID                        string `json:"object.id"`
		Object_Path                      string `json:"object.path"`
		Object_Serial                    string `json:"object.serial"`
		Port_Group                       string `json:"port.group"`
		Priority_Driver                  string `json:"priority.driver"`
		Priority_Session                 string `json:"priority.session"`
		Sysfs_Path                       string `json:"sysfs.path"`
	} `json:"properties"`
	SampleSpecification string `json:"sample_specification"`
	State               string `json:"state"`
	Volume              struct {
		Aux0 struct {
			DB           string  `json:"db"`
			Value        float64 `json:"value"`
			ValuePercent string  `json:"value_percent"`
		} `json:"aux0"`
	} `json:"volume"`
}
