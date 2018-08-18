package services

import (
	"log"
	"encoding/json"
	"products/config"
	"products/models/request"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaConsumer struct {
	config      *config.Config
	productServ *ProductService
}

func NewKafkaConsumer(config *config.Config, ps *ProductService) *KafkaConsumer {
	return &KafkaConsumer{
		config:      config,
		productServ: ps,
	}
}

func (kc *KafkaConsumer) Run() {

	log.Println("Start receiving from Kafka")

	configConsumer := kafka.ConfigMap{
		"bootstrap.servers":       kc.config.BootstrapServers,
		"group.id":                kc.config.GroupID,
		"auto.offset.reset":       kc.config.AutoOffsetReset,
		"auto.commit.enable":      kc.config.AutoCommitEnable,
		"auto.commit.interval.ms": kc.config.AutoCommitInterval,
	}

	c, err := kafka.NewConsumer(&configConsumer)

	if err != nil {
		panic(err)
	}

	topicsSubs := kc.config.TopicsSubscribed
	err = c.SubscribeTopics(topicsSubs, nil)

	if err != nil {
		panic(err)
	}

	for {
		msg, err := c.ReadMessage(-1)

		if err == nil {

			topic := *msg.TopicPartition.Topic

			switch topic {
			case "products":
				log.Println("Reading a products message")
				products, err := kc.parseProductsMessage(msg.Value)
				if err != nil {
					log.Printf("Error parsing event message value. Message %v \n Error: %s\n", msg.Value, err.Error())
					break
				}

				// save products to database
				_, e := kc.productServ.CreateMany(products)
				if e != nil {
					log.Printf("Error saving products to database\n Error: %s\n", e.Response)
					break
				}
			default: //ignore any other topics
			}
		} else {
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	c.Close()
}

func (kc *KafkaConsumer) parseProductsMessage(messageValue []byte) (*[]*mrequest.ProductCreate, error) {
	products := make([]*mrequest.ProductCreate, 0)
	err := json.Unmarshal(messageValue, &products)

	if err != nil {
		return nil, err
	}

	return &products, nil
}
