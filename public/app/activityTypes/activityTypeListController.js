angular.module('tima').controller('activityTypeListController', ['$scope', '$http', 'messageService', 'popupService', function ($scope, $http, messageService, popupService) {

    $scope.activityTypes = [];

    $scope.list = function() {
        $http.get('/activityTypes')
        .success(function(data, status, headers, config) {
            $scope.activityTypes = data;
            $scope.initializePagination();
        });
    };
    $scope.list();

    $scope.delete = function(id) {
        popupService.show('Delete Activity Type', 'Do you really want to delete this activity type?', 'Delete', 'Cancel')
        .result.then(function() {
            $http.delete('/activityTypes/' + id)
            .success(function() {
                $scope.list();
            })
            .error(function(data, status) {
                messageService.add('danger', data);
            });
        });
    };

    $scope.initializePagination = function() {
        $scope.currentPage = 1;
        $scope.totalItems = $scope.activityTypes.length;
        $scope.itemsPerPage = 10;
    };

    $scope.$watch('filteredActivityTypes ', function(newVal, oldVal) {
        $scope.currentPage = 1;
    }, true);

}]);
