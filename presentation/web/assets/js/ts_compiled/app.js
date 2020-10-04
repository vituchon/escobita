/// <reference path='./third_party_definitions/_definitions.ts' />
/// <reference path='./util.ts' />
var escobita = angular.module('escobita', ['ui.router']);
var App;
(function (App) {
    function setup($state, $urlRouterProvider, $location) {
        $location.html5Mode({ enabled: true, requireBase: false });
        $urlRouterProvider.otherwise('/');
        var lobby = {
            name: 'lobby',
            url: 'lobby',
            templateUrl: '/presentation/web/assets/html/lobby.html',
            controller: "LobbyController",
            controllerAs: "ctr"
        };
        var game = {
            name: 'game',
            url: 'game',
            templateUrl: '/presentation/web/assets/html/game.html',
            controller: "GameController",
            controllerAs: "ctr",
            params: {
                game: null,
                player: null
            }
        };
        var about = {
            name: 'about',
            url: 'about',
            templateUrl: '/presentation/web/assets/html/about.html',
            controller: "AboutController",
            controllerAs: "ctr"
        };
        $state.state(lobby);
        $state.state(about);
        $state.state(game);
    }
    ;
    escobita.config(['$stateProvider', '$urlRouterProvider', '$locationProvider', setup]);
})(App || (App = {}));
escobita.run(['$state', function ($state) {
        $state.go("lobby");
        $("#vitu-logo").animate({ height: 140 }, 2000);
    }]);
escobita.directive('loading', [function () {
        return {
            restrict: 'E',
            replace: true,
            scope: {
                message: '@?'
            },
            template: "<div class=\"verticalLayout center\" style=\"opacity:0.7\">\n        <div class=\"loader\">\n          <div class=\"bounce1\"></div>\n          <div class=\"bounce2\"></div>\n          <div class=\"bounce3\"></div>\n        </div>\n        <span style=\"font-size:18px;font-weight:200\" ng-show=\"message\">{{message}}</span>\n     </div>"
        };
    }]);
// Wraps toastr calls using site "custom look and feel" parameters
var Toastr;
(function (Toastr) {
    function success(message) {
        return toastr.success(message, '', {
            positionClass: 'toast-bottom-left',
            toastClass: "toast-container"
        });
    }
    Toastr.success = success;
    function info(message) {
        return toastr.info(message, '', {
            positionClass: 'toast-bottom-left',
            toastClass: "toast-container"
        });
    }
    Toastr.info = info;
    function warn(message) {
        return toastr.warning(message, '', {
            positionClass: 'toast-bottom-left',
            toastClass: "toast-container"
        });
    }
    Toastr.warn = warn;
    function error(message) {
        return toastr.error(message, '', {
            positionClass: 'toast-bottom-left',
            toastClass: "toast-container"
        });
    }
    Toastr.error = error;
    function clear() {
        return toastr.clear();
    }
    Toastr.clear = clear;
    function chat(playerName, message) {
        return toastr.info(message, "De " + playerName, {
            positionClass: 'toast-bottom-left',
            toastClass: "toastr-chat-class",
            titleClass: "toasrt-chat-tittle",
            messageClass: "toasrt-chat-message",
            timeOut: 10000
        });
    }
    Toastr.chat = chat;
})(Toastr || (Toastr = {}));
