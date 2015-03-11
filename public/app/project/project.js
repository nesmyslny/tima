angular.module('tima').factory('Project',
['$resource', 'resourceSaveInterceptor',
function($resource, resourceSaveInterceptor) {
    return $resource("/projects/:id", {}, {
        save: {
            method: "POST",
            interceptor: resourceSaveInterceptor
        },
    });
}]);
