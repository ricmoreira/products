package config

import (
	"os"
	"strconv"
	"strings"
)

var (
	// GENERAL
	HOST           string = "HOST"
	MONGO_HOST     string = "MONGO_HOST"
	MONGO_DATABASE string = "MONGO_DATABASE"

	// KAFKA
	GROUP_ID             string = "GROUP_ID"
	TOPICS_SUBSCRIBED    string = "TOPICS_SUBSCRIBED"
	BOOTSTRAP_SERVERS    string = "BOOTSTRAP_SERVERS"
	REQUEST_TIMEOUT      string = "REQUEST_TIMEOUT"
	RETRIES              string = "RETRIES"
	BATCH_SIZE           string = "BATCH_SIZE"
	LINGER               string = "LINGER"
	BUFFER_MEMORY        string = "BUFFER_MEMORY"
	AUTO_COMMIT_INTERVAL string = "AUTO_COMMIT_INTERVAL"
	AUTO_COMMIT_ENABLE   string = "AUTO_COMMIT_ENABLE"
	AUTO_OFFSET_RESET    string = "AUTO_OFFSET_RESET"
)

type Config struct {
	Host              string
	MongoHost         string
	MongoDatabaseName string
	*KafkaConsumerConfig
}

// KafkaConsumerConfig is the struct that represents Kafka parameters for the Kafka consumer
// these parameters are explained in more detail in: https://kafka.apache.org/documentation/#newconsumerconfigs
type KafkaConsumerConfig struct {
	GroupID            string   // a unique string that identifies the Connect cluster group this worker belongs to
	TopicsSubscribed   []string // topics subscribed
	BootstrapServers   string   // kafka brokers endpoints (separated by ",")
	RequestTimeout     int      // controls the maximum amount of time the client will wait for the response of a request
	Retries            int      // number of retries if response to request failed
	BatchSize          int      // size of the aggregation on records aggregation
	Linger             int      // delay before the producer batch requests to Kafka
	BufferMemory       int      // total bytes of memory the producer can use to buffer records waiting to be sent to the server
	AutoCommitInterval int      // the frequency in milliseconds that the consumer offsets are auto-committed to Kafka
	AutoCommitEnable   bool     // if true, periodically commit to ZooKeeper the offset of messages already fetched by the consumer

	AutoOffsetReset string // what to do when there is no initial offset in ZooKeeper or if an offset is out of range
	// IMPORTANT: if set to "earliest", this means that if Kafka loses its commit history, some events may dealed twice.
	// This is a trade-of for trying not to loose any event produced to Kafka
}

func NewConfig() *Config {
	topicsEnv := MustGetEnv(TOPICS_SUBSCRIBED)
	topics := strings.Split(topicsEnv, ",")

	bootServ := MustGetEnv(BOOTSTRAP_SERVERS)
	reqTimeOut, _ := strconv.Atoi(MustGetEnv(REQUEST_TIMEOUT))
	retries, _ := strconv.Atoi(MustGetEnv(RETRIES))
	batchSize, _ := strconv.Atoi(MustGetEnv(BATCH_SIZE))
	linger, _ := strconv.Atoi(MustGetEnv(LINGER))
	bufferMemory, _ := strconv.Atoi(MustGetEnv(BUFFER_MEMORY))
	autoCommitInt, _ := strconv.Atoi(MustGetEnv(AUTO_COMMIT_INTERVAL))
	autoCommitEnable, _ := strconv.ParseBool(MustGetEnv(AUTO_COMMIT_ENABLE))
	autoOffsetReset := MustGetEnv(AUTO_OFFSET_RESET)

	kafkaConfig := &KafkaConsumerConfig{
		GroupID:            MustGetEnv(GROUP_ID),
		TopicsSubscribed:   topics,
		BootstrapServers:   bootServ,
		RequestTimeout:     reqTimeOut,
		Retries:            retries,
		BatchSize:          batchSize,
		Linger:             linger,
		BufferMemory:       bufferMemory,
		AutoCommitInterval: autoCommitInt,
		AutoCommitEnable:   autoCommitEnable,
		AutoOffsetReset:    autoOffsetReset,
	}

	return &Config{
		Host:              MustGetEnv(HOST),
		MongoHost:         MustGetEnv(MONGO_HOST),
		MongoDatabaseName: MustGetEnv(MONGO_DATABASE),
		KafkaConsumerConfig:    kafkaConfig,
	}
}

func MustGetEnv(envVarName string) string {
	res, found := os.LookupEnv(envVarName)

	if !found {
		panic("Environment variable " + envVarName + " not found")
	}

	return res
}
