angular.module('tima').controller('PopupFormController',
['$scope', '$modalInstance', 'title', 'template', 'data', 'acceptButton', 'cancelButton',
function ($scope, $modalInstance, title, template, data, acceptButton, cancelButton) {

    $scope.title = title;
    $scope.template = template;
    $scope.data = data;
    $scope.acceptButton = acceptButton;
    $scope.cancelButton = cancelButton;

    $scope.accept = function () {
        $scope.$broadcast('show-errors-check-validity');
        if (!$scope.form.$valid) {
            return;
        }

        $modalInstance.close(true);
    };

    $scope.cancel = function () {
        $modalInstance.dismiss();
    };

}]);
