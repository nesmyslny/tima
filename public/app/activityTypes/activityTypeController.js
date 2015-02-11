angular.module('tima').controller('activityTypeController', ['$scope', '$http', '$routeParams', '$location', function ($scope, $http, $routeParams, $location) {

    $scope.activityType = {
        id: -1,
        title: ''
    };

    $scope.fetch = function() {
        var id = parseInt($routeParams.id);

        if (id > -1) {
            $http.get('/activityTypes/' + id)
            .success(function(data, status, headers, config) {
                $scope.activityType = data;
            });
        }
    };
    $scope.fetch();

    $scope.save = function() {
        $scope.$broadcast('show-errors-check-validity');
        if (!$scope.formActivityType.$valid) {
            return;
        }

        $http.post('/activityTypes', $scope.activityType)
        .success(function(data, status, headers, config) {
            $location.path('/activityTypes');
        });
    };

}]);
