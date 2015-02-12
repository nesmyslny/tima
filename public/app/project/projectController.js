angular.module('tima').controller('ProjectController',
['$scope', '$routeParams', '$location', 'Project', 'ActivityType',
function ($scope, $routeParams, $location, Project, ActivityType) {

    $scope.project = {
        id: -1,
        title: '',
        activityTypes: []
    };

    $scope.activityTypes = [];
    $scope.selectedActivityType = {};

    $scope.fetch = function() {
        var id = parseInt($routeParams.id);

        if (id > -1) {
            $scope.project = Project.get({id:id});
        }
    };
    $scope.fetch();

    $scope.list = function() {
        $scope.activityTypes = ActivityType.query();
    };
    $scope.list();

    $scope.addActivityType = function() {
        if (typeof $scope.selectedActivityType.selected === "undefined") {
            return;
        }

        var activityType = $scope.selectedActivityType.selected;
        var alreadyInList = $scope.project.activityTypes.some(function(at) {
            return at.id == activityType.id;
        });

        if (!alreadyInList) {
            $scope.project.activityTypes.push(activityType);
        }

        $scope.selectedActivityType = {};
    };

    $scope.deleteActivityType = function(activityType) {
        var index = $scope.project.activityTypes.indexOf(activityType);
        $scope.project.activityTypes.splice(index, 1);
    };

    $scope.save = function() {
        $scope.$broadcast('show-errors-check-validity');
        if (!$scope.formProject.$valid) {
            return;
        }

        Project.save($scope.project, function() {
            $location.path('/projects');
        });
    };

}]);
