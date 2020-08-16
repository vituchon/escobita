/// <reference path='../app.ts' />
/// <reference path='../services/services.ts' />

module Lobby {
  class Controller {

    public games: Games.Game[]

    constructor(private $state: ng.ui.IStateService, private gamesService: Games.Service) {
      this.games = [];
    }

    public createGame(game: Api.Game) {
      this.gamesService.createGame(game).then((createdGame) => {
        this.games.push(createdGame)
      })
    }

    public updateGameList() {
      this.gamesService.getGames().then((games) => {
        this.games = games
      })
    }
  }

  escobita.controller('LobbyController', ['$state', 'GamesService', Controller]);
}