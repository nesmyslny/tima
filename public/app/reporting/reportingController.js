angular.module('tima').controller('ReportingController',
['$scope', '_', 'reportingService',
function ($scope, _, reportingService) {

    $scope.criteria = reportingService.createCriteria();
    $scope.overview = {};
    $scope.openDatePopups = [];

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
