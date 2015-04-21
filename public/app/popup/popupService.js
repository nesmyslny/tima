angular.module('tima').factory('popupService',
['$modal',
function($modal) {

    var service = {
        showSimple: function(title, body, acceptButton, cancelButton) {
            return $modal.open({
                templateUrl: 'app/popup/popupSimpleTemplate.html',
                controller: 'PopupSimpleController',
                resolve: {
                    title: function() { return title; },
                    body: function() { return body; },
                    acceptButton: function() { return acceptButton; },
                    cancelButton: function() { return cancelButton; }
                }
            });
        },

        showForm: function(title, template, data, acceptButton, cancelButton) {
            return $modal.open({
                templateUrl: 'app/popup/popupFormTemplate.html',
                controller: 'PopupFormController',
                resolve: {
                    title: function() { return title; },
                    template: function() { return template; },
                    data: function() { return data; },
                    acceptButton: function() { return acceptButton; },
                    cancelButton: function() { return cancelButton; }
                }
            });
        },

        showSelectList: function(title, items, acceptButton, cancelButton) {
            return $modal.open({
                templateUrl: 'app/popup/popupSelectListTemplate.html',
                controller: 'PopupSelectListController',
                resolve: {
                    title: function() { return title; },
                    items: function() { return items; },
                    acceptButton: function() { return acceptButton; },
                    cancelButton: function() { return cancelButton; }
                }
            });
        },

        showMarkdown: function(title, markdown, acceptButton, cancelButton) {
            return $modal.open({
                templateUrl: 'app/popup/popupMarkdownTemplate.html',
                controller: 'PopupMarkdownController',
                resolve: {
                    title: function() { return title; },
                    markdown: function() { return markdown; },
                    acceptButton: function() { return acceptButton; },
                    cancelButton: function() { return cancelButton; }
                }
            });
        }
    };

    return service;
}]);
