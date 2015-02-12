angular.module('tima').factory('activityService',
['$filter', 'Activity', 'ProjectActivityType',
function($filter, Activity, ProjectActivityType) {

    function removeDeletedActivities(source, dest) {
        dest.forEach(function(activity) {
            var actFound = $filter('filter')(source, {id: activity.id}, true);
            if (!actFound.length) {
                var index = dest.indexOf(activity);
                dest.splice(index, 1);
            }
        });
    }

    function mergeActivities(source, dest) {
        source.forEach(function(activity) {
            var actFound = $filter('filter')(dest, {id: activity.id}, true);
            if (actFound.length) {
                actFound[0].duration = activity.duration;
            } else {
                dest.push(activity);
            }
        });
    }

    function refreshActivitiesViewValues(activities) {
        activities.forEach(function(activity) {
            var m = moment.duration(activity.duration, 'minutes');
            activity.durationHours = m.hours();
            activity.durationMinutes = m.minutes();
            activity.durationFormatted = getTimeFormatted(activity.durationHours, activity.durationMinutes);
        });
    }

    function getTotalDuration(activities) {
        var totalDuration = 0;
        activities.forEach(function(activity) {
            totalDuration += activity.duration;
        });

        return {
            minutes: totalDuration,
            formatted: getDurationFormatted(totalDuration)
        };
    }

    function getDurationFormatted(duration) {
        var m = moment.duration(duration, 'minutes');
        return getTimeFormatted(m.hours(), m.minutes());
    }

    function getTimeFormatted(hours, minutes) {
        var durationFormatted = hours > 0 ? hours + 'h' : '';
        durationFormatted += minutes > 0 ? ' ' + minutes + 'min' : '';
        return durationFormatted;
    }

    var service = {
        refresh: function(day, activities) {
            return Activity.query({day:day})
            .$promise.then(function(data) {
                removeDeletedActivities(data, activities);
                mergeActivities(data, activities);
                refreshActivitiesViewValues(activities);
                return getTotalDuration(activities);
            });
        },

        getProjectActivityList: function() {
            return ProjectActivityType.query().$promise;
        },

        save: function(activity) {
            return Activity.save(activity).$promise;
        },

        delete: function(id) {
            return Activity.delete({id:id}).$promise;
        },

        createNew: function(day, userId, projectId, activityTypeId, hours, minutes) {
            return {
                id: -1,
                day: moment(day, 'YYYY-MM-DD').format('YYYY-MM-DD[T]00:00:00.000[Z]'),
                userId: userId,
                projectId: projectId,
                activityTypeId: activityTypeId,
                duration: service.calculateDuration(hours, minutes)
            };
        },

        calculateDuration: function(hours, minutes) {
            return hours * 60 + minutes;
        }
    };

    return service;
}]);
