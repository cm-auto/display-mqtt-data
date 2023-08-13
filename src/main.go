package main

import (
	"display-csv-import/src/models"
	mqttlogic "display-csv-import/src/mqtt_logic"
	"display-csv-import/src/routes"
	"fmt"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	engine := gin.Default()

	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "3000"
	}
	mqttServerAddress, ok := os.LookupEnv("MQTT_SERVER_ADDRESS")
	if !ok {
		panic("Env variable 'MQTT_SERVER_ADDRESS' not set")
	}
	topic, ok := os.LookupEnv("MQTT_TOPIC")
	if !ok {
		panic("Env variable 'MQTT_TOPIC' not set")
	}

	engine.RedirectTrailingSlash = true

	personContainer := models.DataContainer[models.Person]{
		RWMutex: sync.RWMutex{},
		Data:    []models.Person{},
	}

	personChannels := models.ChannelsContainer[models.Person]{}

	engine.Use(func(c *gin.Context) {
		c.Set("personContainer", &personContainer)
		c.Set("personChannels", &personChannels)
	})

	routes.RegisterRoutes(engine)

	client, err := mqttlogic.CreateClientAndConnect(mqttServerAddress)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(250)

	messageHandler := mqttlogic.GeneratePersonMessageHandler(&personChannels, &personContainer)

	mqttlogic.Subscribe(client, topic, messageHandler)

	address := fmt.Sprintf(":%s", port)
	engine.Run(address)
}
