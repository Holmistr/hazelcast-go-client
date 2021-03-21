package sql

type QueryId struct {

	memberIdHigh int64
	memberIdLow int64
	localIdHigh int64
	localIdLow int64

}

func (q *QueryId) LocalIdLow() int64 {
	return q.localIdLow
}

func (q *QueryId) SetLocalIdLow(localIdLow int64) {
	q.localIdLow = localIdLow
}

func (q *QueryId) LocalIdHigh() int64 {
	return q.localIdHigh
}

func (q *QueryId) SetLocalIdHigh(localIdHigh int64) {
	q.localIdHigh = localIdHigh
}

func (q *QueryId) MemberIdLow() int64 {
	return q.memberIdLow
}

func (q *QueryId) SetMemberIdLow(memberIdLow int64) {
	q.memberIdLow = memberIdLow
}

func (q *QueryId) MemberIdHigh() int64 {
	return q.memberIdHigh
}

func (q *QueryId) SetMemberIdHigh(memberIdHigh int64) {
	q.memberIdHigh = memberIdHigh
}

func NewQueryId(memberIdHigh int64, memberIdLow int64, localIdHigh int64, localIdLow int64) QueryId {
	return QueryId{memberIdHigh: memberIdHigh, memberIdLow: memberIdLow, localIdHigh: localIdHigh, localIdLow: localIdLow}
}
