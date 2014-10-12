angular
.module('gnomon', ['ngRoute', 'ui.bootstrap.showErrors'])
.config(['$routeProvider', function($routeProvider) {

    $routeProvider
    .when('/', {
        templateUrl: 'app/components/signin/signin.html',
        controller: 'SigninController'
    })
    .otherwise({
        redirectTo: '/'
    });
}]);
