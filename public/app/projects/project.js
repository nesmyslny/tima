angular.module('tima').factory('Project', ['$resource', function($resource) {
    return $resource("/projects/:id");
}]);
