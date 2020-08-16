/// <reference path='./third_party_definitions/_definitions.ts' />
var escobita = angular.module('escobita', ['ui.router']);
var App;
(function (App) {
    function setup($state, $urlRouterProvider, $location) {
        $location.html5Mode({ enabled: true, requireBase: false });
        $urlRouterProvider.otherwise('/');
        /*const landing : ng.ui.IState = {
          name: 'landing',
          url: '/',
          templateUrl: '/presentation/web/assets/html/landing.html',
          controller: "LandingController",
          controllerAs: "ctr",
        };*/
        var lobby = {
            name: 'lobby',
            url: 'lobby',
            templateUrl: '/presentation/web/assets/html/lobby.html',
            controller: "LobbyController",
            controllerAs: "ctr"
        };
        var about = {
            name: 'about',
            url: 'about',
            templateUrl: '/presentation/web/assets/html/about.html',
            controller: "AboutController",
            controllerAs: "ctr"
        };
        //$state.state(landing);
        $state.state(lobby);
        $state.state(about);
    }
    ;
    escobita.config(['$stateProvider', '$urlRouterProvider', '$locationProvider', setup]);
})(App || (App = {}));
