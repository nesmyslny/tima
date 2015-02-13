angular.module('tima').controller('ProjectCategoryListController',
['$scope', 'ProjectCategory', 'popupService',
function ($scope, ProjectCategory, popupService) {

    $scope.projectCategories = [];
    var focusedCategoryId = null;

    $scope.list = function() {
        $scope.projectCategories = ProjectCategory.query();
    };
    $scope.list();

    $scope.focus = function(id) {
        focusedCategoryId = id;
    };

    $scope.hasFocus = function(id) {
        return focusedCategoryId === id;
    };

    $scope.show = function(projectCategory) {
        projectCategory.isVisible = !projectCategory.isVisible;
    };

}]);
