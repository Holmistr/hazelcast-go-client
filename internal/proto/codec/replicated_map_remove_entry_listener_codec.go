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
	"github.com/hazelcast/hazelcast-go-client/v4/internal/core"
	"github.com/hazelcast/hazelcast-go-client/v4/internal/proto"
)

const (
	// hex: 0x0D0E00
	ReplicatedMapRemoveEntryListenerCodecRequestMessageType = int32(855552)
	// hex: 0x0D0E01
	ReplicatedMapRemoveEntryListenerCodecResponseMessageType = int32(855553)

	ReplicatedMapRemoveEntryListenerCodecRequestRegistrationIdOffset = proto.PartitionIDOffset + proto.IntSizeInBytes
	ReplicatedMapRemoveEntryListenerCodecRequestInitialFrameSize     = ReplicatedMapRemoveEntryListenerCodecRequestRegistrationIdOffset + proto.UuidSizeInBytes

	ReplicatedMapRemoveEntryListenerResponseResponseOffset = proto.ResponseBackupAcksOffset + proto.ByteSizeInBytes
)

// Removes the specified entry listener. If there is no such listener added before, this call does no change in the
// cluster and returns false.
type replicatedmapRemoveEntryListenerCodec struct{}

var ReplicatedMapRemoveEntryListenerCodec replicatedmapRemoveEntryListenerCodec

func (replicatedmapRemoveEntryListenerCodec) EncodeRequest(name string, registrationId core.UUID) *proto.ClientMessage {
	clientMessage := proto.NewClientMessageForEncode()
	clientMessage.SetRetryable(true)

	initialFrame := proto.NewFrame(make([]byte, ReplicatedMapRemoveEntryListenerCodecRequestInitialFrameSize))
	FixSizedTypesCodec.EncodeUUID(initialFrame.Content, ReplicatedMapRemoveEntryListenerCodecRequestRegistrationIdOffset, registrationId)
	clientMessage.AddFrame(initialFrame)
	clientMessage.SetMessageType(ReplicatedMapRemoveEntryListenerCodecRequestMessageType)
	clientMessage.SetPartitionId(-1)

	StringCodec.Encode(clientMessage, name)

	return clientMessage
}

func (replicatedmapRemoveEntryListenerCodec) DecodeResponse(clientMessage *proto.ClientMessage) bool {
	frameIterator := clientMessage.FrameIterator()
	initialFrame := frameIterator.Next()

	return FixSizedTypesCodec.DecodeBoolean(initialFrame.Content, ReplicatedMapRemoveEntryListenerResponseResponseOffset)
}
