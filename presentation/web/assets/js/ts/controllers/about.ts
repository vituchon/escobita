/// <reference path='../app.ts' />

module About {
  class Controller {
    constructor(private $state: ng.ui.IStateService) {
    }
  }

  escobita.controller('AboutController', ['$state', Controller]);
}