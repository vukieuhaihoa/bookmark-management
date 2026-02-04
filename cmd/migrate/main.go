package main

import "github.com/vukieuhaihoa/bookmark-management/internal/infrastructure"

func main() {
	_ = infrastructure.CreateSQLDBAndMigration()
}
