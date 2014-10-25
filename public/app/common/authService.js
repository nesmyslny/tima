angular.module('tima').factory('authService', ['$http', '$location', 'sessionService', function($http, $location, sessionService) {

    var service = {

        getUsername: function() {
            return sessionService.getUsername();
        },

        signIn: function(credentials, redirectPath) {
            $http.post('/signin', credentials)
            .success(function(data, status, headers, config) {
                var tokenData = jwt_decode(data.StringResult);
                sessionService.init(data.StringResult, tokenData.username);
                $location.path(redirectPath);
                credentials.clear();
            })
            .error(function(data, status, headers, config) {
                sessionService.delete();
                credentials.clear();
            });
        },

        signOut: function() {
            sessionService.delete();
            $location.path('signin');
        },

        isSignedIn : function($q, $timeout, $http, $location, $rootScope){
            var deferred = $q.defer();
            $http.get('/issignedin')
            .success(function(data, status, headers, config) {
                if (data.BoolResult) {
                    $timeout(deferred.resolve, 0);
                } else {
                    $timeout(function(){deferred.reject();}, 0);
                    service.signOut();
                }
            })
            .error(function(data, status, headers, config) {
                $timeout(function(){deferred.reject();}, 0);
                service.signOut();
            });

            return deferred.promise;
        }

    };

    return service;
}]);
