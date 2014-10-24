angular.module('tima').factory('authInterceptor', ['$rootScope', '$q', '$window', '$location', '$injector', function ($rootScope, $q, $window, $location, $injector) {
    return {
        request: function (config) {
            config.headers = config.headers || {};
            if ($window.sessionStorage.token) {
                config.headers.Authorization = 'Bearer ' + $window.sessionStorage.token;
            }
            return config;
        },
        responseError: function (response) {
            if (response.status == 401) {
                // manuelly getting authService because of circular dependency ($http).
                // todo: investigate / refactor authService?
                authService = $injector.get('authService');
                authService.signOut();
            }
            return $q.reject(response);
        }
    };
}]);
