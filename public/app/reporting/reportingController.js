angular.module('tima').controller('ReportingController',
['$scope', '_', '$moment', 'Reporting',
function ($scope, _, $moment, Reporting) {

    $scope.criteria = {
        startDate: null,
        endDate: null
    };

    $scope.overview = {};
    $scope.overviewChartType = "Line";
    $scope.openDatePopups = [];

    $scope.refreshReport = function() {
        if ($scope.criteria.startDate) {
            $scope.criteria.startDate = $moment($scope.criteria.startDate).format("YYYY-MM-DD[T]00:00:00.000[Z]");
        }

        if ($scope.criteria.endDate) {
            $scope.criteria.endDate = $moment($scope.criteria.endDate).format("YYYY-MM-DD[T]00:00:00.0000[Z]");
        }

        $scope.overview = Reporting.createOverview($scope.criteria, function() {
            var duration = $scope.overview.durationTotal;

            if (duration < 1440) {
                var hours = $moment.duration(duration, "minutes").asHours();
                $scope.overview.durationHours = +hours.toFixed(2);

            } else {
                var days = $moment.duration(duration, "minutes").asDays();
                $scope.overview.durationDays = +days.toFixed(2);
            }

            $scope.overview.timelineData = [_.pluck($scope.overview.timeline, "value")];
            $scope.overview.timelineLabels = _.pluck($scope.overview.timeline, "label");
        });
    };
    $scope.refreshReport();

    $scope.openDatePopup = function($event) {
        $event.preventDefault();
        $event.stopPropagation();

        var id = $event.currentTarget.id;
        $scope.openDatePopups[id] = !$scope.openDatePopups[id];
    };

}]);
