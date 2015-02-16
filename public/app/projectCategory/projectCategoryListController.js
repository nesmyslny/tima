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
        projectCategory.showChildren = !projectCategory.showChildren;
    };

    $scope.add = function(parent) {
        var parentId = null;
        var categories = $scope.projectCategories;
        if (parent) {
            parentId = parent.id;
            categories = parent.projectCategories;
        }

        data = {
            title: ""
        };

        popupService.showForm("Add Project Category", "app/projectCategory/projectCategoryPopupTemplate.html", data, "Add", "Cancel")
        .result.then(function() {
            var category = {
                id: -1,
                parentId: parentId,
                title: data.title,
                projectCategories: []
            };

            ProjectCategory.save(category, function(response) {
                category.id = response.intResult;
                categories.push(category);

                if (parent) {
                    parent.showChildren = true;
                }

            });
        });
    };

    $scope.edit = function(category) {
        var data = {
            id: category.id,
            parentId: category.parentId,
            title: category.title
        };

        popupService.showForm("Edit Project Category", "app/projectCategory/projectCategoryPopupTemplate.html", data, "Save", "Cancel")
        .result.then(function() {
            ProjectCategory.save(data, function() {
                category.title = data.title;
            });
        });
    };

}]);
