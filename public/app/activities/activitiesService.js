angular.module('tima').factory('activitiesService', ['$http', '$q', '$filter', function($http, $q, $filter) {
    var service = {

        refreshActivities: function(day, activities) {
            var deferred = $q.defer();

            $http.get('/activities/' + day)
            .success(function(data, status, headers, config) {
                var totalDuration = 0;

                data.forEach(function(activity) {
                    var actFound = $filter('filter')(activities, {id: activity.id}, true);
                    if (actFound.length) {
                        actFound[0].duration = activity.duration;
                    } else {
                        activities.push(activity);
                    }
                    totalDuration += activity.duration;
                });

                deferred.resolve({
                    totalDuration: totalDuration
                });
            })
            .error(function(data, status, header, config){
                // todo: error handling
                deferred.reject(data, status);
            });

            return deferred.promise;
        },

        saveActivity: function(day, text, hours, minutes) {
            var deferred = $q.defer();
            var activity = {
                day: moment(day, 'YYYY-MM-DD').format('YYYY-MM-DD[T]00:00:00.000[Z]'),
                text: text,
                duration: hours * 60 + minutes
            };

            $http.post('/activities', activity)
            .success(function(data, status, headers, config) {
                deferred.resolve();
            })
            .error(function(data, status, header, config) {
                // todo: error handling
                deferred.reject(data, status);
            });

            return deferred.promise;
        }

    };

    return service;
}]);
