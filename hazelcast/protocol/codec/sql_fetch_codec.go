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
    "github.com/hazelcast/hazelcast-go-client/v4/internal/sql"
)


const(
    // hex: 0x210500
    SqlFetchCodecRequestMessageType  = int32(2163968)
    // hex: 0x210501
    SqlFetchCodecResponseMessageType = int32(2163969)

    SqlFetchCodecRequestCursorBufferSizeOffset = proto.PartitionIDOffset + proto.IntSizeInBytes
    SqlFetchCodecRequestInitialFrameSize = SqlFetchCodecRequestCursorBufferSizeOffset + proto.IntSizeInBytes

)

// Fetches the next row page.

func EncodeSqlFetchRequest(queryId sql.QueryId, cursorBufferSize int32) *proto.ClientMessage {
    clientMessage := proto.NewClientMessageForEncode()
    clientMessage.SetRetryable(false)

    initialFrame := proto.NewFrameWith(make([]byte, SqlFetchCodecRequestInitialFrameSize), proto.UnfragmentedMessage)
    FixSizedTypesCodec.EncodeInt(initialFrame.Content, SqlFetchCodecRequestCursorBufferSizeOffset, cursorBufferSize)
    clientMessage.AddFrame(initialFrame)
    clientMessage.SetMessageType(SqlFetchCodecRequestMessageType)
    clientMessage.SetPartitionId(-1)

    EncodeSqlQueryId(clientMessage, queryId)

    return clientMessage
}

func DecodeSqlFetchResponse(clientMessage *proto.ClientMessage) (rowPage sql.SqlPage, error sql.SqlError) {
    frameIterator := clientMessage.FrameIterator()
    frameIterator.Next()

    rowPage = CodecUtil.DecodeNullableForSqlPage(frameIterator)
    //error = CodecUtil.DecodeNullableForSqlError(frameIterator)

    return rowPage, error }
