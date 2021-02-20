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
	"github.com/hazelcast/hazelcast-go-client/v4/internal/proto"

	"github.com/hazelcast/hazelcast-go-client/v4/internal/serialization"
)

const (
	// hex: 0x013700
	MapFetchKeysCodecRequestMessageType = int32(79616)
	// hex: 0x013701
	MapFetchKeysCodecResponseMessageType = int32(79617)

	MapFetchKeysCodecRequestBatchOffset      = proto.PartitionIDOffset + proto.IntSizeInBytes
	MapFetchKeysCodecRequestInitialFrameSize = MapFetchKeysCodecRequestBatchOffset + proto.IntSizeInBytes
)

// Fetches specified number of keys from the specified partition starting from specified table index.
type mapFetchKeysCodec struct{}

var MapFetchKeysCodec mapFetchKeysCodec

func (mapFetchKeysCodec) EncodeRequest(name string, iterationPointers []proto.Pair, batch int32) *proto.ClientMessage {
	clientMessage := proto.NewClientMessageForEncode()
	clientMessage.SetRetryable(true)

	initialFrame := proto.NewFrame(make([]byte, MapFetchKeysCodecRequestInitialFrameSize))
	FixSizedTypesCodec.EncodeInt(initialFrame.Content, MapFetchKeysCodecRequestBatchOffset, batch)
	clientMessage.AddFrame(initialFrame)
	clientMessage.SetMessageType(MapFetchKeysCodecRequestMessageType)
	clientMessage.SetPartitionId(-1)

	StringCodec.Encode(clientMessage, name)
	EntryListIntegerIntegerCodec.Encode(clientMessage, iterationPointers)

	return clientMessage
}

func (mapFetchKeysCodec) DecodeResponse(clientMessage *proto.ClientMessage) (iterationPointers []proto.Pair, keys []serialization.Data) {
	frameIterator := clientMessage.FrameIterator()
	frameIterator.Next()

	iterationPointers = EntryListIntegerIntegerCodec.Decode(frameIterator)
	keys = ListMultiFrameCodec.DecodeForData(frameIterator)

	return iterationPointers, keys
}
