/// <reference path='../app.ts' />
/// <reference path='../api-model.ts' />


namespace Root {

  export class Controller {
    //private currentGame: Api.Game;
    public currentStatename: string;
    constructor(private $scope: ng.IScope) {
    }

    public $onInit() {
      this.$scope.$on('$stateChangeSuccess', (event :ng.IAngularEvent, toState: ng.ui.IState, toParams, fromState, fromParams) => {
        this.currentStatename = toState.name
      });
    }

  }

  escobita.controller('rootController', ['$scope', Controller]);
}
