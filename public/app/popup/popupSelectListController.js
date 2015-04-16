angular.module('tima').controller('PopupSelectListController',
['$scope', '$modalInstance', 'title', 'items', 'acceptButton', 'cancelButton',
function ($scope, $modalInstance, title, items, acceptButton, cancelButton) {

    $scope.title = title;
    $scope.items = items;
    $scope.acceptButton = acceptButton;
    $scope.cancelButton = cancelButton;

    $scope.currentPage = 1;
    $scope.totalItems = $scope.items.length;
    $scope.itemsPerPage = 10;

    $scope.$watch('filteredItems', function(newVal, oldVal) {
        $scope.currentPage = 1;
    }, true);

    $scope.accept = function () {
        $modalInstance.close(true);
    };

    $scope.cancel = function () {
        $modalInstance.dismiss();
    };

}]);
