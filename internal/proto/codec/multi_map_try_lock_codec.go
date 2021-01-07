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
	// hex: 0x021100
	MultiMapTryLockCodecRequestMessageType = int32(135424)
	// hex: 0x021101
	MultiMapTryLockCodecResponseMessageType = int32(135425)

	MultiMapTryLockCodecRequestThreadIdOffset    = proto.PartitionIDOffset + proto.IntSizeInBytes
	MultiMapTryLockCodecRequestLeaseOffset       = MultiMapTryLockCodecRequestThreadIdOffset + proto.LongSizeInBytes
	MultiMapTryLockCodecRequestTimeoutOffset     = MultiMapTryLockCodecRequestLeaseOffset + proto.LongSizeInBytes
	MultiMapTryLockCodecRequestReferenceIdOffset = MultiMapTryLockCodecRequestTimeoutOffset + proto.LongSizeInBytes
	MultiMapTryLockCodecRequestInitialFrameSize  = MultiMapTryLockCodecRequestReferenceIdOffset + proto.LongSizeInBytes

	MultiMapTryLockResponseResponseOffset = proto.ResponseBackupAcksOffset + proto.ByteSizeInBytes
)

// Tries to acquire the lock for the specified key for the specified lease time. After lease time, the lock will be
// released. If the lock is not available, then the current thread becomes disabled for thread scheduling purposes
// and lies dormant until one of two things happens:the lock is acquired by the current thread, or the specified
// waiting time elapses.
type multimapTryLockCodec struct{}

var MultiMapTryLockCodec multimapTryLockCodec

func (multimapTryLockCodec) EncodeRequest(name string, key serialization.Data, threadId int64, lease int64, timeout int64, referenceId int64) *proto.ClientMessage {
	clientMessage := proto.NewClientMessageForEncode()
	clientMessage.SetRetryable(true)

	initialFrame := proto.NewFrame(make([]byte, MultiMapTryLockCodecRequestInitialFrameSize))
	internal.FixSizedTypesCodec.EncodeLong(initialFrame.Content, MultiMapTryLockCodecRequestThreadIdOffset, threadId)
	internal.FixSizedTypesCodec.EncodeLong(initialFrame.Content, MultiMapTryLockCodecRequestLeaseOffset, lease)
	internal.FixSizedTypesCodec.EncodeLong(initialFrame.Content, MultiMapTryLockCodecRequestTimeoutOffset, timeout)
	internal.FixSizedTypesCodec.EncodeLong(initialFrame.Content, MultiMapTryLockCodecRequestReferenceIdOffset, referenceId)
	clientMessage.AddFrame(initialFrame)
	clientMessage.SetMessageType(MultiMapTryLockCodecRequestMessageType)
	clientMessage.SetPartitionId(-1)

	internal.StringCodec.Encode(clientMessage, name)
	internal.DataCodec.Encode(clientMessage, key)

	return clientMessage
}

func (multimapTryLockCodec) DecodeResponse(clientMessage *proto.ClientMessage) bool {
	frameIterator := clientMessage.FrameIterator()
	initialFrame := frameIterator.Next()

	return internal.FixSizedTypesCodec.DecodeBoolean(initialFrame.Content, MultiMapTryLockResponseResponseOffset)
}