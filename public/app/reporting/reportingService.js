angular.module('tima').factory('reportingService',
['_', '$moment', 'Reporting',
function(_, $moment, Reporting) {

    function prepareCriteria(criteria) {
        if (criteria.startDate) {
            criteria.startDate = $moment(criteria.startDate).format("YYYY-MM-DD[T]00:00:00.000[Z]");
        }
        if (criteria.endDate) {
            criteria.endDate = $moment(criteria.endDate).format("YYYY-MM-DD[T]00:00:00.0000[Z]");
        }
    }

    function addTimelineGrouping(obj, timeline) {
        obj.timelineDay = {
            labels: timeline.labels,
            values: timeline.data
        };

        obj.timelineWeek = groupTimeline(timeline, "YYYY[W]WW");
        obj.timelineMonth = groupTimeline(timeline, "YYYY-MM");

        if (timeline.labels.length > 180) {
            return obj.timelineMonth;
        } else if (timeline.labels.length > 31) {
            return obj.timelineWeek;
        }
        return obj.timelineDay;
    }

    function groupTimeline(timeline, groupByFormat) {
        var groupedTimeline = {};

        for (var i = 0; i < timeline.labels.length; i++) {
            var group = $moment(timeline.labels[i]).format(groupByFormat);

            if (_.isUndefined(groupedTimeline[group])) {
                groupedTimeline[group] = [];
            }

            for (var j = 0; j < timeline.data.length; j++) {
                if (_.isUndefined(groupedTimeline[group][j])) {
                    groupedTimeline[group][j] = 0;
                }
                groupedTimeline[group][j] += timeline.data[j][i];
            }
        }

        return {
            labels: _.keys(groupedTimeline),
            values: _.unzip(_.values(groupedTimeline))
        };
    }

    var service = {

        getReportOverview: function(criteria, callback) {

            prepareCriteria(criteria);

            return Reporting.createOverview(criteria, function(overview) {
                overview.chartType = "Line";
                var duration = overview.durationTotal;

                if (duration < 1440) {
                    var hours = $moment.duration(duration, "minutes").asHours();
                    overview.durationHours = +hours.toFixed(2);
                } else {
                    var days = $moment.duration(duration, "minutes").asDays();
                    overview.durationDays = +days.toFixed(2);
                }

                overview.currentTimeline = addTimelineGrouping(overview, overview.timeline);
            });
        },

        getReportProjects: function(criteria, callback) {

            prepareCriteria(criteria);

            return Reporting.createProjects(criteria, function(projectsView) {
                projectsView.chartType = "Line";
                projectsView.noData = projectsView.timeline.labels.length === 0;
                projectsView.currentTimeline = addTimelineGrouping(projectsView, projectsView.timeline);


                projectsView.pie = {
                    labels: projectsView.timeline.series,
                    data: []
                };

                for (var i = 0; i < projectsView.timeline.series.length; i++) {
                    projectsView.pie.data.push(_.sum(projectsView.timeline.data[i]));
                }
            });
        },

        createCriteria: function() {
            return {
                startDate: $moment().startOf("month"),
                endDate: $moment().endOf("month"),
                projects: []
            };
        }
    };

    return service;
}]);
