angular.module('tima').controller('NavbarController', ['$scope', '$window', 'authService', function ($scope, $window, authService) {

    $scope.username = '';
    $scope.isSignedIn = false;
    $scope.$watch(authService.getUsername, usernameWatchCallback, true);

    function usernameWatchCallback(newVal) {
        $scope.username = newVal;
        $scope.isSignedIn = Boolean(newVal);
    }

    $scope.signout = function(){
        authService.signOut();
    };

}]);
