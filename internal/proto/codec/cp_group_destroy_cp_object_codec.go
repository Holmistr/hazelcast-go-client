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
)

const (
	// hex: 0x1E0200
	CPGroupDestroyCPObjectCodecRequestMessageType = int32(1966592)
	// hex: 0x1E0201
	CPGroupDestroyCPObjectCodecResponseMessageType = int32(1966593)

	CPGroupDestroyCPObjectCodecRequestInitialFrameSize = proto.PartitionIDOffset + proto.IntSizeInBytes
)

// Destroys the distributed object with the given name on the requested
// CP group
type cpgroupDestroyCPObjectCodec struct{}

var CPGroupDestroyCPObjectCodec cpgroupDestroyCPObjectCodec

func (cpgroupDestroyCPObjectCodec) EncodeRequest(groupId proto.RaftGroupId, serviceName string, objectName string) *proto.ClientMessage {
	clientMessage := proto.NewClientMessageForEncode()
	clientMessage.SetRetryable(true)

	initialFrame := proto.NewFrame(make([]byte, CPGroupDestroyCPObjectCodecRequestInitialFrameSize))
	clientMessage.AddFrame(initialFrame)
	clientMessage.SetMessageType(CPGroupDestroyCPObjectCodecRequestMessageType)
	clientMessage.SetPartitionId(-1)

	internal.RaftGroupIdCodec.Encode(clientMessage, groupId)
	internal.StringCodec.Encode(clientMessage, serviceName)
	internal.StringCodec.Encode(clientMessage, objectName)

	return clientMessage
}