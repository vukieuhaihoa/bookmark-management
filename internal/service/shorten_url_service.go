package service

import (
	"context"

	"github.com/vukieuhaihoa/bookmark-management/internal/repository"
	"github.com/vukieuhaihoa/bookmark-management/pkg/stringutils"
)

const (
	defaultURLCodeLength = 8
)

type ShortenURL interface {
	ShortenURL(ctx context.Context, originalURL string) (string, error)
}

type shortenURL struct {
	repo repository.UrlStorage
}

func NewShortenURL(repo repository.UrlStorage) ShortenURL {
	return &shortenURL{
		repo: repo,
	}
}

func (s *shortenURL) ShortenURL(ctx context.Context, originalURL string) (string, error) {
	urlCode, err := stringutils.GenerateCode(defaultURLCodeLength)
	if err != nil {
		return "", err
	}

	err = s.repo.StoreURL(ctx, urlCode, originalURL, 0)
	if err != nil {
		return "", err
	}

	return urlCode, nil
}
