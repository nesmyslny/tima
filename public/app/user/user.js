angular.module('tima').factory('User',
['$resource',
function($resource) {
    return $resource("/users/:id");
}]);
