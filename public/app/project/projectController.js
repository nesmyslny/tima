angular.module('tima').controller('ProjectController',
['$scope', '$routeParams', '$location', '$q', '_', 'Project', 'ProjectCategory', 'ActivityType',
function ($scope, $routeParams, $location, $q, _, Project, ProjectCategory, ActivityType) {

    $scope.project = {
        id: -1,
        title: '',
        activityTypes: []
    };

    $scope.activityTypes = [];
    $scope.selectedActivityType = {};

    $scope.projectCategories = [];
    $scope.selectedProjectCategory = {};

    $scope.fetch = function() {
        var id = parseInt($routeParams.id);

        $scope.projectCategories = ProjectCategory.queryList();
        $scope.activityTypes = ActivityType.query();

        if (id > -1) {
            $scope.project = Project.get({id:id});
        }

        $q.all([
            $scope.project.$promise,
            $scope.projectCategories.$promise
        ]).then(function() {
            $scope.selectedProjectCategory.selected = _.find($scope.projectCategories, { 'id': $scope.project.projectCategoryId });
        });
    };
    $scope.fetch();

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

        $scope.project.projectCategoryId = $scope.selectedProjectCategory.selected.id;

        Project.save($scope.project, function() {
            $location.path('/projects');
        });
    };

}]);
