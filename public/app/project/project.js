angular.module('tima').factory('Project',
['$resource', 'resourceSaveInterceptor', 'sessionService', 'util',
function($resource, resourceSaveInterceptor, sessionService, util) {

    function setProjectUserFlags(project) {
        project.isResponsible = project.responsibleUserId === sessionService.user.id;
        project.isManager = project.managerUserId === sessionService.user.id;
    }

    return $resource("/projects/:id", {}, {
        save: {
            method: "POST",
            interceptor: resourceSaveInterceptor
        },
        queryAdmin: {
            method: "GET",
            url: "/projects/admin",
            isArray: true
        },
        queryUser: {
            method: "GET",
            url: "/projects/user",
            isArray: true,
            transformResponse: function(data, headers) {
                if (util.isJsonResponse(headers)) {
                    data = angular.fromJson(data);
                    _.forEach(data, function(project) {
                        setProjectUserFlags(project);
                    });
                }

                return data;
            }
        },
        querySelectList: {
            method: "GET",
            url: "/projects/selectList",
            isArray: true
        },
        get: {
            method: "GET",
            url: "/projects/:id",
            transformResponse: function(data, headers) {
                if (util.isJsonResponse(headers)) {
                    data = angular.fromJson(data);
                    setProjectUserFlags(data);
                }

                return data;
            }
        }
    });
}]);
