angular.module('tima').controller('activitiesController', ['activitiesService', '$scope', '$routeParams', '$location', function (activitiesService, $scope, $routeParams, $location) {

    $scope.day = $routeParams.day;
    $scope.dayHeader = moment($scope.day, 'YYYY-MM-DD').format('dddd, MMMM Do YYYY');
    $scope.durationHeader = '';
    $scope.totalDuration = 0;
    $scope.activities = [];

    $scope.formData = {
        text: '',
        hours: null,
        minutes: null,
        clear: function() {
            this.text = '';
            this.hours = '';
            this.minutes = '';
        }
    };

    $scope.list = function() {
        var promise = activitiesService.refreshActivities($scope.day, $scope.activities);
        promise.then(function(data) {
            $scope.totalDuration = data.totalDuration.minutes;
            $scope.durationHeader = data.totalDuration.formatted;
        });
    };
    $scope.list();

    $scope.add = function() {
        $scope.$broadcast('show-errors-check-validity');
        if (!$scope.formAddActivity.$valid) {
            return;
        }

        activitiesService.saveActivity($scope.day, $scope.formData.text, $scope.formData.hours, $scope.formData.minutes)
        .then(function() {
            $scope.list();
        });

        $scope.formData.clear();
        $scope.$broadcast('show-errors-reset');
    };

    $scope.delete = function(id) {
        activitiesService.deleteActivity(id)
        .then(function() {
            $scope.list();
        });
    };

    $scope.changeDuration = function(activity) {
        activitiesService.changeActivityDuration(activity)
        .then(function() {
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
