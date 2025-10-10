package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	"github.com/lohanguedes/movie-microservices/gen"
	"github.com/lohanguedes/movie-microservices/metadata/pkg/model"
	"google.golang.org/protobuf/proto"
)

var metadata = &model.Metadata{
	ID:          "123",
	Title:       "The Movie 2",
	Description: "Sequel of the legendary movie \"The MOVIE\"!",
	Director:    "Foo Bar Baz - AKA FBB",
}

var genMetadata = &gen.Metadata{
	Id:          "123",
	Title:       "The Movie 2",
	Description: "Sequel of the legendary movie \"The MOVIE\"!",
	Director:    "Foo Bar Baz - AKA FBB",
}

// Test application, ignoring errors
func main() {
	jsonBytes, _ := json.Marshal(metadata)
	xmlBytes, _ := xml.Marshal(metadata)
	protoBytes, _ := proto.Marshal(genMetadata)

	fmt.Printf("Json size:\t%dB\n", len(jsonBytes))
	fmt.Printf("XML size:\t%dB\n", len(xmlBytes))
	fmt.Printf("proto size:\t%dB\n", len(protoBytes))
}
