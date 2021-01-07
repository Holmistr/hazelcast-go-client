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
	// hex: 0x030E00
	QueueCompareAndRetainAllCodecRequestMessageType = int32(200192)
	// hex: 0x030E01
	QueueCompareAndRetainAllCodecResponseMessageType = int32(200193)

	QueueCompareAndRetainAllCodecRequestInitialFrameSize = proto.PartitionIDOffset + proto.IntSizeInBytes

	QueueCompareAndRetainAllResponseResponseOffset = proto.ResponseBackupAcksOffset + proto.ByteSizeInBytes
)

// Retains only the elements in this collection that are contained in the specified collection (optional operation).
// In other words, removes from this collection all of its elements that are not contained in the specified collection.
type queueCompareAndRetainAllCodec struct{}

var QueueCompareAndRetainAllCodec queueCompareAndRetainAllCodec

func (queueCompareAndRetainAllCodec) EncodeRequest(name string, dataList []serialization.Data) *proto.ClientMessage {
	clientMessage := proto.NewClientMessageForEncode()
	clientMessage.SetRetryable(false)

	initialFrame := proto.NewFrame(make([]byte, QueueCompareAndRetainAllCodecRequestInitialFrameSize))
	clientMessage.AddFrame(initialFrame)
	clientMessage.SetMessageType(QueueCompareAndRetainAllCodecRequestMessageType)
	clientMessage.SetPartitionId(-1)

	internal.StringCodec.Encode(clientMessage, name)
	internal.ListMultiFrameCodec.EncodeForData(clientMessage, dataList)

	return clientMessage
}

func (queueCompareAndRetainAllCodec) DecodeResponse(clientMessage *proto.ClientMessage) bool {
	frameIterator := clientMessage.FrameIterator()
	initialFrame := frameIterator.Next()

	return internal.FixSizedTypesCodec.DecodeBoolean(initialFrame.Content, QueueCompareAndRetainAllResponseResponseOffset)
}