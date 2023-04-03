package turnserver

import (
	"context"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/gravitl/netmaker/logger"
	"github.com/gravitl/netmaker/servercfg"
	"github.com/pion/turn/v2"
)

var (
	UsersMap = make(map[string][]byte)
)

func RegisterNewHostWithTurn(hostID, hostPass string) {
	UsersMap[hostID] = turn.GenerateAuthKey(hostID, servercfg.GetTurnHost(), hostPass)
}

func Start(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	// Create a UDP listener to pass into pion/turn
	// pion/turn itself doesn't allocate any UDP sockets, but lets the user pass them in
	// this allows us to add logging, storage or modify inbound/outbound traffic
	udpListener, err := net.ListenPacket("udp4", "0.0.0.0:"+strconv.Itoa(servercfg.GetTurnPort()))
	if err != nil {
		log.Panicf("Failed to create TURN server listener: %s", err)
	}

	s, err := turn.NewServer(turn.ServerConfig{
		Realm: servercfg.GetTurnHost(),
		// Set AuthHandler callback
		// This is called every time a user tries to authenticate with the TURN server
		// Return the key for that user, or false when no user is found
		AuthHandler: func(username string, realm string, srcAddr net.Addr) ([]byte, bool) {
			if key, ok := UsersMap[username]; ok {
				return key, true
			}
			return nil, false
		},
		// PacketConnConfigs is a list of UDP Listeners and the configuration around them
		PacketConnConfigs: []turn.PacketConnConfig{
			{
				PacketConn: udpListener,
				RelayAddressGenerator: &turn.RelayAddressGeneratorStatic{
					RelayAddress: net.ParseIP("64.227.178.89"), // Claim that we are listening on IP passed by user (This should be your Public IP)
					Address:      "0.0.0.0",                    // But actually be listening on every interface
				},
			},
		},
	})
	if err != nil {
		log.Panic(err)
	}
	go func() {
		for {
			time.Sleep(time.Second * 10)
			log.Print(s.AllocationCount())

		}
	}()

	// Block until user sends SIGINT or SIGTERM
	<-ctx.Done()
	logger.Log(0, "## Stopping Turn Server...")
	if err = s.Close(); err != nil {
		log.Panic(err)
	}
}