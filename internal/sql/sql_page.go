package sql

type SqlPage struct {

	columnTypes []SqlColumnType
	data DataHolder
	last bool

}
