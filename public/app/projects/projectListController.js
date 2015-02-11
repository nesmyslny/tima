angular.module('tima').controller('projectListController', ['$scope', '$http', 'messageService', 'popupService', function ($scope, $http, messageService, popupService) {

    $scope.projects = [];

    $scope.list = function() {
        $http.get('/projects')
        .success(function(data, status, headers, config) {
            $scope.projects = data;
            $scope.initializePagination();
        });
    };
    $scope.list();

    $scope.delete = function(id) {
        popupService.show('Delete Project', 'Do you really want to delete this project?', 'Delete', 'Cancel')
        .result.then(function() {
            $http.delete('/projects/' + id)
            .success(function() {
                $scope.list();
            })
            .error(function(data, status) {
                messageService.add('danger', data);
            });
        });
    };

    $scope.initializePagination = function() {
        $scope.currentPage = 1;
        $scope.totalItems = $scope.projects.length;
        $scope.itemsPerPage = 10;
    };

    $scope.$watch('filteredProjects', function(newVal, oldVal) {
        $scope.currentPage = 1;
    }, true);

}]);
