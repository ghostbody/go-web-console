package routes

import (
	"fmt"
	"go-web-console/exceptions"
	"go-web-console/sshclient"
	"net/http"

	"github.com/gin-gonic/gin"
)

func v1ApiPing(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

type createLinkRequest struct {
	Hostname   string `json:"hostname"`
	Port       int    `json:"port"`
	Username   string `json:"username"`
	Type       string `json:"type"`
	Password   string `json:"password,omitempty"`
	PrivateKey string `json:"private_key,omitempty"`
	Passphrase string `json:"passphrase,omimtempty"`
}

type createLinkResponse struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type clientData struct {
	ClientID string `json:"clientID"`
}

type getLinkResponse struct {
	Clients map[string]clientData `json:"clients"`
}

func v1ApiCreateLink(c *gin.Context) {
	var json createLinkRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		c.Error(exceptions.NewRespErrorWithStr(http.StatusBadRequest, "can not parse json: %s", err.Error()))
		return
	}
	sshClientManager := sshclient.GetSSHClientsManager()
	// TODO: these params should be checked
	targetAddr := fmt.Sprintf("%s:%d", json.Hostname, json.Port)
	var err error
	var client *sshclient.SSHClient
	if json.Type == "password" {
		client, err = sshClientManager.Connet2TargetWithPassword(
			json.Username,
			targetAddr,
			json.Password,
		)
	} else if json.Type == "key" {
		client, err = sshClientManager.Connet2TargetWithPrivateKey(
			json.Username,
			targetAddr,
			json.PrivateKey,
			json.Passphrase,
		)
	} else {
		c.Error(exceptions.NewRespErrorWithStr(http.StatusBadRequest, "unkown auth type: %s", json.Type))
		return
	}
	if err != nil {
		c.Error(
			exceptions.NewRespErrorWithStr(
				http.StatusForbidden,
				"can not connect to ssh server: %s", err.Error(),
			),
		)
		return
	}
	c.JSON(http.StatusOK, createLinkResponse{
		ClientID:     client.ClientID,
		ClientSecret: client.ClientSecret,
	})
}

func v1ApiGetLinks(c *gin.Context) {
	sshClientManager := sshclient.GetSSHClientsManager()

	clients := sshClientManager.GetSSHClients()

	responseClients := make(map[string]clientData, len(clients))

	for clientID, client := range clients {
		responseClients[clientID] = clientData{
			ClientID: client.ClientID,
		}
	}

	c.JSON(http.StatusOK, getLinkResponse{
		Clients: responseClients,
	})

}
