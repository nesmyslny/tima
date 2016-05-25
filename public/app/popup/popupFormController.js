angular.module('tima').controller('PopupFormController',
['$scope', '$uibModalInstance', 'title', 'template', 'data', 'acceptButton', 'cancelButton',
function ($scope, $uibModalInstance, title, template, data, acceptButton, cancelButton) {

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

        $uibModalInstance.close(true);
    };

    $scope.cancel = function () {
        $uibModalInstance.dismiss();
    };

}]);
