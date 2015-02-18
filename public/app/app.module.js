angular
.module('tima', ['ngRoute', 'ngSanitize', 'ngResource', 'ui.bootstrap', 'ui.bootstrap.showErrors', 'ui.select', 'jwt-decode'])
.config(['$routeProvider', '$httpProvider', 'uiSelectConfig', function($routeProvider, $httpProvider, uiSelectConfig) {

    uiSelectConfig.theme = 'bootstrap';

    var checkSignedIn = function(authService, $q, $timeout, $http, $location, $rootScope) {
        return authService.isSignedIn($q, $timeout, $http, $location, $rootScope);
    };

    $routeProvider
    .when('/signin', {
        templateUrl: 'app/signin/signin.html',
        controller: 'SigninController'
    })
    .when('/activities/:day', {
        templateUrl: 'app/activity/activityDay.html',
        controller: 'ActivityController',
        resolve: {
            signedIn: checkSignedIn
        }
    })
    .when('/projects', {
        templateUrl: 'app/project/projectList.html',
        controller: 'ProjectListController',
        resolve: {
            signedIn: checkSignedIn
        }
    })
    .when('/projects/:id', {
        templateUrl: 'app/project/project.html',
        controller: 'ProjectController',
        resolve: {
            signedIn: checkSignedIn
        }
    })
    .when('/projectCategories', {
        templateUrl: 'app/projectCategory/projectCategoryList.html',
        controller: 'ProjectCategoryListController',
        resolve: {
            signedIn: checkSignedIn
        }
    })
    .when('/activityTypes', {
        templateUrl: 'app/activityType/activityTypeList.html',
        controller: 'ActivityTypeListController',
        resolve: {
            signedIn: checkSignedIn
        }
    })
    .when('/activityTypes/:id', {
        templateUrl: 'app/activityType/activityType.html',
        controller: 'ActivityTypeController',
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

    $httpProvider.interceptors.push('httpAuthInterceptor');
    $httpProvider.interceptors.push('httpErrorInterceptor');

}]);
