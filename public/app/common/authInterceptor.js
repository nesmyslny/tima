angular.module('tima').factory('authInterceptor', ['$q', '$injector', 'sessionService', function ($q, $injector, sessionService) {
    return {
        request: function (config) {
            config.headers = config.headers || {};
            if (sessionService.getToken()) {
                config.headers.Authorization = 'Bearer ' + sessionService.getToken();
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
