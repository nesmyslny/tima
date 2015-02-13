angular.module('tima').factory('ProjectCategory',
['$resource',
function($resource) {
    return $resource("/projectCategories/:id");
}]);
