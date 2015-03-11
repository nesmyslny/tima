angular.module('tima').factory('ActivityType',
['$resource', 'resourceSaveInterceptor',
function($resource, resourceSaveInterceptor) {
    return $resource("/activityTypes/:id", {}, {
        save: {
            method: "POST",
            interceptor: resourceSaveInterceptor
        }
    });
}]);
