angular.module('tima').factory('httpErrorInterceptor',
['$q', '$injector', 'messageService',
function ($q, $injector, messageService) {
    return {
        responseError: function (response) {
            if (response.status == 401) {
                // manually getting authService because of circular dependency ($http).
                // todo: investigate / refactor authService?
                authService = $injector.get('authService');
                authService.signOut();
            } else {
                messageService.add('danger', response.data);
            }
            return $q.reject(response);
        }
    };
}]);
