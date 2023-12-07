package queryTemplates

import "fmt"

var (
	CitiesQueryTemplate = &CommonFilterQueryTemplate{
		baseQuery: `SELECT v."location", count(*) AS cnt FROM hnh_data.vacancy v`,
		groupBy:   `GROUP BY v."location"`,
		orderBy:   `ORDER BY cnt`,
	}

	ExperienceQueryTemplate = &CommonFilterQueryTemplate{
		baseQuery: `SELECT v.experience, count(*) AS cnt FROM hnh_data.vacancy v`,
		groupBy:   `GROUP BY v.experience`,
	}

	EmploymentQueryTemplate = &CommonFilterQueryTemplate{
		baseQuery: `SELECT v.employment, count(*) AS cnt FROM hnh_data.vacancy v`,
		groupBy:   `GROUP BY v.employment`,
	}

	EducationTypeQueryTemplate = &CommonFilterQueryTemplate{
		baseQuery: `SELECT v.education_type, count(*) AS cnt FROM hnh_data.vacancy v`,
		groupBy:   `GROUP BY v.education_type`,
	}
)

type CommonFilterQueryTemplate struct {
	baseQuery   string
	whereClause string
	groupBy     string
	orderBy     string
}

func (qt *CommonFilterQueryTemplate) BuildQuery(searchCondition bool) string {
	if searchCondition {
		qt.whereClause = "WHERE plainto_tsquery($1) @@ v.fts"
	}
	return fmt.Sprintf("%s %s %s %s", qt.baseQuery, qt.whereClause, qt.groupBy, qt.orderBy)
}
