angular.module('tima').controller('DepartmentListController',
['$scope', '_', 'Department', 'popupService',
function ($scope, _, Department, popupService) {

    $scope.departments = Department.queryTree();
    var focusedDepartmentId = null;

    $scope.focus = function(id) {
        focusedDepartmentId = id;
    };

    $scope.hasFocus = function(id) {
        return focusedDepartmentId === id;
    };

    $scope.showToggle = function(department) {
        department.showChildren = !department.showChildren;
    };

    $scope.add = function(parent) {
        var parentId = null;
        var refIdPrefix = "";
        var departments = $scope.departments;
        if (parent) {
            parentId = parent.id;
            departments = parent.departments;
        }

        var department = {
            id: -1,
            parentId: parentId,
            title: "",
            departments: []
        };

        popupService.showForm("Add Department", "app/department/departmentPopupTemplate.html", department, "Add", "Cancel")
        .result.then(function() {

            Department.save(department, function(response) {
                department.id = response.intResult;
                departments.push(department);

                if (parent) {
                    parent.showChildren = true;
                }

            });
        });
    };

    $scope.edit = function(department) {
        var data = _.clone(department);

        popupService.showForm("Edit Department", "app/department/departmentPopupTemplate.html", data, "Save", "Cancel")
        .result.then(function() {
            Department.save(data, function() {
                _.assign(department, data);
            });
        });
    };

    $scope.delete = function(department) {
        var departments = $scope.departments;
        var parent = findDepartment(department.parentId, $scope.departments);
        if (parent) {
            departments = parent.departments;
        }

        popupService.showSimple("Delete Department", "Do you really want to delete this project category?", "Delete", "Cancel")
        .result.then(function() {
            Department.delete({id:department.id}, function() {
                var index = departments.indexOf(department);
                departments.splice(index, 1);

                if (parent && departments.length === 0) {
                    parent.showChildren = false;
                }
            });
        });
    };

    function findDepartment(id, departments) {
        if (id === null) {
            return null;
        }

        var department = null;

        _.forEach(departments, function(x) {
            if (x.id == id) {
                department = x;
                return false;
            }

            department = findDepartment(id, x.departments);
            if (department) {
                return false;
            }
        });

        return department;
    }

}]);
