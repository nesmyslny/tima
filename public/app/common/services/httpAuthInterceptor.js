angular.module('tima').factory('httpAuthInterceptor',
['sessionService',
function (sessionService) {
    return {
        request: function (config) {
            config.headers = config.headers || {};
            if (sessionService.token) {
                config.headers.Authorization = 'Bearer ' + sessionService.token;
            }
            return config;
        }
    };
}]);
