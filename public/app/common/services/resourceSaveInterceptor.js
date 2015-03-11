angular.module('tima').factory('resourceSaveInterceptor',
['_',
function (_) {
    return {
        response: function (response) {
            if (response.config.method == "POST" && response.status == 200) {
                _.assign(response.config.data, response.data);
            }
            return response;
        }
    };
}]);
