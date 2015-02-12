angular.module('tima').controller('signinController', ['$scope', 'authService', function ($scope, authService) {
    $scope.formData = {
        username: '',
        password: '',
        clear: function() {
            this.username = '';
            this.password = '';
        }
    };

    $scope.signin = function() {
        $scope.$broadcast('show-errors-check-validity');

        if (!$scope.formSignin.$valid) {
            return;
        }

        authService.signIn($scope.formData, '/');
    };

    $scope.signout = function(){
        authService.signOut();
    };
}]);
