angular.module('tima').controller('ActivityController',
['activityService', '$scope', '$routeParams', '$location', 'authService',
function (activityService, $scope, $routeParams, $location, authService) {

    $scope.day = $routeParams.day;
    $scope.dayHeader = moment($scope.day, 'YYYY-MM-DD').format('dddd, MMMM Do YYYY');
    $scope.durationHeader = '';
    $scope.totalDuration = 0;
    $scope.activities = [];
    $scope.projectActivityList = [];

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
        var promise = activityService.refresh($scope.day, $scope.activities);
        promise.then(function(data) {
            $scope.totalDuration = data.minutes;
            $scope.durationHeader = data.formatted;
        });
    };
    $scope.list();

    $scope.fetchProjectActivityList = function() {
        activityService.getProjectActivityList().then(function(data) {
            $scope.projectActivityList = data;
            $scope.projectActivityList.forEach(function(item) {
                item.text = item.projectTitle + ": " + item.activityTypeTitle;
            });
        });
    };
    $scope.fetchProjectActivityList();

    $scope.add = function() {
        $scope.$broadcast('show-errors-check-validity');
        if (!$scope.formAddActivity.$valid) {
            return;
        }

        var activity = activityService.createNew(
            $scope.day,
            authService.getUser().id,
            $scope.formData.projectActivity.selected.projectId,
            $scope.formData.projectActivity.selected.activityTypeId,
            $scope.formData.hours,
            $scope.formData.minutes
        );

        activityService.save(activity).then(function() {
            $scope.list();
        });

        $scope.formData.clear();
        $scope.$broadcast('show-errors-reset');

        // workaround: ui-select doesn't play nice with showErrors.js
        // todo: find solution (replace/remove showErrors.js? fix/enahnce showErrors.js? replace ui-select?)
        $scope.formAddActivity.$setPristine();
    };

    $scope.delete = function(id) {
        activityService.delete(id).then(function() {
            $scope.list();
        });
    };

    $scope.save = function(activity) {
        activityService.save(activity).then(function() {
            $scope.list();
        });
    };

    $scope.openDatePicker = function($event) {
        $event.preventDefault();
        $event.stopPropagation();
        $scope.datePickerOpened = !$scope.datePickerOpened;
    };

    $scope.today = function() {
        $scope.day = moment().format('YYYY-MM-DD');
    };

    $scope.navigateDay = function(forward) {
        var i = forward ? 1 : -1;
        $scope.day = moment($scope.day).add(i, 'days').format('YYYY-MM-DD');
    };

    $scope.$watch("day", dayWatchCallback, true);
    function dayWatchCallback(newVal, oldVal) {
        if (newVal != oldVal) {
            var date = moment($scope.day).format('YYYY-MM-DD');
            $location.path('/activities/' + date);
        }
    }

}]);
