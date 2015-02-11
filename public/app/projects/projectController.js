angular.module('tima').controller('projectController', ['$scope', '$http', '$routeParams', '$location', 'messageService', function ($scope, $http, $routeParams, $location, messageService) {

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
            $http.get('/projects/' + id)
            .success(function(data, status, headers, config) {
                $scope.project = data;
            });
        }
    };
    $scope.fetch();

    $scope.list = function() {
        $http.get('/activityTypes')
        .success(function(data, status, headers, config) {
            $scope.activityTypes = data;
        });
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

        $http.post('/projects', $scope.project)
        .success(function(data, status, headers, config) {
            $location.path('/projects');
        })
        .error(function(data, status) {
            messageService.add('danger', data);
        });
    };

}]);
