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
        queryMyProjects: {
            method: "GET",
            url: "/myprojects",
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
