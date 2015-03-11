angular.module('tima').factory('Department',
['$resource', 'resourceSaveInterceptor',
function($resource, resourceSaveInterceptor) {
    return $resource("/departments/:id", {}, {
        save: {
            method: "POST",
            interceptor: resourceSaveInterceptor
        },
        queryTree: {
            url: "/departments/tree",
            method: "GET",
            isArray: true
        },
        queryList: {
            url: "/departments/list",
            method: "GET",
            isArray: true
        }
    });
}]);
