angular.module('tima').controller('ActivityController',
['activityService', '$scope', '$routeParams', '$location', '_', '$moment',
function (activityService, $scope, $routeParams, $location, _, $moment) {

    $scope.day = $routeParams.day;
    $scope.dayHeader = $moment($scope.day, 'YYYY-MM-DD').format('dddd, MMMM Do YYYY');
    $scope.durationHeader = '';
    $scope.totalDuration = 0;
    $scope.activities = [];
    $scope.projectActivityList = activityService.getProjectActivityList();

    $scope.formData = {
        projectActivity: null,
        hours: null,
        minutes: null,
        clear: function() {
            this.projectActivity = null;
            this.hours = null;
            this.minutes = null;
        }
    };

    $scope.list = function() {
        $scope.activities = activityService.get($scope.day, function() {
            refreshTotalDuration();
        });
    };
    $scope.list();

    $scope.add = function() {
        $scope.$broadcast("show-errors-check-validity");
        if (!$scope.formAddActivity.$valid) {
            return;
        }

        activityService.add($scope.activities, $scope.day, $scope.formData.projectActivity.selected, $scope.formData.hours, $scope.formData.minutes,
            function() {
                refreshTotalDuration();
            });

        $scope.formData.clear();
        $scope.$broadcast("show-errors-reset");

        // workaround: ui-select doesn't play nice with showErrors.js
        // todo: find solution (replace/remove showErrors.js? fix/enahnce showErrors.js? replace ui-select?)
        $scope.formAddActivity.$setPristine();
    };

    $scope.delete = function(activity) {
        activityService.delete(activity.id, function() {
            var index = _.indexOf($scope.activities, activity);
            $scope.activities.splice(index, 1);
            refreshTotalDuration();
        });
    };

    $scope.save = function(activity) {
        activityService.save($scope.activities, activity, function() {
            refreshTotalDuration();
        });
    };

    function refreshTotalDuration() {
        var totalDuration = activityService.calculateTotalDuration($scope.activities);
        $scope.totalDuration = totalDuration.duration;
        $scope.durationHeader = totalDuration.durationFormatted;
    }

    $scope.openDatePicker = function($event) {
        $event.preventDefault();
        $event.stopPropagation();
        $scope.datePickerOpened = !$scope.datePickerOpened;
    };

    $scope.today = function() {
        $scope.day = $moment().format('YYYY-MM-DD');
    };

    $scope.navigateDay = function(forward) {
        var i = forward ? 1 : -1;
        $scope.day = $moment($scope.day).add(i, 'days').format('YYYY-MM-DD');
    };

    $scope.$watch("day", dayWatchCallback, true);
    function dayWatchCallback(newVal, oldVal) {
        if (newVal != oldVal) {
            var date = $moment($scope.day).format('YYYY-MM-DD');
            $location.path('/activities/' + date);
        }
    }

}]);
