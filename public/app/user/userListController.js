angular.module('tima').controller('UserListController',
['$scope', 'User',
function ($scope, User) {

    $scope.users = [];
    $scope.search = "";

    $scope.list = function() {
        $scope.users = User.query(function() {
            $scope.initializePagination();
        });
    };
    $scope.list();

    $scope.initializePagination = function() {
        $scope.currentPage = 1;
        $scope.totalItems = $scope.users.length;
        $scope.itemsPerPage = 10;
    };

    $scope.$watch('filteredUsers', function(newVal, oldVal) {
        $scope.currentPage = 1;
    }, true);

}]);
