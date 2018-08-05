package main

import (
	"products/containers"
	"products/server"
	"products/services"
)

func main() {

	container := containers.BuildContainer()

	// Fire server
	err := container.Invoke(func(server *server.Server) {

		// Fire Kafka consumer
		e := container.Invoke(func(kafkaConsumer *services.KafkaConsumer) {
			go func() {
				kafkaConsumer.Run()
			}()
		})

		if e != nil {
			panic(e)
		}

		server.Run()
	})

	if err != nil {
		panic(err)
	}

}
