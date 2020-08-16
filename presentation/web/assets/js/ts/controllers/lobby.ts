/// <reference path='../app.ts' />

module Lobby {
  class Controller {
    constructor(private $state: ng.ui.IStateService) {
    }
  }

  escobita.controller('LobbyController', ['$state', Controller]);
}