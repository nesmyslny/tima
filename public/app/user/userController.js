angular.module('tima').controller('UserController',
['$scope', '$routeParams', '$location', '$q', '_', 'User',
function ($scope, $routeParams, $location, $q, _, User) {

    $scope.user = {
        id: -1,
        title: '',
    };

    $scope.fetch = function() {
        var id = parseInt($routeParams.id);

        if (id > -1) {
            $scope.user = User.get({id:id});
        }
    };
    $scope.fetch();

    $scope.save = function() {
        $scope.$broadcast('show-errors-check-validity');
        if (!$scope.formUser.$valid) {
            return;
        }

        User.save($scope.user, function() {
            $location.path('/users');
        });
    };

}]);
