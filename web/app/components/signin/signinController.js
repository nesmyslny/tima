angular.module('gnomon').controller('SigninController', ['$scope', '$http', 'authService', function ($scope, $http, authService) {
    $scope.formData = {
        username: '',
        password: '',
        clear: function() {
            this.username = '';
            this.password = '';
        }
    };

    $scope.secretData = null;

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

    $scope.secret = function() {
        $http.get('/secret')
        .success(function(data, status, headers, config) {
            $scope.secretData = data;
        });
    };

}]);
