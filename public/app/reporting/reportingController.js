angular.module('tima').controller('ReportingController',
['$scope', '_', 'reportingService', 'Project', 'multiSelect',
function ($scope, _, reportingService, Project, multiSelect) {

    $scope.multiSelect = multiSelect;
    $scope.criteria = reportingService.createCriteria();
    $scope.overview = {};
    $scope.openDatePopups = [];

    // todo: keep user rights in mind! users aren't allowed to analyse all projects.
    $scope.projects = Project.queryAdmin();
    $scope.selectedProject = {};

    $scope.refreshReport = function() {
        $scope.overview = reportingService.getOverview($scope.criteria);
    };
    $scope.refreshReport();

    $scope.openDatePopup = function($event) {
        $event.preventDefault();
        $event.stopPropagation();

        var id = $event.currentTarget.id;
        $scope.openDatePopups[id] = !$scope.openDatePopups[id];
    };

}]);
