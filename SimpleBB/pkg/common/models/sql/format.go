package sql

func BoolToSql(b bool) bool {
	return b
}

func ByteToSql(b byte) int64 {
	return int64(b)
}

func IntToSql(i int) int64 {
	return int64(i)
}

func StringToSql(s string) string {
	return s
}
