package server

import (
	"net/http"
	"time"

	"gopkg.in/gorp.v1"
)

type ChartData struct {
	Date  string  `db:"date" json:"date"`
	Value float32 `db:"value" json:"value"`
}

type ReportCriteria struct {
	StartDate *time.Time `json:"startDate"`
	EndDate   *time.Time `json:"endDate"`
	Projects  []Project  `json:"projects"`
}

type ReportOverview struct {
	DurationTotal     int         `db:"duration_total" json:"durationTotal"`
	ProjectCount      int         `db:"project_count" json:"projectCount"`
	ActivityTypeCount int         `db:"activity_type_count" json:"activityTypeCount"`
	UserCount         int         `db:"user_count" json:"userCount"`
	DepartmentCount   int         `db:"department_count" json:"departmentCount"`
	Timeline          []ChartData `db:"-" json:"timeline"`
}

type Reporting struct {
	dbMap *gorp.DbMap
}

func NewReporting(db *DB) *Reporting {
	return &Reporting{db.dbMap}
}

func (reporting *Reporting) CreateOverview(context *HandlerContext) (interface{}, *HandlerError) {
	var criteria ReportCriteria
	err := context.GetReqBodyJSON(&criteria)
	if err != nil {
		return nil, &HandlerError{err, "couldn't retrieve report overview", http.StatusInternalServerError}
	}

	sqlCriteria, paramsCriteria := reporting.createSqlWhere(&criteria)

	sql := "select count(distinct a.user_id) user_count, count(distinct a.project_id) project_count, " +
		"count(distinct a.activity_type_id) activity_type_count, sum(a.duration) duration_total, " +
		"count(distinct u.department_id) department_count from activity a, user u where a.user_id = u.id"
	if len(sqlCriteria) > 0 {
		sql += " and " + sqlCriteria
	}

	var overview *ReportOverview
	err = reporting.dbMap.SelectOne(&overview, sql, paramsCriteria...)
	if err != nil {
		return nil, &HandlerError{err, "couldn't retrieve report overview", http.StatusInternalServerError}
	}

	sqlTimeline := "select date_format(day, '%Y-%m-%d') date, round(sum(a.duration) / 60, 2) value from activity a"
	if len(sqlCriteria) > 0 {
		sqlTimeline += " where " + sqlCriteria
	}
	sqlTimeline += " group by day order by date"

	var timeline []ChartData
	_, err = reporting.dbMap.Select(&timeline, sqlTimeline, paramsCriteria...)
	if err != nil {
		return nil, &HandlerError{err, "couldn't retrieve timeline in report overview", http.StatusInternalServerError}
	}
	overview.Timeline = timeline

	return overview, nil
}

func (reporting *Reporting) createSqlWhere(criteria *ReportCriteria) (string, []interface{}) {
	var whereConditions []string
	var params []interface{}

	if criteria.StartDate != nil {
		whereConditions = append(whereConditions, "a.day >= ?")
		params = append(params, criteria.StartDate)
	}

	if criteria.EndDate != nil {
		whereConditions = append(whereConditions, "a.day <= ?")
		params = append(params, criteria.EndDate)
	}

	if len(criteria.Projects) > 0 {
		var projectsIn string
		for _, p := range criteria.Projects {
			if len(projectsIn) > 0 {
				projectsIn += ", "
			}
			projectsIn += "?"
			params = append(params, p.ID)
		}
		whereConditions = append(whereConditions, "a.project_id in("+projectsIn+")")
	}

	var where string
	for _, v := range whereConditions {
		if len(where) > 0 {
			where += " and "
		}
		where += v
	}

	return where, params
}
