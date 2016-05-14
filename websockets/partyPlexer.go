package websockets

import (
	"encoding/json"
	"fmt"
)

func parseResponse(stringifiedJSON []byte) map[string]interface{} {
	var messageInterface interface{}
  err := json.Unmarshal(stringifiedJSON, &messageInterface); if err != nil {
  	fmt.Println(err)
  }
  incomingMessage := messageInterface.(map[string]interface{})
  return incomingMessage
}

func RouteIncomingMessage(stringifiedJSON []byte, c *Connection) {

}

