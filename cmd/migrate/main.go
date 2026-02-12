package main

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/vukieuhaihoa/bookmark-management/cmd/migrate/script"
	"github.com/vukieuhaihoa/bookmark-management/internal/infrastructure"
	"github.com/vukieuhaihoa/bookmark-management/pkg/common"
)

func main() {
	ctx := context.Background()
	dbClient := infrastructure.CreateSQLDBAndMigration()
	script.BackfillForCodeShortenCol(dbClient)

	log.Info().Msg("Clearing Redis cache after backfill.")
	redisClient := infrastructure.CreateRedisCon()
	err := redisClient.FlushAll(ctx).Err()
	common.HandlerError(err)
	log.Info().Msg("Redis cache cleared successfully.")
}
