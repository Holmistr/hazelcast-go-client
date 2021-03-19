/*
* Copyright (c) 2008-2020, Hazelcast, Inc. All Rights Reserved.
*
* Licensed under the Apache License, Version 2.0 (the "License")
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/
package codec

import (
    "github.com/hazelcast/hazelcast-go-client/v4/internal/proto"
    "github.com/hazelcast/hazelcast-go-client/v4/internal/sql"
)


const(
    SqlColumnMetadataCodecTypeFieldOffset = 0
    SqlColumnMetadataCodecNullableFieldOffset = SqlColumnMetadataCodecTypeFieldOffset + proto.IntSizeInBytes
    SqlColumnMetadataCodecNullableInitialFrameSize = SqlColumnMetadataCodecNullableFieldOffset + proto.BooleanSizeInBytes
)

/*
type sqlcolumnmetadataCodec struct {}

var SqlColumnMetadataCodec sqlcolumnmetadataCodec
*/

func EncodeSqlColumnMetadata(clientMessage *proto.ClientMessage, sqlColumnMetadata sql.SqlColumnMetadata){
    clientMessage.AddFrame(proto.BeginFrame.Copy())
    initialFrame := proto.NewFrame(make([]byte,SqlColumnMetadataCodecNullableInitialFrameSize))
    FixSizedTypesCodec.EncodeInt(initialFrame.Content, SqlColumnMetadataCodecTypeFieldOffset, int32(sqlColumnMetadata.Type()))
    FixSizedTypesCodec.EncodeBoolean(initialFrame.Content, SqlColumnMetadataCodecNullableFieldOffset, sqlColumnMetadata.Nullable())
    clientMessage.AddFrame(initialFrame)

    EncodeString(clientMessage, sqlColumnMetadata.Name())

    clientMessage.AddFrame(proto.EndFrame.Copy())
}

func DecodeSqlColumnMetadata(frameIterator *proto.ForwardFrameIterator) sql.SqlColumnMetadata {
    // begin frame
    frameIterator.Next()
    initialFrame := frameIterator.Next()
    _type := FixSizedTypesCodec.DecodeInt(initialFrame.Content, SqlColumnMetadataCodecTypeFieldOffset)
    isNullableExists := false
    let nullable = false    if (initialFrame.content.length >= NULLABLE_OFFSET + BitsUtil.BOOLEAN_SIZE_IN_BYTES) {
    nullable = FixSizedTypesCodec.DecodeBoolean(initialFrame.Content, NULLABLE_OFFSET)
    isNullableExists = true
    }

    name := DecodeString(frameIterator)
    CodecUtil.FastForwardToEndFrame(frameIterator)
    return sql.NewSqlColumnMetadata(name, _type, isNullableExists, nullable)
}