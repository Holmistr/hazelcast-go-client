package sql

type QueryId struct {

	memberIdHigh uint64
	memberIdLow uint64
	localIdHigh uint64
	localIdLow uint64

}

func (q *QueryId) LocalIdLow() uint64 {
	return q.localIdLow
}

func (q *QueryId) SetLocalIdLow(localIdLow uint64) {
	q.localIdLow = localIdLow
}

func (q *QueryId) LocalIdHigh() uint64 {
	return q.localIdHigh
}

func (q *QueryId) SetLocalIdHigh(localIdHigh uint64) {
	q.localIdHigh = localIdHigh
}

func (q *QueryId) MemberIdLow() uint64 {
	return q.memberIdLow
}

func (q *QueryId) SetMemberIdLow(memberIdLow uint64) {
	q.memberIdLow = memberIdLow
}

func (q *QueryId) MemberIdHigh() uint64 {
	return q.memberIdHigh
}

func (q *QueryId) SetMemberIdHigh(memberIdHigh uint64) {
	q.memberIdHigh = memberIdHigh
}

func NewQueryId(memberIdHigh uint64, memberIdLow uint64, localIdHigh uint64, localIdLow uint64) QueryId {
	return QueryId{memberIdHigh: memberIdHigh, memberIdLow: memberIdLow, localIdHigh: localIdHigh, localIdLow: localIdLow}
}
