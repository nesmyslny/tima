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

    function insertMissingDays(timeline, startDate, endDate) {
        var start = null;
        if (startDate) {
            start = $moment(startDate).format("YYYY-MM-DD");
        } else if (timeline.length > 0) {
            start = timeline[0].date;
        }

        var end = null;
        if (endDate) {
            end = $moment(endDate).format("YYYY-MM-DD");
        } else if (timeline.length > 0) {
            end = timeline[timeline.length - 1].date;
        }

        var date = start;

        for (i = 0; date != end; i++) {
            var tlDate = null;
            if (!_.isUndefined(timeline[i])) {
                tlDate = timeline[i].date;
            }

            if (tlDate != date) {
                timeline.splice(i, 0, {date: date, value: 0});
            }
            date = $moment(date).add(1, "d").format("YYYY-MM-DD");
        }
    }

    function initializeGrouping(overview) {
        overview.timelineDay = {
            labels: _.pluck(overview.timeline, "date"),
            values: [_.pluck(overview.timeline, "value")]
        };

        overview.timelineWeek = groupTimeline(overview.timeline, "YYYY[W]WW");
        overview.timelineMonth = groupTimeline(overview.timeline, "YYYY-MM");

        if (overview.timeline.length > 180) {
            return overview.timelineMonth;
        } else if (overview.timeline.length > 31) {
            return overview.timelineWeek;
        }
        return overview.timelineDay;
    }

    function groupTimeline(timeline, groupByFormat) {
        var groupedTimeline = {};

        _.forEach(timeline, function(item) {
            var group = $moment(item.date).format(groupByFormat);

            if (_.isUndefined(groupedTimeline[group])) {
                groupedTimeline[group] = 0;
            }

            groupedTimeline[group] += item.value;
        });

        return {
            labels: _.keys(groupedTimeline),
            values: [_.values(groupedTimeline)]
        };
    }

    var service = {

        getOverview: function(criteria, callback) {

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

                insertMissingDays(overview.timeline, criteria.startDate, criteria.endDate);
                overview.currentTimeline = initializeGrouping(overview);
            });
        },

        createCriteria: function() {
            return {
                startDate: $moment().startOf("month"),
                endDate: $moment().endOf("month")
            };
        }
    };

    return service;
}]);
