angular.module('tima').controller('activityTypeListController', ['$scope','ActivityType', 'messageService', 'popupService', function ($scope, ActivityType, messageService, popupService) {

    $scope.activityTypes = [];

    $scope.list = function() {
        $scope.activityTypes = ActivityType.query(function(data) {
            $scope.initializePagination();
        });
    };
    $scope.list();

    $scope.delete = function(id) {
        popupService.show('Delete Activity Type', 'Do you really want to delete this activity type?', 'Delete', 'Cancel')
        .result.then(function() {
            ActivityType.delete({id:id}, function() {
                $scope.list();
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
