package svc

import "fmt"

type DB struct {
	config map[string]string
}

func NewDB(config map[string]string) *DB {
	inst := new(DB)
	inst.config = config

	return inst
}

func (db *DB) Connect() {
	fmt.Println("  - Database connected")
}