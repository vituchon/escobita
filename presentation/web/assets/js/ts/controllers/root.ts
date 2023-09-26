/// <reference path='../app.ts' />
/// <reference path='../api-model.ts' />


namespace Root {

  export class Controller {
    public currentStatename: string;
    public currentGame: Api.Game;
    public player: Api.Player;
    public hasAMatchInProgress: boolean = false;

    constructor(private $scope: ng.IScope, private appStateService : AppState.Service) {
    }

    public $onInit() {
      this.$scope.$on('$stateChangeSuccess', (event :ng.IAngularEvent, toState: ng.ui.IState, toParams, fromState, fromParams) => {
        this.currentStatename = toState.name
        if (toState.name === "game") {
          const params= <{game: Api.Game, player: Api.Player}> toParams
          this.currentGame = params.game
        }
      });

      this.$scope.$watch(() => {
        return this.currentGame
      }, (newGame) => {
        if (!_.isUndefined(newGame)) {
          this.hasAMatchInProgress = Games.hasMatchInProgress(newGame)
        }
      })

      this.$scope.$watch(() => {
        return this.appStateService.get<boolean>("isAMatchInProgress")
      }, (isAMatchInProgres) => {
        if (!_.isUndefined(isAMatchInProgres)) {
          this.hasAMatchInProgress = isAMatchInProgres
        }
      })

      this.$scope.$watch(() => {
        return this.appStateService.get<Api.Player>("player")
      }, (player) => {
        if (!_.isUndefined(player)) {
          this.player = player;
        }
      })
    }
  }

  escobita.controller('rootController', ['$scope', 'AppStateService', Controller]);
}
