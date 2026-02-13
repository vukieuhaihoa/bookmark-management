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
	err := script.BackfillForCodeShortenCol(dbClient)
	if err != nil {
		log.Error().Err(err).Msg("Failed to backfill code shorten column")
		common.HandlerError(err)
		return
	}

	log.Info().Msg("Clearing Redis cache after backfill.")
	redisClient := infrastructure.CreateRedisCon()
	err = redisClient.FlushAll(ctx).Err()
	if err != nil {
		log.Error().Err(err).Msg("Failed to clear Redis cache after backfill")
		common.HandlerError(err)
		return
	}
	log.Info().Msg("Redis cache cleared successfully.")
}
