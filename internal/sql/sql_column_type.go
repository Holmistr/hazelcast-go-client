package sql

type SqlColumnType int

const (

	VARCHAR SqlColumnType = 0
	BOOLEAN SqlColumnType = 1
	TINYINT SqlColumnType = 2
	SMALLINT SqlColumnType = 3
	INTEGER SqlColumnType = 4
	BIGINT SqlColumnType = 5
	DECIMAL SqlColumnType = 6
	REAL SqlColumnType = 7
	DOUBLE SqlColumnType = 8
	DATE SqlColumnType = 9
	TIME SqlColumnType = 10
	TIMESTAMP SqlColumnType = 11
	TIMESTAMP_WITH_TIME_ZONE SqlColumnType = 12
	OBJECT SqlColumnType = 13
	NULL SqlColumnType = 14

)