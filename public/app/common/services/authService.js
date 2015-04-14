angular.module('tima').factory('authService',
['$location', '$q', 'JwtDecode', 'sessionService', 'Auth', 'userRoles',
function($location, $q, JwtDecode, sessionService, Auth, userRoles) {
    var service = {
        getUser: function() {
            return sessionService.user;
        },

        signIn: function(credentials, redirectPath) {
            Auth.signIn(credentials, function(data) {
                var tokenData = JwtDecode.decode(data.value);
                sessionService.init(data.value, tokenData.user);
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
                return data.value;
            }, function() {
                return false;
            });
        },

        isAuthorized: function(role) {
            return sessionService.user.role >= role.id;
        },

        isRole: function(role) {
            return sessionService.user.role === role.id;
        },

        isAdmin: function() {
            return service.isRole(userRoles.admin);
        },

        isManager: function() {
            return service.isRole(userRoles.manager);
        },

        isDeptManager: function() {
            return service.isRole(userRoles.deptManager);
        },

        isUser: function() {
            return service.isRole(userRoles.user);
        },

        checkPermission: function(role) {
            return service.isAuthenticated().then(function(authenticated) {
                if (!authenticated) {
                    service.signOut();
                    return $q.reject();
                }

                if (!service.isAuthorized(role)) {
                    $location.path("/");
                    return $q.reject();
                }
            });
        }
    };

    return service;
}]);
