angular.module('tima').directive('userAccess',
['_', 'authService', 'userRoles',
function(_, authService, userRoles) {
    return {
        restrict: 'A',
        link: function(scope, element, attrs) {
            var makeVisible = function() {
                element.removeClass('hidden');
            };

            var makeHidden = function() {
                element.addClass('hidden');
            };

            var determineVisibility = function() {
                if (authService.isAuthorized(role)) {
                    makeVisible();
                } else {
                    makeHidden();
                }
            };

            var role = userRoles[attrs.userAccess];

            if (!_.isUndefined(role)) {
                determineVisibility();
            }
        }
    };
}]);
