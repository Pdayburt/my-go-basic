package domain

import "time"

const (
	ArticleStatusUnknown ArticleStatus = iota
	ArticleStatusUnpublished
	ArticleStatusPublished
	ArticleStatusPrivate
)

type Article struct {
	Id      int64
	Title   string
	Content string
	//Author从用户来
	Author Author
	Status ArticleStatus
	Ctime  time.Time
	Utime  time.Time
}

func (a *Article) Abstract() string {
	//摘要取 前几句
	cs := []rune(a.Content)
	if len(cs) < 100 {
		return a.Content
	}
	return string(cs[:100])
}

type Author struct {
	Id   int64
	Name string
}

type ArticleStatus uint8

func (as ArticleStatus) ToUint8() uint8 {
	return uint8(as)
}

func (as ArticleStatus) NonPublish() bool {
	return as != ArticleStatusUnpublished
}

func (as ArticleStatus) String() string {
	switch as {
	case ArticleStatusUnpublished:
		return "Unpublished"
	case ArticleStatusPublished:
		return "Published"
	case ArticleStatusPrivate:
		return "Private"
	default:
		return "Unknown"
	}

}
