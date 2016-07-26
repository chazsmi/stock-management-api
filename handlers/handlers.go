package handlers

import (
	"log"

	"github.com/chazsmi/socket-events/events"
	"github.com/chazsmi/stock-service/proto"
	"github.com/micro/go-micro/broker"
	"github.com/micro/protobuf/proto"

	"golang.org/x/net/websocket"
)

type Events struct {
	EventRec      chan events.Event
	EventsHandler *events.Handler
}

var topic = "charlieplc.topic.stock"

func (c Events) Stock(ws *websocket.Conn) {
	req := ws.Request()
	req.ParseForm()
	sku := req.Form.Get("sku")

	conn := &events.Connection{
		Ref:  sku,
		Ws:   ws,
		Done: make(chan bool),
	}

	c.EventsHandler.RegisterEvent(conn)

	for {
		select {
		case <-conn.Done:
			// Kill the connection upon termination
			log.Println("Killing connection")
			break
		}
	}
}

func (c Events) Sub() {
	_, err := broker.Subscribe(topic, func(p broker.Publication) error {
		st := &stock.StockReadResponse{}
		err := proto.Unmarshal(p.Message().Body, st)
		if err != nil {
			log.Println("proto unmarshaling error: ", err.Error())
			return err
		}
		// Construct an events struct and send to the chnannel
		c.EventRec <- events.Event{
			Ref: &st.Sku,
			Data: map[string]interface{}{
				"Sku":    st.Sku,
				"Amount": st.Amount,
			},
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
}
