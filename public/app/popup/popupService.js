angular.module('tima').factory('popupService',
['$uibModal',
function($uibModal) {

    var service = {
        showSimple: function(title, body, acceptButton, cancelButton) {
            return $uibModal.open({
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
            return $uibModal.open({
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
            return $uibModal.open({
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
            return $uibModal.open({
                templateUrl: 'app/popup/popupMarkdownTemplate.html',
                controller: 'PopupMarkdownController',
                resolve: {
                    title: function() { return title; },
                    markdown: function() { return markdown; },
                    acceptButton: function() { return acceptButton; },
                    cancelButton: function() { return cancelButton; }
                }
            });
        },

        showText: function(title, text, acceptButton, cancelButton) {
            return $uibModal.open({
                templateUrl: 'app/popup/popupTextTemplate.html',
                controller: 'PopupTextController',
                resolve: {
                    title: function() { return title; },
                    text: function() { return text; },
                    acceptButton: function() { return acceptButton; },
                    cancelButton: function() { return cancelButton; }
                }
            });
        }
    };

    return service;
}]);
