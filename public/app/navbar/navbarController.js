angular.module('tima').controller('NavbarController',
['$scope', 'authService',
function ($scope, authService) {
    $scope.username = '';
    $scope.isSignedIn = false;
    $scope.$watch(authService.getUser, usernameWatchCallback, true);
    $scope.activitiesPathToday = 'activities/' + moment().format('YYYY-MM-DD');

    function usernameWatchCallback(newVal) {
        $scope.username = newVal ? newVal.username : '';
        $scope.isSignedIn = Boolean(newVal);
    }

    $scope.signout = function() {
        authService.signOut();
    };
}]);
