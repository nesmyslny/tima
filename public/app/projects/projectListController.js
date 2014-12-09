angular.module('tima').controller('projectListController', ['$scope', '$http', '$location', function ($scope, $http, $location) {

    $scope.projects = [];

    $scope.list = function() {
        $http.get('/projects')
        .success(function(data, status, headers, config) {
            $scope.projects = data;
        });
    };
    $scope.list();

    $scope.new = function() {
        $location.path('projects/-1');
    };

    $scope.delete = function(id) {
        $http.delete('/projects/' + id)
        .success(function() {
            $scope.list();
        })
        .error(function(data, status) {
            // todo: error handling
        });
    };

}]);
