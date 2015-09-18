angular.module('tima').controller('ProjectListController',
['$scope', 'Project', 'popupService',
function ($scope, Project, popupService) {

    $scope.projects = [];

    $scope.list = function() {
        $scope.projects = Project.queryAdmin(function() {
            $scope.initializePagination();
        });
    };
    $scope.list();

    $scope.delete = function(id) {
        popupService.showSimple('Delete Project', 'Do you really want to delete this project?', 'Delete', 'Cancel')
        .result.then(function() {
            Project.delete({id:id}, function() {
                $scope.list();
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
