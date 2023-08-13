package mqttlogic

import (
	"display-csv-import/src/models"
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// the following two functions return or take in a client
// and not a pointer to one, because the methods of the
// clients are not callable if a pointer is passed
func CreateClientAndConnect(server string) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions().AddBroker(server)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return client, nil
}

func Subscribe(client mqtt.Client, topic string, handler func(client mqtt.Client, msg mqtt.Message)) error {
	token := client.Subscribe(topic, 0, handler)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func GeneratePersonMessageHandler(
	websocketChannels *models.ChannelsContainer[models.Person],
	personContainer *models.DataContainer[models.Person],
) func(client mqtt.Client, msg mqtt.Message) {
	return func(client mqtt.Client, msg mqtt.Message) {
		person := models.Person{}
		if err := json.Unmarshal(msg.Payload(), &person); err != nil {
			fmt.Printf("error: %v", err)
			return
		}
		websocketChannels.Send(person)
		personContainer.RWMutex.Lock()
		defer personContainer.RWMutex.Unlock()
		personContainer.Data = append(personContainer.Data, person)
	}
}
