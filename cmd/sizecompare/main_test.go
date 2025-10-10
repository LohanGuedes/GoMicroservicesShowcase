package main

import (
	"encoding/json"
	"encoding/xml"
	"testing"

	"google.golang.org/protobuf/proto"
)

func BenchmarkSerializetoJSON(b *testing.B) {
	for b.Loop() {
		json.Marshal(metadata)
	}
}

func BenchmarkSerializetoXML(b *testing.B) {
	for b.Loop() {
		xml.Marshal(metadata)
	}
}

func BenchmarkSerializetoProto(b *testing.B) {
	for b.Loop() {
		proto.Marshal(genMetadata)
	}
}
