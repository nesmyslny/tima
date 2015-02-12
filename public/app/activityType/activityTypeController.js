angular.module('tima').controller('ActivityTypeController',
['$scope', '$routeParams', '$location', 'ActivityType',
function ($scope, $routeParams, $location, ActivityType) {

    $scope.activityType = {
        id: -1,
        title: ''
    };

    $scope.fetch = function() {
        var id = parseInt($routeParams.id);

        if (id > -1) {
            $scope.activityType = ActivityType.get({id:id});
        }
    };
    $scope.fetch();

    $scope.save = function() {
        $scope.$broadcast('show-errors-check-validity');
        if (!$scope.formActivityType.$valid) {
            return;
        }

        ActivityType.save($scope.activityType, function() {
            $location.path('/activityTypes');
        });
    };

}]);
