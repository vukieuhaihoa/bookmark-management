package model

const (
	RedisShortenPrefix    = "r" //  this prefix is used for Redis keys related to bookmark shorten codes
	BookmarkShortenPrefix = "p" //	this prefix is used for the CodeShortenEncoded field in the Bookmark model to distinguish it from other types of codes
)
