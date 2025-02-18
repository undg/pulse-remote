package pactl

import (
	"os/exec"
	"regexp"
	"testing"
	"time"
)

func playsilence() *exec.Cmd {
	cmd := exec.Command("paplay", "--raw", "--channels=2", "--format=s16le", "--rate=44100", "/dev/zero")
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second / 100) // Give time to appear in pactl
	return cmd
}

func TestGetStatus(t *testing.T) {
	t.Run("Outputs", func(t *testing.T) {
		status := GetStatus()
		if status.Outputs == nil {
			t.Errorf("Missing Outputs in Status struct")
		}
	})
	t.Run("Apps", func(t *testing.T) {
		cmd := playsilence()
		defer cmd.Process.Kill()
		status := GetStatus()
		if status.Apps == nil {
			t.Errorf("Missing Apps in Status struct")
		}
	})
}

func TestGetOutputs(t *testing.T) {
	t.Run("Volume", func(t *testing.T) {
		sink, _ := GetOutputs()
		if sink[0].Volume < 0 {
			t.Errorf("Expected volume more than 0, but got %d", sink[0].Volume)
		}
	})
	namePattern := `^(alsa_output|bluez_sink|bluez_output|combined)\..*`
	re := regexp.MustCompile(namePattern)

	tests := []struct {
		Name  string
		Input string
		Want  bool
	}{
		{"ValidSink", "alsa_output.pci-0000_0c_00.4.analog-stereo", true},
		{"InvalidSink", "dupa", false},
		{"InvalidEmptySink", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			sink, _ := GetOutputs()
			if got := re.MatchString(tt.Input); got != tt.Want {
				t.Errorf("Expected name should starts with pattern, but got [%s]", sink[0].Name)
			}
		})
	}
}

func TestGetSources(t *testing.T) {
	t.Run("Volume", func(t *testing.T) {
		source, _ := GetSources()
		if source[0].Volume < 0 {
			t.Errorf("Expected volume more than 0, but got %d", source[0].Volume)
		}
	})
	namePattern := `^(alsa_output|bluez_sink|bluez_output|combined)\..*`
	re := regexp.MustCompile(namePattern)

	tests := []struct {
		Name  string
		Input string
		Want  bool
	}{
		{"ValidSink", "alsa_output.pci-0000_0c_00.4.analog-stereo", true},
		{"InvalidSink", "dupa", false},
		{"InvalidEmptySink", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			sink, _ := GetOutputs()
			if got := re.MatchString(tt.Input); got != tt.Want {
				t.Errorf("Expected name should starts with pattern, but got [%s]", sink[0].Name)
			}
		})
	}
}
