angular.module('tima').controller('MyProjectListController',
['$scope', 'Project',
function ($scope, Project) {
    $scope.projects = [];
    $scope.search = "";

    $scope.list = function() {
        $scope.projects = Project.queryUser(function() {
            $scope.initializePagination();
        });
    };
    $scope.list();

    $scope.initializePagination = function() {
        $scope.currentPage = 1;
        $scope.totalItems = $scope.projects.length;
        $scope.itemsPerPage = 10;
    };

    $scope.$watch('filteredProjects', function(newVal, oldVal) {
        $scope.currentPage = 1;
    }, true);
}]);
