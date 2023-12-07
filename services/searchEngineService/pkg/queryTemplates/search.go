package queryTemplates

import (
	"HnH/services/searchEngineService/pkg/searchOptions"
	"HnH/services/searchEngineService/searchEnginePB"
	"fmt"
	"strings"
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

func (sqt SearchQueryTemplates) BuildTemplate(limit, offset int64, options *searchEnginePB.SearchOptions) (string, []interface{}) {
	// searchTerm := ""
	// if searchCondition {
	// 	// sqt.whereClause = fmt.Sprintf("WHERE %s", searchTerm)
	// 	searchTerm = "plainto_tsquery($3) @@ tbl.fts"
	// }

	var initArgs []interface{}
	initArgs = append(initArgs, limit, offset)
	whereOptions, args := searchOptions.GetWhereTerms(options, 3)
	if len(whereOptions) > 0 {
		sqt.whereClause = fmt.Sprintf("WHERE %s", strings.Join(whereOptions, " AND "))
	}
	fmt.Printf("where clause: %s\n", sqt.whereClause)
	initArgs = append(initArgs, args...)

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
	), initArgs
}
