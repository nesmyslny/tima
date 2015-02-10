angular.module('tima').controller('projectController', ['$scope', '$http', '$routeParams', '$location', function ($scope, $http, $routeParams, $location) {

    $scope.project = {
        id: -1,
        title: '',
        activityTypes: []
    };

    $scope.fetch = function() {
        var id = parseInt($routeParams.id);

        if (id > -1) {
            $http.get('/projects/' + id)
            .success(function(data, status, headers, config) {
                $scope.project = data;
            });
        }
    };
    $scope.fetch();

    $scope.deleteActivityType = function(activityType) {
        var index = $scope.project.activityTypes.indexOf(activityType);
        $scope.project.activityTypes.splice(index, 1);
    };

    $scope.save = function() {
        $scope.$broadcast('show-errors-check-validity');
        if (!$scope.formProject.$valid) {
            return;
        }

        $http.post('/projects', $scope.project)
        .success(function(data, status, headers, config) {
            $location.path('/projects');
        });
    };

}]);
