/// <reference path='../app.ts' />
/// <reference path='../services/services.ts' />

module Lobby {
  class Controller {

    public games: Games.Game[];
    public players: Players.Player[];

    public player: Players.Player;

    constructor(private $state: ng.ui.IStateService, private gamesService: Games.Service, private playersService: Players.Service) {
      this.games = [];
      this.playersService.getClientPlayer().then((player) => {
        this.player = player;
      })
      this.gamesService.getGames().then((games) => {
        this.games = games
      })
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

    public updatePlayer(player: Players.Player) {
      this.playersService.updatePlayer(player).then((player) => {
        this.player = player;
      })
    }

    public updatePlayersList() {
      this.playersService.getPlayers().then((players) => {
        this.players = players
      })
    }

    public doesGameAcceptPlayers(game: Games.Game) {
      return !Games.isStarted(game)
    }

    public joinGame(player: Players.Player, game: Games.Game) {
      console.log(player,game)
      // this.$state.go => game screen with game as data!
    }
  }

  escobita.controller('LobbyController', ['$state', 'GamesService', 'PlayersService', Controller]);
}