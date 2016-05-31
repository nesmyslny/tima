angular.module('tima').controller('ReportingController',
['$scope', '_', 'reportingService', 'Project', 'multiSelect',
function ($scope, _, reportingService, Project, multiSelect) {

    $scope.multiSelect = multiSelect;
    $scope.criteria = reportingService.createCriteria();
    $scope.overview = {};
    $scope.projectsView = {};
    $scope.openDatePopups = [];

    // todo: keep user rights in mind! users aren't allowed to analyse all projects.
    $scope.projects = Project.queryAdmin();
    $scope.selectedProject = {};

    // Chart.js on tabs causing issues (flickering of different states; switching 'Days'/'Weeds'/'Months')
    // Workaround: The content of the tab is in a div, which is only shown, when the tab is active (via ng-if!)
    $scope.selectedTab = "overview";

    $scope.refreshReport = function() {
        $scope.overview = reportingService.getReportOverview($scope.criteria);
        $scope.projectsView = reportingService.getReportProjects($scope.criteria);
    };

    $scope.openDatePopup = function($event) {
        $event.preventDefault();
        $event.stopPropagation();

        var id = $event.currentTarget.id;
        $scope.openDatePopups[id] = !$scope.openDatePopups[id];
    };

}]);
