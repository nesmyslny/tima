angular
.module('gnomon', ['ngRoute', 'ui.bootstrap.showErrors'])
.config(['$routeProvider', '$httpProvider', function($routeProvider, $httpProvider) {

    $routeProvider
    .when('/signin', {
        templateUrl: 'app/components/signin/signin.html',
        controller: 'SigninController'
    })
    .when('/', {
        templateUrl: 'app/components/secret/secret.html',
        controller: 'SigninController',
        resolve: {
            loggedin: function(authService, $q, $timeout, $http, $location, $rootScope) {
                return authService.isSignedIn($q, $timeout, $http, $location, $rootScope);
            }
        }
    })
    .otherwise({
        redirectTo: '/signin'
    });

    $httpProvider.interceptors.push('authInterceptor');
}]);