package sql

type SqlColumnMetadata struct {

	name string
	columnType SqlColumnType
	nullable bool

}

func NewSqlColumnMetadata(name string, columnType SqlColumnType, nullable bool) SqlColumnMetadata {
	return SqlColumnMetadata{name: name, columnType: columnType, nullable: nullable}
}

func (s *SqlColumnMetadata) Nullable() bool {
	return s.nullable
}

func (s *SqlColumnMetadata) SetNullable(nullable bool) {
	s.nullable = nullable
}

func (s *SqlColumnMetadata) ColumnType() SqlColumnType {
	return s.columnType
}

func (s *SqlColumnMetadata) SetColumnType(columnType SqlColumnType) {
	s.columnType = columnType
}

func (s *SqlColumnMetadata) Name() string {
	return s.name
}

func (s *SqlColumnMetadata) SetName(name string) {
	s.name = name
}
