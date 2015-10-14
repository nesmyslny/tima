angular.module('tima').factory('sessionService',
['$window', '_',
function($window, _) {
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

        updateUser: function(user) {
            // todo: is there a lodash-function for this? (copying only properties, which are in source.)
            _.forOwn(service.user, function(value, key, object) {
                object[key] = user[key];
            });
            $window.sessionStorage.user = JSON.stringify(service.user);
        },

        setToken: function(token) {
            $window.sessionStorage.token = service.token = token;
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
