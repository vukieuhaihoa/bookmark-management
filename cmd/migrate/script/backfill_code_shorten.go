package script

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/vukieuhaihoa/bookmark-management/pkg/encoding"
	"gorm.io/gorm"
)

type bookmark struct {
	ID string `gorm:"column:id"`
}

func (bookmark) TableName() string {
	return "bookmarks"
}

type schema_migrations struct {
	Version int64 `gorm:"column:version"`
}

func (schema_migrations) TableName() string {
	return "schema_migrations"
}

func BackfillForCodeShortenCol(db *gorm.DB) error {
	sm := &schema_migrations{}
	err := db.Model(&schema_migrations{}).Where("version = ?", 3).First(sm).Error
	if err != nil {
		log.Error().Err(err).Msg("Failed to verify migration version for backfilling code_shorten column.")
		return err
	}

	log.Info().Msg("Starting backfill for code_shorten column in bookmarks table.")
	var bookmarks []bookmark
	if err := db.Model(&bookmark{}).Find(&bookmarks).Error; err != nil {
		log.Error().Err(err).Msg("Failed to fetch bookmarks for backfilling code_shorten column.")
		return err
	}

	if len(bookmarks) == 0 {
		log.Info().Msg("No bookmarks found for backfilling code_shorten column.")
		return nil
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		for i, bm := range bookmarks {
			seq := int64(i + 1)
			encoded, err := encoding.StdEncoding.EncodeInt64ToString(seq)
			if err != nil {
				return fmt.Errorf("failed to encode %d: %w", seq, err)
			}

			if err := tx.Model(&bookmark{}).
				Where("id = ?", bm.ID).
				Updates(map[string]interface{}{
					"code_shorten":         seq,
					"code_shorten_encoded": "p_" + encoded,
				}).Error; err != nil {
				return fmt.Errorf("failed to update bookmark %s: %w", bm.ID, err)
			}

			fmt.Printf("updated %s: code_shorten=%d, code_shorten_encoded=%s\n", bm.ID, seq, "p_"+encoded)
		}

		// Sync sequence so next auto-increment continues after max
		if err := tx.Exec(
			"SELECT setval(pg_get_serial_sequence('bookmarks', 'code_shorten'), ?)",
			int64(len(bookmarks)),
		).Error; err != nil {
			return fmt.Errorf("failed to reset sequence: %w", err)
		}

		return nil
	})

	if err != nil {
		log.Error().Err(err).Msg("Transaction failed while backfilling code_shorten column.")
		return err
	}

	log.Info().Msg("Successfully backfilled code_shorten column for all bookmarks.")
	return nil
}
