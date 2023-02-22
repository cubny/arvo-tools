# Makefile

# Go files
GO_FILES := $(wildcard *.go)

# Binary names
AVRO_ENCODE_BIN := avro_encode
AVRO_GENERATE_BIN := avro_generate

# Docker image names
AVRO_ENCODE_IMAGE := avro-encode
AVRO_GENERATE_IMAGE := avro-generate

# Build flags
BUILD_FLAGS := -ldflags="-s -w"

.PHONY: all clean build-docker-images

all: $(AVRO_ENCODE_BIN) $(AVRO_GENERATE_BIN)

clean:
	rm -f $(AVRO_ENCODE_BIN) $(AVRO_GENERATE_BIN)
	docker rmi $(AVRO_ENCODE_IMAGE) $(AVRO_GENERATE_IMAGE)

$(AVRO_ENCODE_BIN): $(GO_FILES)
	go build $(BUILD_FLAGS) -o $(AVRO_ENCODE_BIN)

$(AVRO_GENERATE_BIN): $(GO_FILES)
	go build $(BUILD_FLAGS) -o $(AVRO_GENERATE_BIN)

build-docker-images: $(AVRO_ENCODE_BIN) $(AVRO_GENERATE_BIN)
	docker build -t $(AVRO_ENCODE_IMAGE) -f Dockerfile.avro-encode .
	docker build -t $(AVRO_GENERATE_IMAGE) -f Dockerfile.avro-generate .
