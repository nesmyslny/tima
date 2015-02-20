angular.module('tima').controller('ActivityTypeListController',
['$scope','ActivityType', 'messageService', 'popupService',
function ($scope, ActivityType, messageService, popupService) {

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

    $scope.add = function() {
        var activityType = {
            id: -1,
            title: "",
        };

        popupService.showForm("Add Activity Type", "app/activityType/activityTypePopupTemplate.html", activityType, "Add", "Cancel")
        .result.then(function() {
            ActivityType.save(activityType, function(response) {
                activityType.id = response.intResult;
                $scope.activityTypes.push(activityType);
            });
        });
    };

    $scope.edit = function(activityType) {
        var data = _.clone(activityType);

        popupService.showForm("Edit Activity Type", "app/activityType/activityTypePopupTemplate.html", data, "Save", "Cancel")
        .result.then(function() {
            ActivityType.save(data, function() {
                _.assign(activityType, data);
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
