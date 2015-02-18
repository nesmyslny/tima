angular.module('tima').controller('ProjectCategoryListController',
['$scope', '_', 'ProjectCategory', 'popupService',
function ($scope, _, ProjectCategory, popupService) {

    $scope.projectCategories = ProjectCategory.queryTree();
    var focusedCategoryId = null;

    $scope.focus = function(id) {
        focusedCategoryId = id;
    };

    $scope.hasFocus = function(id) {
        return focusedCategoryId === id;
    };

    $scope.showToggle = function(projectCategory) {
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

    $scope.delete = function(category) {
        var categories = $scope.projectCategories;
        var parent = findCategory(category.parentId, $scope.projectCategories);
        if (parent) {
            categories = parent.projectCategories;
        }

        popupService.show('Delete Project Category', 'Do you really want to delete this project category?', 'Delete', 'Cancel')
        .result.then(function() {
            ProjectCategory.delete({id:category.id}, function() {
                var index = categories.indexOf(category);
                categories.splice(index, 1);

                if (parent && categories.length === 0) {
                    parent.showChildren = false;
                }
            });
        });
    };

    function findCategory(id, categories) {
        if (id === null) {
            return null;
        }

        var category = null;

        _.forEach(categories, function(x) {
            if (x.id == id) {
                category = x;
                return false;
            }

            category = findCategory(id, x.projectCategories);
            if (category) {
                return false;
            }
        });

        return category;
    }

}]);
