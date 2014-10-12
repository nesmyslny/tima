angular.module('gnomon').controller('SigninController', ['$scope', '$http', function ($scope, $http) {
    $scope.formData = {
        username: '',
        password: ''
    };

    $scope.signin = function() {
        $scope.$broadcast('show-errors-check-validity');

        if (!$scope.formSignin.$valid) {
            return;
        }

        $http.post('/signin', $scope.formData)
        .success(function(data, status, headers, config) {
            alert('success');
        })
        .error(function(data, status, headers, config) {
            $scope.formData.username = '';
            $scope.formData.password = '';
        });
    };
}]);
