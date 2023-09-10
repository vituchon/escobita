/// <reference path='../app.ts' />
/// <reference path='../services/_services.d.ts' />

module Root {
  class Controller {
    constructor() {

    }

    public $onInit() {
      console.log("A way")
    }
  }

  escobita.controller('RootController', [Controller]);
}