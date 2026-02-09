package main

import (
	"github.com/vukieuhaihoa/bookmark-management/cmd/migrate/script"
	"github.com/vukieuhaihoa/bookmark-management/internal/infrastructure"
)

func main() {
	dbClient := infrastructure.CreateSQLDBAndMigration()
	script.BackfillForCodeShortenCol(dbClient)
}
