// Copyright (c) 2008-2020, Hazelcast, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package codec

import (
	"github.com/hazelcast/hazelcast-go-client/internal/proto"
	"github.com/hazelcast/hazelcast-go-client/internal/proto/codec/internal"
	"github.com/hazelcast/hazelcast-go-client/serialization"
)

const (
	// hex: 0x013100
	MapExecuteWithPredicateCodecRequestMessageType = int32(78080)
	// hex: 0x013101
	MapExecuteWithPredicateCodecResponseMessageType = int32(78081)

	MapExecuteWithPredicateCodecRequestInitialFrameSize = proto.PartitionIDOffset + proto.IntSizeInBytes
)

// Applies the user defined EntryProcessor to the entries in the map which satisfies provided predicate.
// Returns the results mapped by each key in the map.
type mapExecuteWithPredicateCodec struct{}

var MapExecuteWithPredicateCodec mapExecuteWithPredicateCodec

func (mapExecuteWithPredicateCodec) EncodeRequest(name string, entryProcessor serialization.Data, predicate serialization.Data) *proto.ClientMessage {
	clientMessage := proto.NewClientMessageForEncode()
	clientMessage.SetRetryable(false)

	initialFrame := proto.NewFrame(make([]byte, MapExecuteWithPredicateCodecRequestInitialFrameSize))
	clientMessage.AddFrame(initialFrame)
	clientMessage.SetMessageType(MapExecuteWithPredicateCodecRequestMessageType)
	clientMessage.SetPartitionId(-1)

	internal.StringCodec.Encode(clientMessage, name)
	internal.DataCodec.Encode(clientMessage, entryProcessor)
	internal.DataCodec.Encode(clientMessage, predicate)

	return clientMessage
}

func (mapExecuteWithPredicateCodec) DecodeResponse(clientMessage *proto.ClientMessage) []proto.Pair {
	frameIterator := clientMessage.FrameIterator()
	// empty initial frame
	frameIterator.Next()

	return internal.EntryListCodec.DecodeForDataAndData(frameIterator)
}