package routes

import (
	"display-csv-import/src/models"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func getPersonContainerFromContext(c *gin.Context) *models.DataContainer[models.Person] {
	// since personContainer will be set in main
	// and is available in all routes
	// we don't need to check if it exists
	personContainer, _ := c.Get("personContainer")
	personContainerTyped := personContainer.(*models.DataContainer[models.Person])
	return personContainerTyped
}

func getPersonChannelsFromContext(c *gin.Context) *models.ChannelsContainer[models.Person] {
	// since personChannels will be set in main
	// and is available in all routes
	// we don't need to check if it exists
	personChannels, _ := c.Get("personChannels")
	personChannelsTyped := personChannels.(*models.ChannelsContainer[models.Person])
	return personChannelsTyped
}

func getPersons(c *gin.Context) {
	personContainer := getPersonContainerFromContext(c)
	personContainer.RWMutex.RLock()
	defer personContainer.RWMutex.RUnlock()
	c.JSON(200, personContainer.Data)
}

func handleWebSocket(conn *websocket.Conn, personChannels *models.ChannelsContainer[models.Person]) {
	defer conn.Close()
	personChannel := personChannels.Add()

	go func() {
		for {
			_, _, err := conn.ReadMessage()
			// includes client leaving the web page
			// or manually closing the websocket
			if err != nil {
				personChannels.RemoveAndClose(personChannel)
				break
			}
		}
	}()

	for person := range personChannel {
		// this does not fail if the connection has been closed
		err := conn.WriteJSON(person)
		if err != nil {
			log.Println("Failed to write message to websocket:", err)
			break
		}
	}
}

func upgradeWebsocketConnection(c *gin.Context) {
	personChannels := getPersonChannelsFromContext(c)
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade to websocket:", err)
		return
	}

	handleWebSocket(conn, personChannels)
}

var upgrader websocket.Upgrader = websocket.Upgrader{}

func RegisterRoutes(engine *gin.Engine) {
	engine.Static("/static", "web-root")

	engine.GET("/persons", getPersons)
	engine.GET("/persons/live", upgradeWebsocketConnection)

	engine.NoRoute(func(ctx *gin.Context) {
		ctx.File("web-root/index.html")
	})
}
