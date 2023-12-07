package queryTemplates

import (
	"fmt"
)

var (
	VacanciesSearchQueryTemplate = &SearchQueryTemplates{
		table_name: "hnh_data.vacancy",
	}

	CVsSearchQueryTemplate = &SearchQueryTemplates{
		table_name: "hnh_data.cv",
	}
)

type SearchQueryTemplates struct {
	table_name  string
	whereClause string
}

func (sqt *SearchQueryTemplates) BuildTemplate(searchCondition bool) string {
	if searchCondition {
		// sqt.whereClause = fmt.Sprintf("WHERE %s", searchTerm)
		sqt.whereClause = "WHERE plainto_tsquery($3) @@ tbl.fts"
	}

	return fmt.Sprintf(
		`WITH filtered_items AS (
		SELECT
			tbl.id,
			tbl.fts
		FROM
			%s tbl
		%s
	),
	count_total AS (
		SELECT
			COUNT(*) AS total
		FROM
			filtered_items
	)
	SELECT
		fi.id,
		ct.total
	FROM
		filtered_items fi,
		count_total ct
	LIMIT $1
	OFFSET $2`,
		sqt.table_name,
		sqt.whereClause,
	)
}
