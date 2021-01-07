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
	"github.com/hazelcast/hazelcast-go-client/serialization"
)

const (
	// hex: 0x0D0D00
	ReplicatedMapAddEntryListenerCodecRequestMessageType = int32(855296)
	// hex: 0x0D0D01
	ReplicatedMapAddEntryListenerCodecResponseMessageType = int32(855297)

	// hex: 0x0D0D02
	ReplicatedMapAddEntryListenerCodecEventEntryMessageType = int32(855298)

	ReplicatedMapAddEntryListenerCodecRequestLocalOnlyOffset  = proto.PartitionIDOffset + proto.IntSizeInBytes
	ReplicatedMapAddEntryListenerCodecRequestInitialFrameSize = ReplicatedMapAddEntryListenerCodecRequestLocalOnlyOffset + proto.BooleanSizeInBytes

	ReplicatedMapAddEntryListenerResponseResponseOffset                  = proto.ResponseBackupAcksOffset + proto.ByteSizeInBytes
	ReplicatedMapAddEntryListenerEventEntryEventTypeOffset               = proto.PartitionIDOffset + proto.IntSizeInBytes
	ReplicatedMapAddEntryListenerEventEntryUuidOffset                    = ReplicatedMapAddEntryListenerEventEntryEventTypeOffset + proto.IntSizeInBytes
	ReplicatedMapAddEntryListenerEventEntryNumberOfAffectedEntriesOffset = ReplicatedMapAddEntryListenerEventEntryUuidOffset + proto.UuidSizeInBytes
)

// Adds an entry listener for this map. The listener will be notified for all map add/remove/update/evict events.
type replicatedmapAddEntryListenerCodec struct{}

var ReplicatedMapAddEntryListenerCodec replicatedmapAddEntryListenerCodec

func (replicatedmapAddEntryListenerCodec) EncodeRequest(name string, localOnly bool) *proto.ClientMessage {
	clientMessage := proto.NewClientMessageForEncode()
	clientMessage.SetRetryable(false)

	initialFrame := proto.NewFrame(make([]byte, ReplicatedMapAddEntryListenerCodecRequestInitialFrameSize))
	internal.FixSizedTypesCodec.EncodeBoolean(initialFrame.Content, ReplicatedMapAddEntryListenerCodecRequestLocalOnlyOffset, localOnly)
	clientMessage.AddFrame(initialFrame)
	clientMessage.SetMessageType(ReplicatedMapAddEntryListenerCodecRequestMessageType)
	clientMessage.SetPartitionId(-1)

	internal.StringCodec.Encode(clientMessage, name)

	return clientMessage
}

func (replicatedmapAddEntryListenerCodec) DecodeResponse(clientMessage *proto.ClientMessage) core.UUID {
	frameIterator := clientMessage.FrameIterator()
	initialFrame := frameIterator.Next()

	return internal.FixSizedTypesCodec.DecodeUUID(initialFrame.Content, ReplicatedMapAddEntryListenerResponseResponseOffset)
}

func (replicatedmapAddEntryListenerCodec) Handle(clientMessage *proto.ClientMessage, handleEntryEvent func(key serialization.Data, value serialization.Data, oldValue serialization.Data, mergingValue serialization.Data, eventType int32, uuid core.UUID, numberOfAffectedEntries int32)) {
	messageType := clientMessage.GetMessageType()
	frameIterator := clientMessage.FrameIterator()
	if messageType == ReplicatedMapAddEntryListenerCodecEventEntryMessageType {
		initialFrame := frameIterator.Next()
		eventType := internal.FixSizedTypesCodec.DecodeInt(initialFrame.Content, ReplicatedMapAddEntryListenerEventEntryEventTypeOffset)
		uuid := internal.FixSizedTypesCodec.DecodeUUID(initialFrame.Content, ReplicatedMapAddEntryListenerEventEntryUuidOffset)
		numberOfAffectedEntries := internal.FixSizedTypesCodec.DecodeInt(initialFrame.Content, ReplicatedMapAddEntryListenerEventEntryNumberOfAffectedEntriesOffset)
		key := internal.CodecUtil.DecodeNullableForData(frameIterator)
		value := internal.CodecUtil.DecodeNullableForData(frameIterator)
		oldValue := internal.CodecUtil.DecodeNullableForData(frameIterator)
		mergingValue := internal.CodecUtil.DecodeNullableForData(frameIterator)
		handleEntryEvent(key, value, oldValue, mergingValue, eventType, uuid, numberOfAffectedEntries)
		return
	}
}