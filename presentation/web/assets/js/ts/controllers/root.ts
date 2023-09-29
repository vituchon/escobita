/// <reference path='../app.ts' />
/// <reference path='../api-model.ts' />

namespace Root {

  export class Controller {
    public currentGame: Api.Game;
    public clientPlayer: Api.Player;

    public isMatchInProgress: boolean;
    public currentStatename: string;

    constructor(private $scope: ng.IScope, private appStateService: AppState.Service) {
    }

    public $onInit() {
      this.$scope.$on('$stateChangeSuccess', (event :ng.IAngularEvent, toState: ng.ui.IState, toParams, fromState, fromParams) => {
        this.currentStatename = toState.name
         // dev notes: is a valid and working approach, but as the lobby screen modifies the player and the game, and the game can be modiied in the game screen IT is not quite correct to handle this values in state transtition
         // so promoting using $watchs on  currentGame and clientPlayer
        /*if (toState.name === "game") {
          const params: {
            game: Api.Game,
            player: Api.Player
          } = toParams;
          this.currentGame = params.game
          this.clientPlayer = params.player
        }*/
      });

      this.$scope.$watch(() => {
        return this.appStateService.get<boolean>("isMatchInProgress")
      }, (isMatchInProgress) => {
        if (!_.isUndefined(isMatchInProgress)) {
          this.isMatchInProgress = isMatchInProgress
        }
      })

      this.$scope.$watch(() => {
        return this.appStateService.get<Api.Game>("currentGame")
      }, (currentGame) => {
        if (!_.isUndefined(currentGame)) {
          this.currentGame = currentGame
        }
      })

      this.$scope.$watch(() => {
        return this.appStateService.get<Api.Player>("clientPlayer")
      }, (clientPlayer) => {
        if (!_.isUndefined(clientPlayer)) {
          this.clientPlayer = clientPlayer
        }
      })
    }

  }

  escobita.controller('rootController', ['$scope', 'AppStateService', Controller]);
}
