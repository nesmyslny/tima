angular.module('tima').factory('popupService',
['$modal',
function($modal) {

    var service = {
        show: function(title, body, acceptButton, cancelButton) {
            return $modal.open({
                templateUrl: 'app/popup/popupTemplate.html',
                controller: 'PopupController',
                resolve: {
                    title: function() { return title; },
                    body: function() { return body; },
                    acceptButton: function() { return acceptButton; },
                    cancelButton: function() { return cancelButton; }
                }
            });
        },

        showForm: function(title, template, data, acceptButton, cancelButton) {
            return $modal.open({
                templateUrl: 'app/popup/popupFormTemplate.html',
                controller: 'PopupFormController',
                resolve: {
                    title: function() { return title; },
                    template: function() { return template; },
                    data: function() { return data; },
                    acceptButton: function() { return acceptButton; },
                    cancelButton: function() { return cancelButton; }
                }
            });
        }
    };

    return service;
}]);
