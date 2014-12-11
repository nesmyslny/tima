angular.module('tima').directive('timaMessages', ['$rootScope', 'messageService', function ($rootScope, messageService) {
    var templateString = '<alert type="{{message.type}}" ng-repeat="message in timaMessages" close="message.close()">{{message.text}}</alert>';

    return {
        restrict: 'EA',
        template: templateString,
        link: function(scope, element, attrs) {
            $rootScope.timaMessages = messageService.messages;
        }
    };
}]);
