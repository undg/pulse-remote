package audiodeprecated

import (
	"errors"
	"log"

	"mrogalski.eu/go/pulseaudio"
)

// clientOpen don't forget to closeClient()
func clientOpen() pulseaudio.Client {
	client, err := pulseaudio.NewClient()
	if err != nil {
		panic(err)
	}

	return *client
}

func clientClose(c pulseaudio.Client) {
	defer c.Close()
}

type CardInfo struct {
	Name  string
	Index uint32
}

func GetCards() ([]CardInfo, error) {
	c := clientOpen()

	cards, err := c.Cards()
	if err != nil {
		log.Println("ERROR GetCards c.Volume", err)
		return nil, errors.New("ERROR GetCards c.Volume")
	}

	clientClose(c)

	cardsInfo := []CardInfo{}

	for _, card := range cards {
		cardInfo := CardInfo{
			Name:  card.Name,
			Index: card.Index,
		}
		cardsInfo = append(cardsInfo, cardInfo)
	}

	return cardsInfo, nil
}

type OutputsInfo struct {
	ActiveIndex int
	CardID      string
	CardName    string
	PortName    string
	Available   bool
	PortID      string
}

func GetOutputs() ([]OutputsInfo, error) {
	c := clientOpen()

	output, activeIndex, err := c.Outputs()
	if err != nil {
		log.Println("ERROR GetOutputs c.Volume", err)
		return nil, errors.New("ERROR GetOutputs c.Volume")
	}

	clientClose(c)

	outputsInfo := []OutputsInfo{}
	for _, output := range output {
		cardInfo := OutputsInfo{
			ActiveIndex: activeIndex,
			CardID:      output.CardID,
			CardName:    output.CardName,
			PortName:    output.PortName,
			Available:   output.Available,
			PortID:      output.PortID,
		}
		outputsInfo = append(outputsInfo, cardInfo)
	}

	return outputsInfo, nil
}
