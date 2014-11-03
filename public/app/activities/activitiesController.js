angular.module('tima').controller('activitiesController', ['activitiesService', '$scope', '$routeParams', '$location', 'authService', function (activitiesService, $scope, $routeParams, $location, authService) {

    $scope.day = $routeParams.day;
    $scope.dayHeader = moment($scope.day, 'YYYY-MM-DD').format('dddd, MMMM Do YYYY');
    $scope.durationHeader = '';
    $scope.totalDuration = 0;
    $scope.activities = [];

    $scope.formData = {
        projectId: null,
        hours: null,
        minutes: null,
        clear: function() {
            this.projectId = null;
            this.hours = null;
            this.minutes = null;
        }
    };

    $scope.list = function() {
        var promise = activitiesService.refresh($scope.day, $scope.activities);
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

        var activity = activitiesService.createNew($scope.day, authService.getUser().id, $scope.formData.projectId, $scope.formData.hours, $scope.formData.minutes);
        activitiesService.save(activity)
        .then(function() {
            $scope.list();
        });

        $scope.formData.clear();
        $scope.$broadcast('show-errors-reset');
    };

    $scope.delete = function(id) {
        activitiesService.delete(id)
        .then(function() {
            $scope.list();
        });
    };

    $scope.changeDuration = function(activity) {
        activity.duration = activitiesService.calculateDuration(activity.durationHours, activity.durationMinutes);
        activitiesService.save(activity)
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
