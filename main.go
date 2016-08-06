package main

import (
	"flag"
	"log"

	"github.com/chazsmi/socket-events/events"
	"github.com/chazsmi/stock-management-api/handlers"
	"github.com/chazsmi/stock-service/config"
	"github.com/micro/go-micro/broker"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
	"github.com/micro/go-web"
	"golang.org/x/net/websocket"
)

func main() {

	// Set up Config
	configFile := flag.String("c", "config.yml", "path to application config file")
	flag.Parse()
	c, err := config.ReadReturn(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	// Start listening for stock updates
	eventsHandler := events.NewEventHandler()
	e := handlers.Events{
		EventReceive:  eventsHandler.Receive,
		EventsHandler: eventsHandler,
	}
	go eventsHandler.Init()

	// Start listening to the msg queue
	if err := broker.Init(broker.Option(func(o *broker.Options) {
		o.Addrs = []string{c.Rabbit.Host}
	})); err != nil {
		log.Fatalf("Broker Init error: %v", err)
	}
	if err := broker.Connect(); err != nil {
		log.Fatalf("Broker Connect error: %v", err)
	}
	go e.Sub()

	service := web.NewService(
		web.Name("api.product.management"),
		web.Version("Latest"),
		web.Address("[::]:8080"),
		web.Advertise("[::]:8080"),
	)

	service.Handle("/stock", websocket.Handler(e.Stock))

	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
