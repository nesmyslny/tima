angular.module('tima').factory('Department',
['$resource',
function($resource) {
    return $resource("/departments/:id", {}, {
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
