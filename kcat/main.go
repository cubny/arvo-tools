package main

import (
	"context"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/linkedin/goavro"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "usage: %s <avro-schema-file> <json-payload-file> <kafka-topic>\n", os.Args[0])
		os.Exit(1)
	}

	schemaFile := os.Args[1]
	payloadFile := os.Args[2]
	topic := os.Args[3]

	// Load the Avro schema from the file
	schemaJSON, err := os.ReadFile(schemaFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read Avro schema file: %s\n", err)
		os.Exit(1)
	}
	codec, err := goavro.NewCodec(string(schemaJSON))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create Avro codec: %s\n", err)
		os.Exit(1)
	}

	// Load the JSON payload from the file and convert it to Avro binary format
	payloadJSON, err := os.ReadFile(payloadFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read JSON payload file: %s\n", err)
		os.Exit(1)
	}
	n, _, err := codec.NativeFromTextual(payloadJSON)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to convert JSON to Avro native format: %s\n", err)
		os.Exit(1)
	}
	avroBinary, err := codec.BinaryFromNative(nil, n)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to convert Avro native format to binary: %s\n", err)
		os.Exit(1)
	}

	// Set up the Kafka producer configuration
	conf := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	}

	// Create the Kafka producer
	p, err := kafka.NewProducer(conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create Kafka producer: %s\n", err)
		os.Exit(1)
	}
	defer p.Close()

	// Produce the message to the Kafka topic
	deliveryChan := make(chan kafka.Event, 1)
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: avroBinary,
	}, deliveryChan)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to produce message to Kafka topic: %s\n", err)
		os.Exit(1)
	}

	// Wait for the delivery report or an error
	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		fmt.Fprintf(os.Stderr, "failed to deliver message: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("delivered message to topic %s [%d] at offset %v\n", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}
}
