package sql

import (
	"fmt"
	"github.com/hazelcast/hazelcast-go-client/v4/hazelcast/protocol/codec"
	"github.com/hazelcast/hazelcast-go-client/v4/internal"
	"github.com/hazelcast/hazelcast-go-client/v4/internal/cluster"
	"github.com/hazelcast/hazelcast-go-client/v4/internal/proto"
	"github.com/hazelcast/hazelcast-go-client/v4/internal/proxy"
	"github.com/hazelcast/hazelcast-go-client/v4/internal/serialization"
)

type SqlService struct {
	proxy proxy.Impl
	connectionManager cluster.ConnectionManager
}

func NewSqlService(connectionManager cluster.ConnectionManager) SqlService {
	return SqlService{connectionManager: connectionManager}
}

func (s *SqlService) Execute(command string) bool {
	fmt.Println("Executing SQL service: ", command)

	memberId := internal.NewUUID()
	localId := internal.NewUUID()
	queryId := NewQueryId(memberId.MostSignificantBits(), memberId.LeastSignificantBits(), localId.MostSignificantBits(), localId.LeastSignificantBits())
	//
	params := []serialization.Data{}
	requestMessage := codec.EncodeSqlExecuteRequest(command, nil, -1, 4096, "", 0, queryId)
	fmt.Println(requestMessage)
	//
	//s.invoke(requestMessage)

	return true
}

func (s *SqlService) invoke(request *proto.ClientMessage) (*proto.ClientMessage, error) {
	inv := s.proxy.InvocationFactory.NewInvocationOnRandomTarget(request)
	s.proxy.RequestCh <- inv
	return inv.Get()
	//select {
	//case p.requestCh <- inv:
	//	return inv.GetWithTimeout(100 * time.Millisecond)
	//case <-time.After(100 * time.Millisecond):
	//	return nil, errors.New("timeout")

	//}

}
