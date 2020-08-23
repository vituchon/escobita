/// <reference path='./third_party_definitions/_definitions.ts' />
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
                game: null
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
