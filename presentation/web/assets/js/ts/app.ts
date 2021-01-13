/// <reference path='./third_party_definitions/_definitions.ts' />
/// <reference path='./util.ts' />

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

escobita.run(['$state', ($state: ng.ui.IStateService) => { // so on landing it goes straight to the lobby
  $state.go("lobby")
  $("#vitu-logo").animate({height:140},2000);
}])


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

// Wraps toastr calls using site "custom look and feel" parameters
namespace Toastr {

  export function success(message: string) {
    return toastr.success(message, '', {
      positionClass: 'toast-bottom-left',
      toastClass: "toast-container",
    });
  }

  export function info(message: string, customOptions?: ToastrOptions) {
    const options = _.extend( {
      positionClass: 'toast-bottom-left',
      toastClass: "toast-container",
    },customOptions)
    return toastr.info(message, '',options);
  }

  export function warn(message: string) {
    return toastr.warning(message, '', {
      positionClass: 'toast-bottom-left',
      toastClass: "toast-container",
    });
  }

  export function error(message: string) {
    return toastr.error(message, '', {
      positionClass: 'toast-bottom-left',
      toastClass: "toast-container",
    });
  }

  export function clear() {
    return toastr.clear();
  }

  export function chat(playerName: string, message: string) {
    return toastr.info(message,`De ${playerName}`, {
      positionClass: 'toast-bottom-left',
      toastClass: "toastr-chat-class",
      titleClass : "toasrt-chat-tittle",
      messageClass: "toasrt-chat-message",
      timeOut: 10000,
    })
  }
}
