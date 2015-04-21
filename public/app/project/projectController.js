angular.module('tima').controller('ProjectController',
['$scope', '$routeParams', '$location', '$q', '_', 'Project', 'ProjectCategory', 'ActivityType', 'Department', 'User', 'authService', 'userRoles', 'popupService',
function ($scope, $routeParams, $location, $q, _, Project, ProjectCategory, ActivityType, Department, User, authService, userRoles, popupService) {

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

    $scope.addActivityType = function() {
        addProjectItem($scope.project.activityTypes, $scope.selectedActivityType.selected);
        $scope.selectedActivityType = {};
    };

    $scope.addMultipleActivityTypes = function() {
        addMultipleProjectItems($scope.activityTypes, $scope.project.activityTypes, "Activity Types", "title");
    };

    $scope.deleteActivityType = function(activityType) {
        deleteProjectItem($scope.project.activityTypes, activityType);
    };

    $scope.addDepartment = function() {
        addProjectItem($scope.project.departments, $scope.selectedDepartment.selected);
        $scope.selectedDepartment = {};
    };

    $scope.addMultipleDepartments = function() {
        addMultipleProjectItems($scope.departments, $scope.project.departments, "Departments", "path");
    };

    $scope.deleteDepartment = function(dept) {
        deleteProjectItem($scope.project.departments, dept);
    };

    $scope.addUser = function() {
        addProjectItem($scope.project.users, $scope.selectedUser.selected);
        $scope.selectedUser = {};
    };

    $scope.addMultipleUsers = function() {
        addMultipleProjectItems($scope.users, $scope.project.users, "Users", "username");
    };

    $scope.deleteUser = function(user) {
        deleteProjectItem($scope.project.users, user);
    };

    function addProjectItem(projectItems, item) {
        if (_.isUndefined(item)) {
            return;
        }

        var alreadyInList = projectItems.some(function(projectItem) {
            return projectItem.id == item.id;
        });

        if (!alreadyInList) {
            projectItems.push(item);
        }
    }

    function addMultipleProjectItems(items, projectItems, popupTitle, valueKey) {
        var popupItems = [];
        _.forEach(items, function(item) {
            popupItems.push({
                value: item[valueKey],
                obj: item,
                checked: _.any(projectItems, "id", item.id)
            });
        });

        popupService.showSelectList(popupTitle, popupItems, "Ok", "Cancel")
        .result.then(function() {
            projectItems.length = 0;
            _.forEach(popupItems, function(item) {
                if (item.checked) {
                    addProjectItem(projectItems, item.obj);
                }
            });
        });
    }

    function deleteProjectItem(projectItems, item) {
        var index = projectItems.indexOf(item);
        projectItems.splice(index, 1);
    }

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
