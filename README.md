# Stock Streaming API

An example streaming API that montors RabbitMQ for changes in stock levels and then exposes a websocket endpoint for clients.
This service is an example service in the following scenero
<p align="center">
  <img src="stock.png" />
</p>
### Dependencies
- Go Web
- Consul
- RabbitMQ

### To run the service
- Make sure you have Consul installed and running
- Make yourself a copy of the config yaml file specifying the location of your Rabbit instance
``` bash
 cp config.example.yml config.yml
```
- Run the service
``` Go
 go run main.go 
```
:8080/stock will listen for all updates
:8080/stock?sku=12345 listen for stock uodates for sku 12345
