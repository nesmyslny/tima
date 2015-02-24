angular.module('tima').controller('UserController',
['$scope', '$routeParams', '$location', '_', 'User', 'sessionService', 'messageService',
function ($scope, $routeParams, $location, _, User, sessionService, messageService) {

    var modeUserSettings = false;
    $scope.user = {
        id: -1,
        title: '',
    };

    $scope.fetch = function() {
        var id = -1;

        if (_.isUndefined($routeParams.id)) {
            id = sessionService.user.id;
            modeUserSettings = true;
        } else {
            id = parseInt($routeParams.id);
        }

        if (id > -1) {
            $scope.user = User.get({id:id}, function() {
                $scope.user.newPassword = undefined;
                $scope.user.newPasswordConfirm = undefined;
            });
        }
    };
    $scope.fetch();

    $scope.save = function() {
        $scope.$broadcast('show-errors-check-validity');
        if (!$scope.formUser.$valid) {
            return;
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
