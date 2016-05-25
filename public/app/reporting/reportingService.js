angular.module('tima').factory('reportingService',
['_', '$moment', 'Reporting',
function(_, $moment, Reporting) {

    function prepareCriteria(criteria) {
        var c = _.clone(criteria);
        if (c.startDate) {
            c.startDate = $moment(c.startDate).format("YYYY-MM-DD[T]00:00:00.000[Z]");
        }
        if (c.endDate) {
            c.endDate = $moment(c.endDate).format("YYYY-MM-DD[T]00:00:00.0000[Z]");
        }
        return c;
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

    function addPieData(obj, timeline) {
        obj.pie = {
            labels: timeline.series,
            data: []
        };

        for (var i = 0; i < timeline.series.length; i++) {
            obj.pie.data.push(_.sum(timeline.data[i]));
        }
    }

    var service = {

        getReportOverview: function(criteria, callback) {

            var c = prepareCriteria(criteria);

            return Reporting.createOverview(c, function(overview) {
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

            var c = prepareCriteria(criteria);

            return Reporting.createProjects(c, function(projectsView) {
                projectsView.chartType = "Line";
                projectsView.noData = projectsView.timeline.labels.length === 0;
                projectsView.currentTimeline = addTimelineGrouping(projectsView, projectsView.timeline);
                addPieData(projectsView, projectsView.timeline);
            });
        },

        createCriteria: function() {
            return {
                startDate: $moment().startOf("month").toDate(),
                endDate: $moment().endOf("month").toDate(),
                projects: []
            };
        }
    };

    return service;
}]);
