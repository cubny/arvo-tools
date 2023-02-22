package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/actgardner/gogen-avro/v8/generator"
)

const defaultNamespace = "default"

func main() {
	// Check that the correct number of command line arguments were provided
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <avro-schema-file>")
		return
	}

	// Read in the Avro schema from the provided file
	schemaFile, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error opening schema file: %v\n", err)
		return
	}
	defer schemaFile.Close()

	schema, err := generator.ParseSchema(schemaFile)
	if err != nil {
		fmt.Printf("Error parsing schema: %v\n", err)
		return
	}

	// Generate a random example JSON message
	rand.Seed(time.Now().UnixNano())
	example := generateRandomRecord(schema)

	// Print the example JSON message to standard output
	jsonBytes, err := json.MarshalIndent(example, "", "    ")
	if err != nil {
		fmt.Printf("Error encoding example message as JSON: %v\n", err)
		return
	}
	os.Stdout.Write(jsonBytes)
}

func generateRandomRecord(schema *generator.RecordDefinition) map[string]interface{} {
	record := make(map[string]interface{})
	for _, field := range schema.Fields {
		value, _ := generateRandomValue(field.Type, field.Nullable, field.NameSpace)
		record[field.Name] = value
	}
	return record
}

func generateRandomValue(fieldType generator.FieldType, nullable bool, namespace string) (interface{}, error) {
	if nullable && rand.Float32() < 0.2 {
		return nil, nil
	}

	switch fieldType.Type {
	case generator.TypeNull:
		return nil, nil

	case generator.TypeBoolean:
		return rand.Float32() < 0.5, nil

	case generator.TypeInt:
		return rand.Int31(), nil

	case generator.TypeLong:
		return rand.Int63(), nil

	case generator.TypeFloat:
		return rand.Float32(), nil

	case generator.TypeDouble:
		return rand.Float64(), nil

	case generator.TypeString:
		return generateRandomString(), nil

	case generator.TypeBytes:
		return generateRandomBytes(), nil

	case generator.TypeRecord:
		record := generateRandomRecord(fieldType.Record())
		if namespace != "" {
			for k, v := range record {
				delete(record, k)
				record[getFullName(namespace, k)] = v
			}
		}
		return record, nil

	case generator.TypeArray:
		return generateRandomArray(fieldType.Items, namespace)

	case generator.TypeMap:
		return generateRandomMap(fieldType.Values, namespace)

	default:
		return nil, fmt.Errorf("unsupported type: %v", fieldType.Type)
	}
}

func generateRandomString() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, rand.Intn(10)+5)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func generateRandomBytes() []byte {
	b := make([]byte, rand.Intn(10)+5)
	rand.Read(b)
	return b
}

func generateRandomArray(itemType generator.FieldType, namespace string) ([]interface{}, error) {
	length := rand.Intn(5) + 1
	array := make([]interface{}, length)
	for i := 0; i < length; i++ {
		item, err := generateRandomValue(itemType, false, namespace)
		if err != nil {
			return nil, err
		}
		array[i] = item
	}
	return array, nil
}

func generateRandomMap(valueType generator.FieldType, namespace string) (map[string]interface{}, error) {
	length := rand.Intn(5) + 1
	m := make(map[string]interface{}, length)
	for i := 0; i < length; i++ {
		key := generateRandomString()
		value, err := generateRandomValue(valueType, false, namespace)
		if err != nil {
			return nil, err
		}
		m[key] = value
	}
	return m, nil
}
