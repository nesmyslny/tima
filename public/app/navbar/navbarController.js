angular.module('tima').controller('NavbarController',
['$scope', '$moment', 'authService',
function ($scope, $moment, authService) {
    $scope.user = {};
    $scope.isSignedIn = false;
    $scope.$watch(authService.getUser, userWatchCallback, true);
    $scope.activitiesPathToday = 'activities/' + $moment().format('YYYY-MM-DD');

    function userWatchCallback(newVal) {
        $scope.user = newVal ? newVal : {};
        $scope.isSignedIn = Boolean(newVal);
    }

    $scope.signout = function() {
        authService.signOut();
    };
}]);
