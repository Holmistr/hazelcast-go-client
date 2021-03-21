package sql

import "github.com/hazelcast/hazelcast-go-client/v4/internal"

type SqlError struct {

	code int32
	message string
	originatingMemberId internal.UUID

}

func NewSqlError(code int32, message string, originatingMemberId internal.UUID) *SqlError {
	return &SqlError{code: code, message: message, originatingMemberId: originatingMemberId}
}

func (s *SqlError) OriginatingMemberId() internal.UUID {
	return s.originatingMemberId
}

func (s *SqlError) SetOriginatingMemberId(originatingMemberId internal.UUID) {
	s.originatingMemberId = originatingMemberId
}

func (s *SqlError) Message() string {
	return s.message
}

func (s *SqlError) SetMessage(message string) {
	s.message = message
}

func (s *SqlError) Code() int32 {
	return s.code
}

func (s *SqlError) SetCode(code int32) {
	s.code = code
}

//func NewSqlError() *sqlError {
//
//}
