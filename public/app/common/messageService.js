angular.module('tima').factory('messageService', ['$rootScope', '$timeout', function($rootScope, $timeout) {
    var service = {
        messages: [],

        add: function(type, text) {
            message = {
                type: type,
                text: text,
                close: function() {
                    service.remove(this);
                }
            };
            this.messages.push(message);
            message.timer = $timeout(function() { message.close(); }, 5000);
        },

        remove: function(message) {
            var index = this.messages.indexOf(message);
            this.messages.splice(index, 1);
        }
    };

    return service;
}]);
