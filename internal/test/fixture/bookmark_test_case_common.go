package fixture

import (
	"time"

	"github.com/vukieuhaihoa/bookmark-management/internal/app/model"
	"gorm.io/gorm"
)

type BookmarkCommonTestDB struct {
	base
}

// Migrate migrates the database schema for the UserCommonTestDB fixture.
//
// Returns:
//   - error: An error if migration fails, otherwise nil
func (b *BookmarkCommonTestDB) Migrate() error {
	return b.db.AutoMigrate(&model.User{}, &model.Bookmark{})
}

// GenerateData populates the test database with common user test data.
//
// Returns:
//   - error: An error if data generation fails, otherwise nil
func (b *BookmarkCommonTestDB) GenerateData() error {
	db := b.db.Session(&gorm.Session{})

	users := []*model.User{
		{
			Base: model.Base{
				ID:        "de305d54-75b4-431b-adb2-eb6b9e546000",
				CreatedAt: TestTime,
				UpdatedAt: TestTime,
			},
			DisplayName: "Alice",
			Username:    "Alice",
			Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
			Email:       "alice@example.com",
		},
		{
			Base: model.Base{
				ID:        "123e4567-e89b-12d3-a456-eb6b9e546001",
				CreatedAt: TestTime,
				UpdatedAt: TestTime,
			},
			DisplayName: "Bob",
			Username:    "Bob",
			Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
			Email:       "bob@example.com",
		},
		{
			Base: model.Base{
				ID:        "987e6543-e21b-12d3-a456-eb6b9e546002",
				CreatedAt: TestTime,
				UpdatedAt: TestTime,
			},
			DisplayName: "Charlie",
			Username:    "Charlie",
			Password:    "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6rS8nY7b1p6K5j5p6v5Q5Z5Z5Z5e",
			Email:       "charlie@example.com",
		},
		{
			Base: model.Base{
				ID:        "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
				CreatedAt: TestTime,
				UpdatedAt: TestTime,
			},
			DisplayName: "Test User 1",
			Username:    "testuser001",
			Password:    "$2a$10$hhuB9rZrp5ikmRb5yAF9hev6AE2tC404jhtP.bdOjme9lECJClzFu",
			Email:       "testuser001@example.com",
		},
	}

	err := db.CreateInBatches(users, 10).Error
	if err != nil {
		return err
	}

	bookmarks := []*model.Bookmark{
		{
			Base: model.Base{
				ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000001",
				CreatedAt: TestTime,
				UpdatedAt: TestTime,
			},
			UserID:             "de305d54-75b4-431b-adb2-eb6b9e546000",
			URL:                "https://example.com/alice",
			Description:        "Bookmark for Alice",
			Code:               "AAAA000001",
			CodeShorten:        1,
			CodeShortenEncoded: "p_1",
		},
		{
			Base: model.Base{
				ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000002",
				CreatedAt: TestTime,
				UpdatedAt: TestTime,
			},
			UserID:             "123e4567-e89b-12d3-a456-eb6b9e546001",
			URL:                "https://example.com/bob",
			Description:        "Bookmark for Bob",
			Code:               "AAAA000002",
			CodeShorten:        2,
			CodeShortenEncoded: "p_2",
		},
		{
			Base: model.Base{
				ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000003",
				CreatedAt: TestTime,
				UpdatedAt: TestTime,
			},
			UserID:             "987e6543-e21b-12d3-a456-eb6b9e546002",
			URL:                "https://example.com/charlie",
			Description:        "Bookmark for Charlie",
			Code:               "AAAA000003",
			CodeShorten:        3,
			CodeShortenEncoded: "p_3",
		},
		{
			Base: model.Base{
				ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000004",
				CreatedAt: TestTime,
				UpdatedAt: TestTime,
			},
			UserID:             "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			URL:                "https://example.com/testuser001",
			Description:        "Bookmark for Test User 1 - record 1",
			Code:               "AAAA000004",
			CodeShorten:        4,
			CodeShortenEncoded: "p_4",
		},
		{
			Base: model.Base{
				ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000005",
				CreatedAt: TestTime,
				UpdatedAt: TestTime,
			},
			UserID:             "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			URL:                "https://example.com/testuser001",
			Description:        "Bookmark for Test User 1 - record 2",
			Code:               "AAAA000005",
			CodeShorten:        5,
			CodeShortenEncoded: "p_5",
		},
		{
			Base: model.Base{
				ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000006",
				CreatedAt: TestTime.Add(2 * time.Hour),
				UpdatedAt: TestTime.Add(2 * time.Hour),
			},
			UserID:             "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			URL:                "https://example.com/testuser003",
			Description:        "Bookmark for Test User 1 - record 3",
			Code:               "AAAA000006",
			CodeShorten:        6,
			CodeShortenEncoded: "p_6",
		},
		{
			Base: model.Base{
				ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000007",
				CreatedAt: TestTime.Add(3 * time.Hour),
				UpdatedAt: TestTime.Add(3 * time.Hour),
			},
			UserID:             "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			URL:                "https://golang.dev/blog/clean-architecture",
			Description:        "Go backend best practices article",
			Code:               "AAAA000007",
			CodeShorten:        7,
			CodeShortenEncoded: "p_7",
		},
		{
			Base: model.Base{
				ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000008",
				CreatedAt: TestTime.Add(4 * time.Hour),
				UpdatedAt: TestTime.Add(4 * time.Hour),
			},
			UserID:             "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			URL:                "https://db-tutorials.dev/postgresql-indexing",
			Description:        "Learn PostgreSQL indexing basics",
			Code:               "AAAA000008",
			CodeShorten:        8,
			CodeShortenEncoded: "p_8",
		},
		{
			Base: model.Base{
				ID:        "a1b2c3d4-e5f6-7890-abcd-ef0000000009",
				CreatedAt: TestTime.Add(5 * time.Hour),
				UpdatedAt: TestTime.Add(5 * time.Hour),
			},
			UserID:             "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
			URL:                "https://redis.io/docs/manual/data-types/",
			Description:        "Redis data types documentation",
			Code:               "AAAA000009",
			CodeShorten:        9,
			CodeShortenEncoded: "p_9",
		},
	}

	return db.CreateInBatches(bookmarks, 10).Error
}
