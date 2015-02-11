angular
.module('tima', ['ngRoute', 'ui.bootstrap', 'ui.bootstrap.showErrors', 'ngSanitize', 'ui.select'])
.config(['$routeProvider', '$httpProvider', 'uiSelectConfig', function($routeProvider, $httpProvider, uiSelectConfig) {

    uiSelectConfig.theme = 'bootstrap';

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
    .when('/projects', {
        templateUrl: 'app/projects/projectList.html',
        controller: 'projectListController',
        resolve: {
            signedIn: checkSignedIn
        }
    })
    .when('/projects/:id', {
        templateUrl: 'app/projects/project.html',
        controller: 'projectController',
        resolve: {
            signedIn: checkSignedIn
        }
    })
    .when('/activityTypes', {
        templateUrl: 'app/activityTypes/activityTypeList.html',
        controller: 'activityTypeListController',
        resolve: {
            signedIn: checkSignedIn
        }
    })
    .when('/activityTypes/:id', {
        templateUrl: 'app/activityTypes/activityType.html',
        controller: 'activityTypeController',
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
