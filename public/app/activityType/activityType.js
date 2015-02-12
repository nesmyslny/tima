angular.module('tima').factory('ActivityType',
['$resource',
function($resource) {
    return $resource("/activityTypes/:id");
}]);
