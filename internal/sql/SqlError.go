package sql

import "github.com/hazelcast/hazelcast-go-client/v4/internal"

type sqlError struct {

	code int32
	message string
	originatingMemberId internal.UUID

}

//func NewSqlError() *sqlError {
//
//}
