angular.module('tima').factory('Activity', ['$resource', function($resource) {
    return $resource("/activities/:id", {}, {
        query: {
            url: "/activities/:day",
            method: "GET",
            isArray: true
        }
    });
}]);
