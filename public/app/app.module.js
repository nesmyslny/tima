angular
.module('tima', ['ngRoute', 'ui.bootstrap', 'ui.bootstrap.showErrors'])
.config(['$routeProvider', '$httpProvider', function($routeProvider, $httpProvider) {

    var checkSignedIn = function(authService, $q, $timeout, $http, $location, $rootScope) {
        return authService.isSignedIn($q, $timeout, $http, $location, $rootScope);
    };

    $routeProvider
    .when('/signin', {
        templateUrl: 'app/signin/signin.html',
        controller: 'signinController'
    })
    .when('/activities/:day', {
        templateUrl: 'app/activities/activities.html',
        controller: 'activitiesController',
        resolve: {
            signedIn: checkSignedIn
        }
    })
    .when('/', {
        redirectTo: '/activities/' + moment().format('YYYY-MM-DD')
    })
    .otherwise({
        redirectTo: '/signin'
    });

    $httpProvider.interceptors.push('authInterceptor');

}]);
