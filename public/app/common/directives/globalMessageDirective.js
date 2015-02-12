angular.module('tima').directive('globalMessages',
['$rootScope', 'messageService',
function ($rootScope, messageService) {
    var templateString = '<alert type="{{message.type}}" ng-repeat="message in globalMessages" close="message.close()">{{message.text}}</alert>';

    return {
        restrict: 'EA',
        template: templateString,
        link: function(scope, element, attrs) {
            $rootScope.globalMessages = messageService.messages;
        }
    };
}]);
