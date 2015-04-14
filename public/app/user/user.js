angular.module('tima').factory('User',
['$resource', 'resourceSaveInterceptor',
function($resource, resourceSaveInterceptor) {
    return $resource("/users/:id", {}, {
        queryAll: {
            method: "GET",
            url: "/users/all",
            isArray: true
        },
        queryDept: {
            method: "GET",
            url: "/users/department",
            isArray: true
        },
        save: {
            method: "POST",
            interceptor: resourceSaveInterceptor
        }
    });
}]);
