package base

type IScannable interface {
	Scan(dest ...any) error
}
