angular.module('tima').controller('PopupController',
['$scope', '$modalInstance', 'title', 'body', 'acceptButton', 'cancelButton',
function ($scope, $modalInstance, title, body, acceptButton, cancelButton) {

    $scope.title = title;
    $scope.body = body;
    $scope.acceptButton = acceptButton;
    $scope.cancelButton = cancelButton;

    $scope.accept = function () {
        $modalInstance.close(true);
    };

    $scope.cancel = function () {
        $modalInstance.dismiss();
    };

}]);
