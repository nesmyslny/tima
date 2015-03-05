angular.module('tima').controller('UserController',
['$scope', '$routeParams', '$location', '$q', '_', 'User', 'Department', 'sessionService', 'messageService', 'userRoles',
function ($scope, $routeParams, $location, $q, _, User, Department, sessionService, messageService, userRoles) {

    var modeUserSettings = false;
    $scope.user = {
        id: -1,
        title: '',
    };

    $scope.departments = [];
    $scope.selectedDepartment = {};

    $scope.roles = _.values(userRoles);
    $scope.selectedRole = {};

    $scope.fetch = function() {
        var id = -1;

        if (_.isUndefined($routeParams.id)) {
            id = sessionService.user.id;
            modeUserSettings = true;
        } else {
            id = parseInt($routeParams.id);
            $scope.departments = Department.queryList();
        }

        if (id > -1) {
            $scope.user = User.get({id:id}, function() {
                $scope.user.newPassword = undefined;
                $scope.user.newPasswordConfirm = undefined;
            });

            if (!modeUserSettings) {
                $q.all([
                    $scope.user.$promise,
                    $scope.departments.$promise
                ]).then(function() {
                    $scope.selectedDepartment.selected = _.find($scope.departments, { 'id': $scope.user.departmentId });
                    $scope.selectedRole.selected = _.find($scope.roles, { 'id': $scope.user.role });
                });
            }
        }
    };
    $scope.fetch();

    $scope.save = function() {
        $scope.$broadcast('show-errors-check-validity');
        if (!$scope.formUser.$valid) {
            return;
        }

        if (!modeUserSettings) {
            $scope.user.departmentId = $scope.selectedDepartment.selected.id;
            $scope.user.role = $scope.selectedRole.selected.id;
        }

        User.save($scope.user, function() {
            $scope.user.newPassword = $scope.user.newPasswordConfirm = "";

            // if the updated user is the signed in user -> update in session
            if ($scope.user.id === sessionService.user.id) {
                sessionService.updateUser($scope.user);
            }

            if (modeUserSettings) {
                messageService.add('success', 'User settings saved.');
            } else {
                $location.path('/users');
            }
        });
    };

}]);
