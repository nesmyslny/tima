angular.module('tima').factory('activityService',
['_', '$moment', 'sessionService', 'Activity', 'ProjectActivityType', 'util',
function(_, $moment, sessionService, Activity, ProjectActivityType, util) {

    function getDurationFormatted(duration) {
        var m = $moment.duration(duration, 'minutes');
        return util.formatTime(m.hours(), m.minutes());
    }

    function calculateDuration(hours, minutes) {
        return (hours || 0) * 60 + (minutes || 0);
    }

    function createNewActivity(day, projectActivity, hours, minutes) {
        return {
            id: -1,
            day: $moment(day, 'YYYY-MM-DD').format('YYYY-MM-DD[T]00:00:00.000[Z]'),
            userId: sessionService.user.id,
            projectId: projectActivity.projectId,
            activityTypeId: projectActivity.activityTypeId,
            duration: calculateDuration(hours, minutes),
            durationHours: hours || 0,
            durationMinutes: minutes || 0,
            durationFormatted: util.formatTime(hours, minutes),
            projectTitle: projectActivity.projectTitle,
            activityTypeTitle: projectActivity.activityTypeTitle
        };
    }

    function extendActivityDuration(activity, hours, minutes) {
        var duration = activity.duration + calculateDuration(hours, minutes);
        setActivityDuration(activity, duration);
    }

    function setActivityDuration(activity, duration) {
        activity.duration = duration;
        var m = $moment.duration(duration, 'minutes');
        activity.durationHours = m.hours();
        activity.durationMinutes = m.minutes();
        activity.durationFormatted = util.formatTime(m.hours(), m.minutes());
    }

    function refreshAfterSave(activities, activity, activityOrig) {
        if (_.isUndefined(activityOrig)) {
            activities.push(activity);
        } else {
            _.assign(activityOrig, activity);
        }
    }

    var service = {
        get: function(day, callback) {
            return Activity.query({day:day}, function() { callback(); });
        },

        getProjectActivityList: function() {
            return ProjectActivityType.query();
        },

        save: function(activities, activity, callback) {
            var clone = _.clone(activity);
            var duration = calculateDuration(clone.durationHours, clone.durationMinutes);
            setActivityDuration(clone, duration);

            Activity.save(clone, function() {
                refreshAfterSave(activities, clone, activity);
                callback();
            });
        },

        add: function(activities, day, projectActivity, hours, minutes, callback) {
            var activityOrig = _.find(activities, { "projectId": projectActivity.projectId, "activityTypeId": projectActivity.activityTypeId});
            var activity = {};

            if (_.isUndefined(activityOrig)) {
                activity = createNewActivity(day, projectActivity, hours, minutes);
            } else {
                activity = _.clone(activityOrig);
                extendActivityDuration(activity, hours, minutes);
            }

            Activity.save(activity, function() {
                refreshAfterSave(activities, activity, activityOrig);
                callback();
            });
        },

        delete: function(id, callback) {
            Activity.delete({id:id}, callback);
        },

        calculateTotalDuration: function(activities) {
            var duration = 0;
            _.forEach(activities, function(activity) {
                duration += activity.duration;
            });
            return {
                duration: duration,
                durationFormatted: getDurationFormatted(duration)
            };
        }
    };

    return service;
}]);
