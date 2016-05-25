angular.module('tima').factory('multiSelect',
['_', 'popupService',
function (_, popupService) {

    function addItem(selectedItems, item) {
        if (_.isUndefined(item)) {
            return;
        }

        var alreadyInList = selectedItems.some(function(selectedItem) {
            return item.id == selectedItem.id;
        });

        if (!alreadyInList) {
            selectedItems.push(item);
        }
    }

    return {
        addSelectedItem: function(selectedItems, selectedItem) {
            addItem(selectedItems, selectedItem.selected);
            selectedItem.selected = undefined;
        },

        addMultipleItems: function(items, selectedItems, popupTitle, valueKey1, valueClass1, valueKey2, valueClass2) {
            var popupItems = [];
            _.forEach(items, function(item) {
                popupItems.push({
                    value1: item[valueKey1],
                    class1: valueClass1,
                    value2: item[valueKey2],
                    class2: valueClass2,
                    obj: item,
                    checked: _.some(selectedItems, ["id", item.id])
                });
            });

            popupService.showSelectList(popupTitle, popupItems, "Ok", "Cancel")
            .result.then(function() {
                selectedItems.length = 0;
                _.forEach(popupItems, function(item) {
                    if (item.checked) {
                        addItem(selectedItems, item.obj);
                    }
                });
            });
        },

        removeItem: function(selectedItems, item) {
            var index = selectedItems.indexOf(item);
            selectedItems.splice(index, 1);
        }
    };
}]);
