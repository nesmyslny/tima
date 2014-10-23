angular.module('gnomon').factory('authService', ['$http', '$window', '$location', function($http, $window, $location) {

    var signinPath = 'signin';

    return {

        signIn: function(credentials, redirectPath) {
            $http.post('/signin', credentials)
            .success(function(data, status, headers, config) {
                $window.sessionStorage.token = data.StringResult;
                $location.path(redirectPath);
                credentials.clear();
            })
            .error(function(data, status, headers, config) {
                delete $window.sessionStorage.token;
                credentials.clear();
            });
        },

        signOut: function() {
            delete $window.sessionStorage.token;
            $location.path(signinPath);
        },

        isSignedIn : function($q, $timeout, $http, $location, $rootScope){
            var deferred = $q.defer();
            $http.get('/issignedin')
            .success(function(data, status, headers, config) {
                if (data.BoolResult) {
                    $timeout(deferred.resolve, 0);
                } else {
                    $timeout(function(){deferred.reject();}, 0);
                    $location.url(signinPath);
                }
            })
            .error(function(data, status, headers, config) {
                $timeout(function(){deferred.reject();}, 0);
                $location.url(signinPath);
            });

            return deferred.promise;
        }
    };

}]);
