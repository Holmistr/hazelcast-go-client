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
	// hex: 0x020E00
	MultiMapAddEntryListenerCodecRequestMessageType = int32(134656)
	// hex: 0x020E01
	MultiMapAddEntryListenerCodecResponseMessageType = int32(134657)

	// hex: 0x020E02
	MultiMapAddEntryListenerCodecEventEntryMessageType = int32(134658)

	MultiMapAddEntryListenerCodecRequestIncludeValueOffset = proto.PartitionIDOffset + proto.IntSizeInBytes
	MultiMapAddEntryListenerCodecRequestLocalOnlyOffset    = MultiMapAddEntryListenerCodecRequestIncludeValueOffset + proto.BooleanSizeInBytes
	MultiMapAddEntryListenerCodecRequestInitialFrameSize   = MultiMapAddEntryListenerCodecRequestLocalOnlyOffset + proto.BooleanSizeInBytes

	MultiMapAddEntryListenerResponseResponseOffset                  = proto.ResponseBackupAcksOffset + proto.ByteSizeInBytes
	MultiMapAddEntryListenerEventEntryEventTypeOffset               = proto.PartitionIDOffset + proto.IntSizeInBytes
	MultiMapAddEntryListenerEventEntryUuidOffset                    = MultiMapAddEntryListenerEventEntryEventTypeOffset + proto.IntSizeInBytes
	MultiMapAddEntryListenerEventEntryNumberOfAffectedEntriesOffset = MultiMapAddEntryListenerEventEntryUuidOffset + proto.UuidSizeInBytes
)

// Adds an entry listener for this multimap. The listener will be notified for all multimap add/remove/update/evict events.
type multimapAddEntryListenerCodec struct{}

var MultiMapAddEntryListenerCodec multimapAddEntryListenerCodec

func (multimapAddEntryListenerCodec) EncodeRequest(name string, includeValue bool, localOnly bool) *proto.ClientMessage {
	clientMessage := proto.NewClientMessageForEncode()
	clientMessage.SetRetryable(false)

	initialFrame := proto.NewFrame(make([]byte, MultiMapAddEntryListenerCodecRequestInitialFrameSize))
	internal.FixSizedTypesCodec.EncodeBoolean(initialFrame.Content, MultiMapAddEntryListenerCodecRequestIncludeValueOffset, includeValue)
	internal.FixSizedTypesCodec.EncodeBoolean(initialFrame.Content, MultiMapAddEntryListenerCodecRequestLocalOnlyOffset, localOnly)
	clientMessage.AddFrame(initialFrame)
	clientMessage.SetMessageType(MultiMapAddEntryListenerCodecRequestMessageType)
	clientMessage.SetPartitionId(-1)

	internal.StringCodec.Encode(clientMessage, name)

	return clientMessage
}

func (multimapAddEntryListenerCodec) DecodeResponse(clientMessage *proto.ClientMessage) core.UUID {
	frameIterator := clientMessage.FrameIterator()
	initialFrame := frameIterator.Next()

	return internal.FixSizedTypesCodec.DecodeUUID(initialFrame.Content, MultiMapAddEntryListenerResponseResponseOffset)
}

func (multimapAddEntryListenerCodec) Handle(clientMessage *proto.ClientMessage, handleEntryEvent func(key serialization.Data, value serialization.Data, oldValue serialization.Data, mergingValue serialization.Data, eventType int32, uuid core.UUID, numberOfAffectedEntries int32)) {
	messageType := clientMessage.GetMessageType()
	frameIterator := clientMessage.FrameIterator()
	if messageType == MultiMapAddEntryListenerCodecEventEntryMessageType {
		initialFrame := frameIterator.Next()
		eventType := internal.FixSizedTypesCodec.DecodeInt(initialFrame.Content, MultiMapAddEntryListenerEventEntryEventTypeOffset)
		uuid := internal.FixSizedTypesCodec.DecodeUUID(initialFrame.Content, MultiMapAddEntryListenerEventEntryUuidOffset)
		numberOfAffectedEntries := internal.FixSizedTypesCodec.DecodeInt(initialFrame.Content, MultiMapAddEntryListenerEventEntryNumberOfAffectedEntriesOffset)
		key := internal.CodecUtil.DecodeNullableForData(frameIterator)
		value := internal.CodecUtil.DecodeNullableForData(frameIterator)
		oldValue := internal.CodecUtil.DecodeNullableForData(frameIterator)
		mergingValue := internal.CodecUtil.DecodeNullableForData(frameIterator)
		handleEntryEvent(key, value, oldValue, mergingValue, eventType, uuid, numberOfAffectedEntries)
		return
	}
}