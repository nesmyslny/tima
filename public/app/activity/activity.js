angular.module('tima').factory('Activity',
['$resource', '_', '$moment', 'resourceSaveInterceptor', 'util',
function($resource, _, $moment, resourceSaveInterceptor, util) {
    return $resource("/activities/:id", {}, {
        query: {
            url: "/activities/:day",
            method: "GET",
            isArray: true,
            transformResponse: function(data, headers) {
                var contentType = headers("content-type");

                if (contentType && _.startsWith(contentType, "application/json")) {
                    data = angular.fromJson(data);

                    _.forEach(data, function(activity) {
                        var m = $moment.duration(activity.duration, 'minutes');
                        activity.durationHours = m.hours();
                        activity.durationMinutes = m.minutes();
                        activity.durationFormatted = util.formatTime(activity.durationHours, activity.durationMinutes);
                    });
                }

                return data;
            }
        },
        save: {
            method: "POST",
            interceptor: resourceSaveInterceptor
        }
    });
}]);
