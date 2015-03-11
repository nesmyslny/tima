angular.module('tima').factory('User',
['$resource', 'resourceSaveInterceptor',
function($resource, resourceSaveInterceptor) {
    return $resource("/users/:id", {}, {
        save: {
            method: "POST",
            interceptor: resourceSaveInterceptor
        }
    });
}]);
