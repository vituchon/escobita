/// <reference path='./third_party_definitions/_definitions.ts' />

const escobita: ng.IModule = angular.module('escobita', ['ui.router']);

module App {

  function setup($state: angular.ui.IStateProvider, $urlRouterProvider: angular.ui.IUrlRouterProvider, $location: angular.ILocationProvider) {
    $location.html5Mode({ enabled: true, requireBase: false });
    $urlRouterProvider.otherwise('/');

    const lobby: ng.ui.IState = {
      name: 'lobby',
      url: 'lobby',
      templateUrl: '/presentation/web/assets/html/lobby.html',
      controller: "LobbyController",
      controllerAs: "ctr"
    };

    const game: ng.ui.IState = {
      name: 'game',
      url: 'game',
      templateUrl: '/presentation/web/assets/html/game.html',
      controller: "GameController",
      controllerAs: "ctr",
      params: {
        game: null,
        player: null,
      }
    };

    const about: ng.ui.IState = {
      name: 'about',
      url: 'about',
      templateUrl: '/presentation/web/assets/html/about.html',
      controller: "AboutController",
      controllerAs: "ctr"
    };

    $state.state(lobby);
    $state.state(about);
    $state.state(game);
  };

  escobita.config(['$stateProvider', '$urlRouterProvider', '$locationProvider', setup]);
}