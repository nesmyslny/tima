angular.module('tima').factory('Reporting',
['$resource',
function($resource) {
    return $resource("/report", {}, {
        createOverview: {
            url: "/report/overview",
            method: "POST"
        },
        createProjects: {
            url: "/report/projects",
            method: "POST"
        }
    });
}]);
