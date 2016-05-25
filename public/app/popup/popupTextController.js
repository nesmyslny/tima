angular.module('tima').controller('PopupTextController',
['$scope', '$uibModalInstance', 'title', 'text', 'acceptButton', 'cancelButton',
function ($scope, $uibModalInstance, title, text, acceptButton, cancelButton) {

    $scope.title = title;
    $scope.text = text;
    $scope.acceptButton = acceptButton;
    $scope.cancelButton = cancelButton;

    $scope.accept = function() {
        $uibModalInstance.close($scope.text);
    };

    $scope.cancel = function() {
        $uibModalInstance.dismiss();
    };

}]);
