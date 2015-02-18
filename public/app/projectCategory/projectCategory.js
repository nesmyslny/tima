angular.module('tima').factory('ProjectCategory',
['$resource',
function($resource) {
    return $resource("/projectCategories/:id", {}, {
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
