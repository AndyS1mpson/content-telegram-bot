package sql

import (
	"fmt"

	"content-telegram-bot/internal/utils/slices"

	sq "github.com/Masterminds/squirrel"
)

// PageQuery параметры запроса для добавления пагинации
type PageQuery struct {
	OffsetQuery *OffsetPageQuery
	LastIDQuery *LastIDPageQuery
	OrderBy     OrderBy
}

// OffsetPageQuery параметры запроса на получение постраничной информации с помощью Limit + Offset
type OffsetPageQuery struct {
	Page     int64
	PageSize int64
}

// LastIDPageQuery параметры запроса на получение постраничной информации с помощью Limit + последний полученный ID
type LastIDPageQuery struct {
	LastRecievedID int64
	PageSize       int64
}

// SimplePage возвращает актуальный лимит и offset
func SimplePage(p PageQuery) (limit, offset int64) {
	return p.OffsetQuery.PageSize, (p.OffsetQuery.Page - 1) * p.OffsetQuery.PageSize
}

func AddPaging(sql sq.SelectBuilder, paging PageQuery, idColumn string) sq.SelectBuilder {
	if paging.OffsetQuery != nil {
		limit, offset := SimplePage(paging)
		sql = sql.Limit(uint64(limit)).Offset(uint64(offset))
	} else if paging.LastIDQuery != nil {
		sql = sql.Where(sq.Gt{idColumn: paging.LastIDQuery.LastRecievedID}).Limit(uint64(paging.LastIDQuery.PageSize))
	}

	return sql.OrderBy(
		slices.Map(paging.OrderBy.Columns, func(c OrderByColumn) string {
			return fmt.Sprintf("%s %s", c.Column, c.Direction)
		})...,
	)
}
