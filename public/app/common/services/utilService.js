angular.module('tima').factory('util',
[
function () {
    return {
        formatTime: function(hours, minutes) {
            var durationFormatted = hours > 0 ? hours + 'h' : '';
            durationFormatted += minutes > 0 ? ' ' + minutes + 'min' : '';
            return durationFormatted;
        }
    };
}]);
