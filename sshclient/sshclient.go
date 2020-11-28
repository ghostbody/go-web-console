package sshclient

import (
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
)

// SSHClient represents the client linked to target server
type SSHClient struct {
	Client       *ssh.Client
	ClientID     string
	ClientSecret string
}

// SSHClientsManager manages all ssh clients
type SSHClientsManager struct {
	clients map[string]*SSHClient
}

var manager *SSHClientsManager = &SSHClientsManager{
	clients: map[string]*SSHClient{},
}

// GetSSHClientsManager get the singleton NewSSHClientsManager
func GetSSHClientsManager() *SSHClientsManager {
	return manager
}

// Connet2TargetWithPassword connect to the target server with password
func (cm *SSHClientsManager) Connet2TargetWithPassword(
	username string, targetAddr string,
	password string,
) (*SSHClient, error) {
	var authMethods []ssh.AuthMethod
	authMethods = []ssh.AuthMethod{ssh.Password(password)}
	return cm.Connect2TargetWithAuthMethods(username, targetAddr, authMethods)
}

// Connet2TargetWithPrivateKey connect to the target server with private key
func (cm *SSHClientsManager) Connet2TargetWithPrivateKey(
	username string, targetAddr string,
	privateKey string, passphrase string,
) (*SSHClient, error) {
	method, err := PublicKeyAuthFunc(privateKey, passphrase)
	if err != nil {
		return nil, err
	}
	authMethods := []ssh.AuthMethod{method}
	return cm.Connect2TargetWithAuthMethods(username, targetAddr, authMethods)
}

// Connect2TargetWithAuthMethods connect to target server with ssh.AuthMethods
func (cm *SSHClientsManager) Connect2TargetWithAuthMethods(
	username string, targetAddr string,
	authMethods []ssh.AuthMethod,
) (*SSHClient, error) {
	config := &ssh.ClientConfig{
		Timeout: time.Second * 3,
		User:    username,
		// TODO: maybe a more secure method
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            authMethods,
	}
	return cm.Connect2TargetWithConfig(username, targetAddr, config)
}

// Connect2TargetWithConfig connect to server with ssh.ClientConfig
func (cm *SSHClientsManager) Connect2TargetWithConfig(
	username string, targetAddr string,
	config *ssh.ClientConfig,
) (*SSHClient, error) {
	client, err := ssh.Dial("tcp", targetAddr, config)
	if err != nil {
		return nil, err
	}
	var clientID string
	for {
		clientID = uuid.New().String()
		if _, ok := cm.clients[clientID]; !ok {
			break
		}
	}
	sshClient := &SSHClient{
		Client:       client,
		ClientID:     clientID,
		ClientSecret: uuid.New().String(),
	}
	cm.clients[clientID] = sshClient
	return sshClient, nil
}

// GetSSHClients gets all clients
func (cm *SSHClientsManager) GetSSHClients() map[string]*SSHClient {
	return cm.clients
}

// GetSSHClientByClientID get client by id
func (cm *SSHClientsManager) GetSSHClientByClientID(clientID string) *SSHClient {
	client, ok := cm.clients[clientID]
	if ok {
		return client
	}
	return nil
}

// PublicKeyAuthFunc parses a string key (maybe with a passphrase) to ssh.AuthMethod
func PublicKeyAuthFunc(key string, passphrase string) (ssh.AuthMethod, error) {
	bkey := []byte(key)
	var signer ssh.Signer
	var err error
	if passphrase == "" {
		signer, err = ssh.ParsePrivateKey(bkey)
	} else {
		bpassphrase := []byte(passphrase)
		signer, err = ssh.ParsePrivateKeyWithPassphrase(bkey, bpassphrase)
	}
	if err != nil {
		log.Fatal("can not load ssh key!")
		return nil, err
	}
	return ssh.PublicKeys(signer), nil
}
