angular
.module('tima', ['ngRoute', 'ngSanitize', 'ngResource', 'ui.bootstrap', 'ui.bootstrap.showErrors', 'validation.match', 'ui.select', 'jwt-decode', 'angular-momentjs', 'hc.marked', 'chart.js'])
.constant('_', window._) // use lodash via DI in controllers, etc.
.constant('userRoles', {
    user: {id: 10, desc: "User"},
    deptManager: {id: 30, desc: "Department Manager"},
    manager: {id: 50, desc: "Manager"},
    admin: {id: 99, desc: "Administrator"}
})
.config(['$routeProvider', '$httpProvider', 'uiSelectConfig', 'markedProvider', 'ChartJsProvider', 'userRoles', function($routeProvider, $httpProvider, uiSelectConfig, markedProvider, ChartJsProvider, userRoles) {

    uiSelectConfig.theme = 'bootstrap';
    markedProvider.setOptions({
        sanitize: true
    });

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
        resolve: createPermissionResolve(userRoles.user)
    })
    .when('/myprojects', {
        templateUrl: 'app/project/myProjectList.html',
        controller: 'MyProjectListController',
        resolve: createPermissionResolve(userRoles.user)
    })
    .when('/projects', {
        templateUrl: 'app/project/projectList.html',
        controller: 'ProjectListController',
        resolve: createPermissionResolve(userRoles.deptManager)
    })
    .when('/projects/:id', {
        templateUrl: 'app/project/project.html',
        controller: 'ProjectController',
        resolve: createPermissionResolve(userRoles.user)
    })
    .when('/projectCategories', {
        templateUrl: 'app/projectCategory/projectCategoryList.html',
        controller: 'ProjectCategoryListController',
        resolve: createPermissionResolve(userRoles.manager)
    })
    .when('/activityTypes', {
        templateUrl: 'app/activityType/activityTypeList.html',
        controller: 'ActivityTypeListController',
        resolve: createPermissionResolve(userRoles.manager)
    })
    .when('/users', {
        templateUrl: 'app/user/userList.html',
        controller: 'UserListController',
        resolve: createPermissionResolve(userRoles.admin)
    })
    .when('/users/:id', {
        templateUrl: 'app/user/userAdministration.html',
        controller: 'UserController',
        resolve: createPermissionResolve(userRoles.admin)
    })
    .when('/userSettings', {
        templateUrl: 'app/user/userSettings.html',
        controller: 'UserController',
        resolve: createPermissionResolve(userRoles.user)
    })
    .when('/departments', {
        templateUrl: 'app/department/departmentList.html',
        controller: 'DepartmentListController',
        resolve: createPermissionResolve(userRoles.admin)
    })
    .when('/reporting', {
        templateUrl: 'app/reporting/reporting.html',
        controller: 'ReportingController',
        resolve: createPermissionResolve(userRoles.deptManager)
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
