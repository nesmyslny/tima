angular.module('tima').controller('PopupMarkdownController',
['$scope', '$modalInstance', 'title', 'markdown', 'acceptButton', 'cancelButton',
function ($scope, $modalInstance, title, markdown, acceptButton, cancelButton) {

    $scope.title = title;
    $scope.acceptButton = acceptButton;
    $scope.cancelButton = cancelButton;
    // make an object for data used in tabs, so it will be in the tab scope as a reference.
    $scope.tabsetData = {
        markdown: markdown
    };

    $scope.accept = function() {
        $modalInstance.close($scope.tabsetData.markdown);
    };

    $scope.cancel = function() {
        $modalInstance.dismiss();
    };

}]);
