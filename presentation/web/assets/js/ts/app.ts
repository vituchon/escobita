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

// TODO : move to presentation/web/assets/js/ts/directives and make proper inclusion
escobita.directive('loading', [() => {
  return {
    restrict: 'E',
    replace: true,
    scope: {
      message: '@?'
    },
    template:
    `<div class="verticalLayout center" style="opacity:0.7">
        <div class="loader">
          <div class="bounce1"></div>
          <div class="bounce2"></div>
          <div class="bounce3"></div>
        </div>
        <span style="font-size:18px;font-weight:200" ng-show="message">{{message}}</span>
     </div>`
  }
}]);