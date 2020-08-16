/// <reference path='./third_party_definitions/_definitions.ts' />

const escobita: ng.IModule = angular.module('escobita', ['ui.router']);

module App {

  function setup($state: angular.ui.IStateProvider, $urlRouterProvider: angular.ui.IUrlRouterProvider, $location: angular.ILocationProvider) {
    $location.html5Mode({ enabled: true, requireBase: false });
    $urlRouterProvider.otherwise('/');

    /*const landing : ng.ui.IState = {
      name: 'landing',
      url: '/',
      templateUrl: '/presentation/web/assets/html/landing.html',
      controller: "LandingController",
      controllerAs: "ctr",
    };*/

    const lobby: ng.ui.IState = {
      name: 'lobby',
      url: 'lobby',
      templateUrl: '/presentation/web/assets/html/lobby.html',
      controller: "LobbyController",
      controllerAs: "ctr"
    };

    const about: ng.ui.IState = {
      name: 'about',
      url: 'about',
      templateUrl: '/presentation/web/assets/html/about.html',
      controller: "AboutController",
      controllerAs: "ctr"
    };

    //$state.state(landing);
    $state.state(lobby);
    $state.state(about);
  };

  escobita.config(['$stateProvider', '$urlRouterProvider', '$locationProvider', setup]);
}