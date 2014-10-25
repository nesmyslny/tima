angular.module('tima').factory('sessionService', ['$window', function($window) {

    storage = {
        token: $window.sessionStorage.token,
        username: $window.sessionStorage.username,
        set: function(token, username) {
            this.token = token;
            this.username = username;
        }
    };

    return {

        getToken: function() {
            return storage.token;
        },

        getUsername: function() {
            return storage.username;
        },

        init: function(token, username) {
            storage.set(token, username);
            $window.sessionStorage.token = token;
            $window.sessionStorage.username = username;
        },

        delete: function() {
            storage.set('', '');
            delete $window.sessionStorage.token;
            delete $window.sessionStorage.username;
        }

    };
}]);
