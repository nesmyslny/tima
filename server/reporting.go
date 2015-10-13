package server

import (
	"net/http"
	"time"

	"gopkg.in/gorp.v1"
)

type DateValue struct {
	Date  time.Time `db:"date" json:"date"`
	Value float32   `db:"value" json:"value"`
}

type Chart struct {
	Series []string    `json:"series"`
	Labels []string    `json:"labels"`
	Data   [][]float32 `json:"data"`
}

type ReportCriteria struct {
	StartDate *time.Time `json:"startDate"`
	EndDate   *time.Time `json:"endDate"`
	Projects  []Project  `json:"projects"`
}

type ReportOverview struct {
	DurationTotal     int   `db:"duration_total" json:"durationTotal"`
	ProjectCount      int   `db:"project_count" json:"projectCount"`
	ActivityTypeCount int   `db:"activity_type_count" json:"activityTypeCount"`
	UserCount         int   `db:"user_count" json:"userCount"`
	DepartmentCount   int   `db:"department_count" json:"departmentCount"`
	Timeline          Chart `db:"-" json:"timeline"`
}

type ReportProjects struct {
	Timeline Chart `db:"-" json:"timeline"`
}

type Reporting struct {
	dbMap *gorp.DbMap
}

func NewReporting(db *DB) *Reporting {
	return &Reporting{db.dbMap}
}

func (reporting *Reporting) CreateReportOverview(context *HandlerContext) (interface{}, *HandlerError) {
	var criteria ReportCriteria
	err := context.GetReqBodyJSON(&criteria)
	if err != nil {
		return nil, &HandlerError{err, "couldn't retrieve report overview", http.StatusInternalServerError}
	}

	sqlCriteria, paramsCriteria := reporting.createSqlWhere(&criteria)

	sql := "select count(distinct a.user_id) user_count, count(distinct a.project_id) project_count, " +
		"count(distinct a.activity_type_id) activity_type_count, ifnull(sum(a.duration), 0) duration_total, " +
		"count(distinct u.department_id) department_count from activity a, user u where a.user_id = u.id"
	if len(sqlCriteria) > 0 {
		sql += " and " + sqlCriteria
	}

	var overview *ReportOverview
	err = reporting.dbMap.SelectOne(&overview, sql, paramsCriteria...)
	if err != nil {
		return nil, &HandlerError{err, "couldn't retrieve report overview", http.StatusInternalServerError}
	}

	sqlTimeline := "select day date, round(sum(a.duration) / 60, 2) value from activity a"
	if len(sqlCriteria) > 0 {
		sqlTimeline += " where " + sqlCriteria
	}
	sqlTimeline += " group by day order by day"

	var dateValues []DateValue
	_, err = reporting.dbMap.Select(&dateValues, sqlTimeline, paramsCriteria...)
	if err != nil {
		return nil, &HandlerError{err, "couldn't retrieve timeline in report overview", http.StatusInternalServerError}
	}

	overview.Timeline.Labels, overview.Timeline.Data = reporting.createLabelsAndData(criteria.StartDate, criteria.EndDate, [][]DateValue{dateValues})
	return overview, nil
}

func (reporting *Reporting) CreateReportProjects(context *HandlerContext) (interface{}, *HandlerError) {
	var criteria ReportCriteria
	err := context.GetReqBodyJSON(&criteria)
	if err != nil {
		return nil, &HandlerError{err, "couln't retieve report of projects", http.StatusInternalServerError}
	}

	sqlCriteria, paramsCriteria := reporting.createSqlWhere(&criteria)

	sql := "select * from project where id in (select distinct a.project_id from activity a"
	if len(sqlCriteria) > 0 {
		sql += " where " + sqlCriteria
	}
	sql += ")"

	var projects []Project
	_, err = reporting.dbMap.Select(&projects, sql, paramsCriteria...)
	if err != nil {
		return nil, &HandlerError{err, "couldn't retrieve projects", http.StatusInternalServerError}
	}

	projectsView := &ReportProjects{}
	projectsView.Timeline.Series = make([]string, len(projects))
	projectsDateValues := make([][]DateValue, len(projects))

	for i, project := range projects {
		projectsView.Timeline.Series[i] = project.Title

		sqlTimeline := "select day date, round(sum(a.duration) / 60, 2) value from activity a where project_id = ?"
		if len(sqlCriteria) > 0 {
			sqlTimeline += " and " + sqlCriteria
		}
		sqlTimeline += " group by day order by day"

		var dateValues []DateValue
		_, err = reporting.dbMap.Select(&dateValues, sqlTimeline, append([]interface{}{project.ID}, paramsCriteria...)...)
		if err != nil {
			return nil, &HandlerError{err, "couldn't retrieve timeline in report of projects", http.StatusInternalServerError}
		}

		projectsDateValues[i] = dateValues
	}

	projectsView.Timeline.Labels, projectsView.Timeline.Data = reporting.createLabelsAndData(criteria.StartDate, criteria.EndDate, projectsDateValues)
	return projectsView, nil
}

func (reporting *Reporting) convertDateValues(data []DateValue) (labels []string, values []float32) {
	labels = make([]string, len(data))
	values = make([]float32, len(data))

	for i, elem := range data {
		labels[i] = elem.Date.Format(dateLayout)
		values[i] = elem.Value
	}

	return labels, values
}

func (reporting *Reporting) insertMissingDates(start time.Time, end time.Time, data *[]DateValue) {
	d := start

	for i := 0; d.Sub(end) <= 0; i++ {
		var elem *time.Time = nil
		if i < len(*data) {
			elem = &(*data)[i].Date
		}

		if elem == nil || !elem.Equal(d) {
			*data = append((*data)[:i], append([]DateValue{{d, 0}}, (*data)[i:]...)...)
		}

		d = d.Add(time.Hour * 24)
	}
}

func (reporting *Reporting) getDateBoundaries(criteriaStart *time.Time, criteriaEnd *time.Time, data [][]DateValue) (start time.Time, end time.Time) {
	for _, elem := range data {
		if len(elem) > 0 {
			if start.IsZero() || elem[0].Date.Sub(start) < 0 {
				start = elem[0].Date
			}

			if end.IsZero() || elem[len(elem)-1].Date.Sub(end) > 0 {
				end = elem[len(elem)-1].Date
			}
		}
	}

	if criteriaStart != nil {
		start = *criteriaStart
	}
	if criteriaEnd != nil {
		end = *criteriaEnd
	}

	return start, end
}

func (reporting *Reporting) createLabelsAndData(criteriaStart *time.Time, criteriaEnd *time.Time, dateValuesArray [][]DateValue) (labels []string, data [][]float32) {
	data = make([][]float32, len(dateValuesArray))
	start, end := reporting.getDateBoundaries(criteriaStart, criteriaEnd, dateValuesArray)

	for i, dateValues := range dateValuesArray {
		reporting.insertMissingDates(start, end, &dateValues)
		labels, data[i] = reporting.convertDateValues(dateValues)
	}

	return labels, data
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
