angular
.module('tima', ['ngRoute', 'ngSanitize', 'ngResource', 'ui.bootstrap', 'ui.bootstrap.showErrors', 'validation.match', 'ui.select', 'jwt-decode', 'angular-momentjs'])
.constant('_', window._) // use lodash via DI in controllers, etc.
.config(['$routeProvider', '$httpProvider', 'uiSelectConfig', function($routeProvider, $httpProvider, uiSelectConfig) {

    uiSelectConfig.theme = 'bootstrap';

    function createPermissionResolve(role) {
        return {
            auth: function(authService) {
                return authService.checkPermission(role);
            }
        };
    }

    $routeProvider
    .when('/signin', {
        templateUrl: 'app/signin/signin.html',
        controller: 'SigninController'
    })
    .when('/activities/:day', {
        templateUrl: 'app/activity/activityDay.html',
        controller: 'ActivityController',
        resolve: createPermissionResolve(0)
    })
    .when('/projects', {
        templateUrl: 'app/project/projectList.html',
        controller: 'ProjectListController',
        resolve: createPermissionResolve(0)
    })
    .when('/projects/:id', {
        templateUrl: 'app/project/project.html',
        controller: 'ProjectController',
        resolve: createPermissionResolve(0)
    })
    .when('/projectCategories', {
        templateUrl: 'app/projectCategory/projectCategoryList.html',
        controller: 'ProjectCategoryListController',
        resolve: createPermissionResolve(0)
    })
    .when('/activityTypes', {
        templateUrl: 'app/activityType/activityTypeList.html',
        controller: 'ActivityTypeListController',
        resolve: createPermissionResolve(0)
    })
    .when('/users', {
        templateUrl: 'app/user/userList.html',
        controller: 'UserListController',
        resolve: createPermissionResolve(0)
    })
    .when('/users/:id', {
        templateUrl: 'app/user/userAdministration.html',
        controller: 'UserController',
        resolve: createPermissionResolve(0)
    })
    .when('/userSettings', {
        templateUrl: 'app/user/userSettings.html',
        controller: 'UserController',
        resolve: createPermissionResolve(0)
    })
    .when('/', {
        redirectTo: '/activities/' + moment().format('YYYY-MM-DD')
    })
    .otherwise({
        redirectTo: '/activities/' + moment().format('YYYY-MM-DD')
    });

    $httpProvider.interceptors.push('httpAuthInterceptor');
    $httpProvider.interceptors.push('httpErrorInterceptor');

}])
.run(['$rootScope', function($rootScope) {
    // use lodash in views
    $rootScope._ = window._;
}]);
