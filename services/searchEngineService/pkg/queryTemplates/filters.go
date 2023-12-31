package queryTemplates

import "fmt"

var (
	// Vacancy filters
	VacCitiesQueryTemplate = &CommonFilterQueryTemplate{
		baseQuery:   `SELECT v."location", count(*) AS cnt FROM hnh_data.vacancy v`,
		whereClause: `WHERE plainto_tsquery($1) @@ v.fts`,
		groupBy:     `GROUP BY v."location"`,
		orderBy:     `ORDER BY cnt`,
	}

	VacExperienceQueryTemplate = &CommonFilterQueryTemplate{
		baseQuery:   `SELECT v.experience, count(*) AS cnt FROM hnh_data.vacancy v`,
		whereClause: `WHERE plainto_tsquery($1) @@ v.fts`,
		groupBy:     `GROUP BY v.experience`,
	}

	VacEmploymentQueryTemplate = &CommonFilterQueryTemplate{
		baseQuery:   `SELECT v.employment, count(*) AS cnt FROM hnh_data.vacancy v`,
		whereClause: `WHERE plainto_tsquery($1) @@ v.fts`,
		groupBy:     `GROUP BY v.employment`,
	}

	VacEducationTypeQueryTemplate = &CommonFilterQueryTemplate{
		baseQuery:   `SELECT v.education_type, count(*) AS cnt FROM hnh_data.vacancy v`,
		whereClause: `WHERE plainto_tsquery($1) @@ v.fts`,
		groupBy:     `GROUP BY v.education_type`,
	}

	VacSalaryQueryTemplate = &CommonFilterQueryTemplate{
		baseQuery: `SELECT
						COALESCE(MIN(v.salary_lower_bound), 0) AS min_salary ,
						COALESCE(MAX(v.salary_upper_bound), 0) AS max_salary,
						COUNT(*)
					FROM
						hnh_data.vacancy v`,
		whereClause: `WHERE plainto_tsquery($1) @@ v.fts`,
	}

	// CV filters
	CvCitiesQueryTemplate = &CommonFilterQueryTemplate{
		baseQuery:   `SELECT c."location", count(*) AS cnt FROM hnh_data.cv c`,
		whereClause: `WHERE plainto_tsquery($1) @@ c.fts`,
		groupBy:     `GROUP BY c."location"`,
		orderBy:     `ORDER BY cnt`,
	}

	CvEducationTypeQueryTemplate = &CommonFilterQueryTemplate{
		baseQuery:   `SELECT c.education_level, count(*) AS cnt FROM hnh_data.cv c`,
		whereClause: `WHERE plainto_tsquery($1) @@ c.fts`,
		groupBy:     `GROUP BY c.education_level`,
		orderBy:     `ORDER BY cnt`,
	}

	CvGenderQueryTemplate = &CommonFilterQueryTemplate{
		baseQuery:   `SELECT c.gender, count(*) AS cnt FROM hnh_data.cv c`,
		whereClause: `WHERE plainto_tsquery($1) @@ c.fts`,
		groupBy:     `GROUP BY c.gender`,
		orderBy:     `ORDER BY cnt`,
	}

	CvExperienceQueryTemplate = &CommonFilterQueryTemplate{
		baseQuery:   `SELECT c.experience, count(*) AS cnt FROM hnh_data.cv c`,
		whereClause: `WHERE plainto_tsquery($1) @@ c.fts`,
		groupBy:     `GROUP BY c.experience`,
	}
)

type CommonFilterQueryTemplate struct {
	baseQuery   string
	whereClause string
	groupBy     string
	orderBy     string
}

func (cfqt CommonFilterQueryTemplate) BuildQuery(searchCondition bool) string {
	cfqt.baseQuery = fmt.Sprintf(cfqt.baseQuery)
	if !searchCondition {
		cfqt.whereClause = ""
	}
	return fmt.Sprintf("%s %s %s %s", cfqt.baseQuery, cfqt.whereClause, cfqt.groupBy, cfqt.orderBy)
}

// type SalaryFilterQueryTemplate struct {
// 	whereClause string
// }

// func (sfqt *SalaryFilterQueryTemplate) BuildQuery(searchCondition bool) string {
// 	if searchCondition {
// 		sfqt.whereClause = "WHERE plainto_tsquery($1) @@ v.fts"
// 	}
// 	return fmt.Sprintf(
// 		`WITH min_max_salary AS (
// 		SELECT
// 			MIN(v.salary_lower_bound) AS min_salary,
// 			MAX(v.salary_upper_bound) AS max_salary
// 		FROM
// 			hnh_data.vacancy v
// 		%s
// 	), salary_range AS (
// 		SELECT
// 			min_salary,
// 			max_salary,
// 			(max_salary - min_salary) / 5 AS range_size
// 		FROM
// 			min_max_salary
// 	), salary_buckets AS (
// 		SELECT
// 			min_salary + range_size * n AS range_start,
// 			min_salary + range_size * (n + 1) AS range_end
// 		FROM
// 			salary_range,
// 			generate_series(0, 4) AS n
// 	)
// 	SELECT
// 		sb.range_start,
// 		sb.range_end,
// 		COUNT(v.*) AS count
// 	FROM
// 		salary_buckets sb
// 	LEFT JOIN hnh_data.vacancy v ON v.salary_lower_bound >= sb.range_start AND v.salary_lower_bound < sb.range_end
// 	GROUP BY
// 		sb.range_start, sb.range_end
// 	ORDER BY
// 		sb.range_start`,
// 		sfqt.whereClause,
// 	)
// }
