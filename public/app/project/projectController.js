angular.module('tima').controller('ProjectController',
['$scope', '$routeParams', '$location', '$q', '_', 'Project', 'ProjectCategory', 'ActivityType', 'Department', 'User', 'authService', 'userRoles', 'popupService', 'multiSelect',
function ($scope, $routeParams, $location, $q, _, Project, ProjectCategory, ActivityType, Department, User, authService, userRoles, popupService, multiSelect) {

    $scope.multiSelect = multiSelect;

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

    $scope.departments = [];
    $scope.selectedDepartment = {};

    $scope.users = [];
    $scope.usersResponsible = [];
    $scope.selectedResponsibleUser = {};
    $scope.selectedManagerUser = {};
    $scope.selectedUser = {};

    $scope.returnPath = $routeParams.returnPath;
    $scope.adminAccess = false;
    if (_.isUndefined($scope.returnPath)) {
        $scope.returnPath = "projects";

        if (authService.isAuthorized(userRoles.deptManager)) {
            $scope.adminAccess = true;
        }
    }

    $scope.fetch = function() {
        var id = parseInt($routeParams.id);

        $scope.projectCategories = ProjectCategory.queryList();
        $scope.activityTypes = ActivityType.query();
        $scope.departments = Department.queryList();
        $scope.users = User.queryAll();

        if ($scope.adminAccess && authService.isDeptManager()) {
            // department manager are allowed to set the responsible, but only users in the spedific department (or a descendant)
            $scope.usersResponsible = User.queryDept();
        } else {
            $scope.usersResponsible = $scope.users;
        }

        if (id > -1) {
            $scope.project = Project.get({id:id});
        }

        $q.all([
            $scope.project.$promise,
            $scope.projectCategories.$promise,
            $scope.users.$promise
        ]).then(function() {
            $scope.selectedProjectCategory.selected = _.find($scope.projectCategories, { 'id': $scope.project.projectCategoryId });
            $scope.selectedResponsibleUser.selected = _.find($scope.usersResponsible, { 'id': $scope.project.responsibleUserId });
            $scope.selectedManagerUser.selected = _.find($scope.users, { 'id': $scope.project.managerUserId });
        });
    };
    $scope.fetch();

    $scope.editDescription = function() {
        var markdown = $scope.project.description;
        popupService.showMarkdown("Description", markdown, "Ok", "Cancel")
        .result.then(function(markdownResult) {
            $scope.project.description = markdownResult;
        });
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
        return !($scope.adminAccess || ($scope.project.isResponsible && enableForResponsible));
    };

}]);
