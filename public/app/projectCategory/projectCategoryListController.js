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
        var refIdPrefix = "";
        var categories = $scope.projectCategories;
        if (parent) {
            parentId = parent.id;
            refIdPrefix = parent.refIdComplete;
            categories = parent.projectCategories;
        }

        var category = {
            id: -1,
            parentId: parentId,
            refId: "",
            title: "",
            refIdPrefix: refIdPrefix,
            projectCategories: []
        };

        popupService.showForm("Add Project Category", "app/projectCategory/projectCategoryPopupTemplate.html", category, "Add", "Cancel")
        .result.then(function() {

            ProjectCategory.save(category, function(response) {
                category.id = response.intResult;
                category.refIdComplete = category.refIdPrefix + category.refId;
                categories.push(category);

                if (parent) {
                    parent.showChildren = true;
                }

            });
        });
    };

    $scope.edit = function(category) {
        var data = _.clone(category);
        parent = findCategory(category.parentId, $scope.projectCategories);
        if (parent) {
            data.refIdPrefix = parent.refIdComplete;
        }

        popupService.showForm("Edit Project Category", "app/projectCategory/projectCategoryPopupTemplate.html", data, "Save", "Cancel")
        .result.then(function() {
            ProjectCategory.save(data, function() {
                _.assign(category, data);
            });
        });
    };

    $scope.delete = function(category) {
        var categories = $scope.projectCategories;
        var parent = findCategory(category.parentId, $scope.projectCategories);
        if (parent) {
            categories = parent.projectCategories;
        }

        popupService.showSimple('Delete Project Category', 'Do you really want to delete this project category?', 'Delete', 'Cancel')
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
