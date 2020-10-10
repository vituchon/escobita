/// <reference path='../app.ts' />
/// <reference path='../services/_services.d.ts' />

module About {
  class Controller {
    constructor(private $state: ng.ui.IStateService) {
    }
  }

  escobita.controller('AboutController', ['$state', Controller]);
}