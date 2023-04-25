package protocols

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func reconnect(config Config, client mqtt.Client, options *mqtt.ClientOptions) {
	go func() {
		for !Exit {
			if !client.IsConnectionOpen() {
				fmt.Println("Disconnected from broker. Trying to reconnect...")
				x, ma := 1, 60
				for i := 1; i < 10; i++ {
					log.Printf("Reconnecting in %ds.\n", x)
					time.Sleep(time.Duration(x) * time.Second)

					token := client.Connect()
					for !token.WaitTimeout(100 * time.Millisecond) {
					}
					if err := token.Error(); err == nil && client.IsConnectionOpen() {
						log.Println("Reconnected!")
						break
					}
					x *= 2
					if x > ma {
						x = ma
					}
				}
				if !client.IsConnectionOpen() {
					log.Println("max try! Exit.")
					Exit = true
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func connectByMQTT(config Config) mqtt.Client {
	opts := mqtt.NewClientOptions()
	broker := fmt.Sprintf("tcp://%s:%d", config.Host, config.Port)
	opts.AddBroker(broker)
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return client
}

func connectByMQTTS(config Config) mqtt.Client {
	var tlsConfig tls.Config
	if config.Tls && config.CaCert == "" {
		log.Fatalln("TLS field in config is required")
	}
	certpool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(config.CaCert)
	if err != nil {
		log.Fatalln(err.Error())
	}
	certpool.AppendCertsFromPEM(ca)
	tlsConfig.RootCAs = certpool

	opts := mqtt.NewClientOptions()
	broker := fmt.Sprintf("ssl://%s:%d", config.Host, config.Port)
	println(broker)
	opts.AddBroker(broker)
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)
	opts.SetTLSConfig(&tlsConfig)

	opts.SetKeepAlive(3 * time.Second)
	opts.SetMaxReconnectInterval(3 * time.Second)

	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(100 * time.Millisecond) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	reconnect(config, client, opts)
	return client
}

func connectByWS(config Config) mqtt.Client {
	opts := mqtt.NewClientOptions()
	broker := fmt.Sprintf("ws://%s:%d/mqtt", config.Host, config.Port)
	opts.AddBroker(broker)
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return client
}

func connectByWSS(config Config) mqtt.Client {
	var tlsConfig tls.Config
	if config.Tls && config.CaCert == "" {
		log.Fatalln("TLS field in config is required")
	}
	certpool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(config.CaCert)
	if err != nil {
		log.Fatalln(err.Error())
	}
	certpool.AppendCertsFromPEM(ca)
	tlsConfig.RootCAs = certpool

	opts := mqtt.NewClientOptions()
	broker := fmt.Sprintf("wss://%s:%d/mqtt", config.Host, config.Port)
	opts.AddBroker(broker)
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)
	opts.SetTLSConfig(&tlsConfig)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return client
}
