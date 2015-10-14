angular.module('tima').controller('PopupTextController',
['$scope', '$modalInstance', 'title', 'text', 'acceptButton', 'cancelButton',
function ($scope, $modalInstance, title, text, acceptButton, cancelButton) {

    $scope.title = title;
    $scope.text = text;
    $scope.acceptButton = acceptButton;
    $scope.cancelButton = cancelButton;

    $scope.accept = function() {
        $modalInstance.close($scope.text);
    };

    $scope.cancel = function() {
        $modalInstance.dismiss();
    };

}]);
