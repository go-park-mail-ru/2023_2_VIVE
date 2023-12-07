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

	SalaryQueryTemplate = &CommonFilterQueryTemplate{
		baseQuery: `SELECT
						COALESCE(MIN(v.salary_lower_bound), 0) AS min_salary ,
						COALESCE(MAX(v.salary_upper_bound), 0) AS max_salary,
						COUNT(*)
					FROM
						hnh_data.vacancy v`,
	}
)

type CommonFilterQueryTemplate struct {
	baseQuery   string
	whereClause string
	groupBy     string
	orderBy     string
}

func (cfqt *CommonFilterQueryTemplate) BuildQuery(searchCondition bool) string {
	if searchCondition {
		cfqt.whereClause = "WHERE plainto_tsquery($1) @@ v.fts"
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
