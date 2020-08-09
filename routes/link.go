package routes

import (
	"fmt"
	"go-web-console/exceptions"
	"go-web-console/sshclient"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func v1ApiPing(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

type linkRequest struct {
	Hostname   string `json:"hostname"`
	Port       int    `json:"port"`
	Username   string `json:"username"`
	Type       string `json:"type"`
	Password   string `json:"password,omitempty"`
	PrivateKey string `json:"private_key,omitempty"`
	Passphrase string `json:"passphrase,omimtempty"`
}

type linkResponse struct {
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
}

func v1ApiLink(c *gin.Context) {
	var json linkRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		c.Error(exceptions.NewRespErrorWithStr(http.StatusBadRequest, "can not parse json: %s", err.Error()))
		return
	}
	log.Print("test")
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
	log.Print("test")
	c.JSON(http.StatusOK, linkResponse{
		ClientID:     client.ClientID,
		ClientSecret: client.ClientSecret,
	})
}
