angular.module('tima').factory('sessionService',
['$window',
function($window) {
    function getUserOrNull() {
        if ($window.sessionStorage.user) {
            return JSON.parse($window.sessionStorage.user);
        } else {
            return null;
        }
    }

    var service = {
        user: getUserOrNull(),
        token: $window.sessionStorage.token,

        init: function(token, user) {
            $window.sessionStorage.token = service.token = token;
            $window.sessionStorage.user = JSON.stringify(user);
            service.user = user;
        },

        delete: function() {
            service.user = null;
            service.token = '';
            delete $window.sessionStorage.token;
            delete $window.sessionStorage.user;
        }
    };

    return service;
}]);
