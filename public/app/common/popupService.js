angular.module('tima').factory('popupService', ['$modal', function($modal) {

    var service = {
        show: function(title, body, acceptButton, cancelButton, acceptFunc, cancelFunc) {
            return $modal.open({
                templateUrl: 'app/common/popupTemplate.html',
                controller: 'popupController',
                resolve: {
                    title: function() { return title; },
                    body: function() { return body },
                    acceptButton: function() { return acceptButton; },
                    cancelButton: function() { return cancelButton; }
                }
            });
        }
    };

    return service;
}]);
