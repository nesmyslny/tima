angular.module('tima').factory('ProjectActivityType',
['$resource',
function($resource) {
    return $resource("/projectActivityTypes");
}]);
