angular.module('tima').factory('ProjectCategory',
['$resource', 'resourceSaveInterceptor',
function($resource, resourceSaveInterceptor) {
    return $resource("/projectCategories/:id", {}, {
        save: {
            method: "POST",
            interceptor: resourceSaveInterceptor
        },
        queryTree: {
            url: "/projectCategories/tree",
            method: "GET",
            isArray: true
        },
        queryList: {
            url: "/projectCategories/list",
            method: "GET",
            isArray: true
        }
    });
}]);
