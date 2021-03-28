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
)


const(
    // hex: 0x170400
    RingbufferCapacityCodecRequestMessageType  = int32(1508352)
    // hex: 0x170401
    RingbufferCapacityCodecResponseMessageType = int32(1508353)

    RingbufferCapacityCodecRequestInitialFrameSize = proto.PartitionIDOffset + proto.IntSizeInBytes

    RingbufferCapacityResponseResponseOffset = proto.ResponseBackupAcksOffset + proto.ByteSizeInBytes
)

// Returns the capacity of this Ringbuffer.

func EncodeRingbufferCapacityRequest(name string) *proto.ClientMessage {
    clientMessage := proto.NewClientMessageForEncode()
    clientMessage.SetRetryable(true)

    initialFrame := proto.NewFrameWith(make([]byte, RingbufferCapacityCodecRequestInitialFrameSize), proto.UnfragmentedMessage)
    clientMessage.AddFrame(initialFrame)
    clientMessage.SetMessageType(RingbufferCapacityCodecRequestMessageType)
    clientMessage.SetPartitionId(-1)

    EncodeString(clientMessage, name)

    return clientMessage
}

func DecodeRingbufferCapacityResponse(clientMessage *proto.ClientMessage) int64 {
    frameIterator := clientMessage.FrameIterator()
    initialFrame := frameIterator.Next()

    return FixSizedTypesCodec.DecodeLong(initialFrame.Content, RingbufferCapacityResponseResponseOffset)
}