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
	"github.com/hazelcast/hazelcast-go-client/core"
	"github.com/hazelcast/hazelcast-go-client/internal/proto"
	"github.com/hazelcast/hazelcast-go-client/internal/proto/codec/internal"
)

const (
	// hex: 0x1D0200
	PNCounterAddCodecRequestMessageType = int32(1901056)
	// hex: 0x1D0201
	PNCounterAddCodecResponseMessageType = int32(1901057)

	PNCounterAddCodecRequestDeltaOffset             = proto.PartitionIDOffset + proto.IntSizeInBytes
	PNCounterAddCodecRequestGetBeforeUpdateOffset   = PNCounterAddCodecRequestDeltaOffset + proto.LongSizeInBytes
	PNCounterAddCodecRequestTargetReplicaUUIDOffset = PNCounterAddCodecRequestGetBeforeUpdateOffset + proto.BooleanSizeInBytes
	PNCounterAddCodecRequestInitialFrameSize        = PNCounterAddCodecRequestTargetReplicaUUIDOffset + proto.UuidSizeInBytes

	PNCounterAddResponseValueOffset        = proto.ResponseBackupAcksOffset + proto.ByteSizeInBytes
	PNCounterAddResponseReplicaCountOffset = PNCounterAddResponseValueOffset + proto.LongSizeInBytes
)

// Adds a delta to the PNCounter value. The delta may be negative for a
// subtraction.
// <p>
// The invocation will return the replica timestamps (vector clock) which
// can then be sent with the next invocation to keep session consistency
// guarantees.
// The target replica is determined by the {@code targetReplica} parameter.
// If smart routing is disabled, the actual member processing the client
// message may act as a proxy.
type pncounterAddCodec struct{}

var PNCounterAddCodec pncounterAddCodec

func (pncounterAddCodec) EncodeRequest(name string, delta int64, getBeforeUpdate bool, replicaTimestamps []proto.Pair, targetReplicaUUID core.UUID) *proto.ClientMessage {
	clientMessage := proto.NewClientMessageForEncode()
	clientMessage.SetRetryable(false)

	initialFrame := proto.NewFrame(make([]byte, PNCounterAddCodecRequestInitialFrameSize))
	internal.FixSizedTypesCodec.EncodeLong(initialFrame.Content, PNCounterAddCodecRequestDeltaOffset, delta)
	internal.FixSizedTypesCodec.EncodeBoolean(initialFrame.Content, PNCounterAddCodecRequestGetBeforeUpdateOffset, getBeforeUpdate)
	internal.FixSizedTypesCodec.EncodeUUID(initialFrame.Content, PNCounterAddCodecRequestTargetReplicaUUIDOffset, targetReplicaUUID)
	clientMessage.AddFrame(initialFrame)
	clientMessage.SetMessageType(PNCounterAddCodecRequestMessageType)
	clientMessage.SetPartitionId(-1)

	internal.StringCodec.Encode(clientMessage, name)
	internal.EntryListUUIDLongCodec.Encode(clientMessage, replicaTimestamps)

	return clientMessage
}

func (pncounterAddCodec) DecodeResponse(clientMessage *proto.ClientMessage) (value int64, replicaTimestamps []proto.Pair, replicaCount int32) {
	frameIterator := clientMessage.FrameIterator()
	initialFrame := frameIterator.Next()

	value = internal.FixSizedTypesCodec.DecodeLong(initialFrame.Content, PNCounterAddResponseValueOffset)
	replicaCount = internal.FixSizedTypesCodec.DecodeInt(initialFrame.Content, PNCounterAddResponseReplicaCountOffset)
	replicaTimestamps = internal.EntryListUUIDLongCodec.Decode(frameIterator)

	return value, replicaTimestamps, replicaCount
}
