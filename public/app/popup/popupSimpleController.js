angular.module('tima').controller('PopupSimpleController',
['$scope', '$uibModalInstance', 'title', 'body', 'acceptButton', 'cancelButton',
function ($scope, $uibModalInstance, title, body, acceptButton, cancelButton) {

    $scope.title = title;
    $scope.body = body;
    $scope.acceptButton = acceptButton;
    $scope.cancelButton = cancelButton;

    $scope.accept = function () {
        $uibModalInstance.close(true);
    };

    $scope.cancel = function () {
        $uibModalInstance.dismiss();
    };

}]);
