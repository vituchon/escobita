/// <reference path='../app.ts' />
/// <reference path='../services/services.ts' />

module Game {
  class Controller {

    public game: Games.Game[];

    constructor(private $state: ng.ui.IStateService, private gamesService: Games.Service, private playersService: Players.Service) {
      console.log("$state.params", $state.params)
      this.game = $state.params["game"]
    }

  }

  escobita.controller('GameController', ['$state', 'GamesService', 'PlayersService', Controller]);
}