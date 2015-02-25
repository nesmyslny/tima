angular.module('tima').factory('Auth',
['$resource',
function($resource) {
    return $resource("", {}, {
        signIn: {
            url: "/signIn",
            method: "POST"
        },
        isSignedIn: {
            url: "isSignedIn",
            method: "GET"
        }
    });
}]);
