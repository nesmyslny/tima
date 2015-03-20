angular.module('tima').controller('ProjectController',
['$scope', '$routeParams', '$location', '$q', '_', 'Project', 'ProjectCategory', 'ActivityType', 'User', 'authService', 'userRoles',
function ($scope, $routeParams, $location, $q, _, Project, ProjectCategory, ActivityType, User, authService, userRoles) {

    $scope.project = {
        id: -1,
        title: '',
        activityTypes: [],
        users: []
    };

    $scope.activityTypes = [];
    $scope.selectedActivityType = {};

    $scope.projectCategories = [];
    $scope.selectedProjectCategory = {};

    $scope.users = [];
    $scope.selectedResponsibleUser = {};
    $scope.selectedManagerUser = {};
    $scope.selectedUser = {};

    $scope.returnPath = $routeParams.returnPath;
    if (_.isUndefined($scope.returnPath)) {
        $scope.returnPath = "projects";
    }

    $scope.fetch = function() {
        var id = parseInt($routeParams.id);

        $scope.projectCategories = ProjectCategory.queryList();
        $scope.activityTypes = ActivityType.query();
        $scope.users = User.query();

        if (id > -1) {
            $scope.project = Project.get({id:id});
        }

        $q.all([
            $scope.project.$promise,
            $scope.projectCategories.$promise,
            $scope.users.$promise
        ]).then(function() {
            $scope.selectedProjectCategory.selected = _.find($scope.projectCategories, { 'id': $scope.project.projectCategoryId });
            $scope.selectedResponsibleUser.selected = _.find($scope.users, { 'id': $scope.project.responsibleUserId });
            $scope.selectedManagerUser.selected = _.find($scope.users, { 'id': $scope.project.managerUserId });
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

    $scope.addUser = function() {
        if (_.isUndefined($scope.selectedUser.selected)) {
            return;
        }

        var user = $scope.selectedUser.selected;
        var alreadyInList = $scope.project.users.some(function(u) {
            return u.id == user.id;
        });

        if (!alreadyInList) {
            $scope.project.users.push(user);
        }

        $scope.selectedUser = {};
    };

    $scope.deleteUser = function(user) {
        var index = $scope.project.users.indexOf(user);
        $scope.project.users.splice(index, 1);
    };

    $scope.save = function() {
        $scope.$broadcast('show-errors-check-validity');
        if (!$scope.formProject.$valid) {
            return;
        }

        $scope.project.projectCategoryId = $scope.selectedProjectCategory.selected.id;
        $scope.project.responsibleUserId = $scope.selectedResponsibleUser.selected.id;
        $scope.project.managerUserId = $scope.selectedManagerUser.selected.id;

        Project.save($scope.project, function() {
            $location.path($scope.returnPath);
        });
    };

    $scope.disableForUsers = function(enableForResponsible) {
        return !(authService.isAuthorized(userRoles.manager) || ($scope.project.isResponsible && enableForResponsible));
    };

}]);
