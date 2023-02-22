package main

import (
"fmt"
"os"

"github.com/linkedin/goavro"
)

func main() {
	// Check that the correct number of command line arguments were provided
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <avro-schema-file> <json-payload-file>")
		return
	}

	// Read in the Avro schema from the provided file
	schemaFile, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error opening schema file: %v\n", err)
		return
	}
	defer schemaFile.Close()

	schema, err := goavro.NewParser().Parse(schemaFile)
	if err != nil {
		fmt.Printf("Error parsing schema: %v\n", err)
		return
	}

	// Read in the JSON payload from the provided file
	payloadFile, err := os.Open(os.Args[2])
	if err != nil {
		fmt.Printf("Error opening payload file: %v\n", err)
		return
	}
	defer payloadFile.Close()

	var payload interface{}
	err = goavro.NewDecoder(schema, payloadFile).Decode(&payload)
	if err != nil {
		fmt.Printf("Error decoding JSON payload: %v\n", err)
		return
	}

	// Encode the payload as Avro
	avroBytes, err := goavro.NewCodec(schema).BinaryFromNative(nil, payload)
	if err != nil {
		fmt.Printf("Error encoding Avro message: %v\n", err)
		return
	}

	// Write the Avro-encoded message to standard output
	os.Stdout.Write(avroBytes)
}

