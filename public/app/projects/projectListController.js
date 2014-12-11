angular.module('tima').controller('projectListController', ['$scope', '$http', 'messageService', function ($scope, $http, messageService) {

    $scope.projects = [];

    $scope.list = function() {
        $http.get('/projects')
        .success(function(data, status, headers, config) {
            $scope.projects = data;
        });
    };
    $scope.list();

    $scope.delete = function(id) {
        $http.delete('/projects/' + id)
        .success(function() {
            $scope.list();
        })
        .error(function(data, status) {
            messageService.add('danger', data);
        });
    };

}]);
