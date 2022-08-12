package models

type ScannedRow interface {
	Scan(dest ...any) error
}
