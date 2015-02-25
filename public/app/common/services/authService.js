angular.module('tima').factory('authService',
['$location', '$q', 'JwtDecode', 'sessionService', 'Auth',
function($location, $q, JwtDecode, sessionService, Auth) {
    var service = {
        getUser: function() {
            return sessionService.user;
        },

        signIn: function(credentials, redirectPath) {
            Auth.signIn(credentials, function(data) {
                var tokenData = JwtDecode.decode(data.stringResult);
                sessionService.init(data.stringResult, tokenData.user);
                $location.path(redirectPath);
                credentials.clear();
            }, function() {
                sessionService.delete();
                credentials.clear();
            });
        },

        signOut: function() {
            sessionService.delete();
            $location.path('signin');
        },

        isAuthenticated: function() {
            return Auth.isSignedIn().$promise.then(function(data) {
                return data.boolResult;
            }, function() {
                return false;
            });
        },

        isAuthorized: function(role) {
            return true;
        },

        checkPermission: function(role) {
            return service.isAuthenticated().then(function(authenticated) {
                if (!authenticated) {
                    service.signOut();
                    return $q.reject();
                }

                if (!service.isAuthorized(role)) {
                    return $q.reject();
                }
            });
        }
    };

    return service;
}]);
