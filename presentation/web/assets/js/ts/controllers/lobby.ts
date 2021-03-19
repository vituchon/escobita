/// <reference path='../app.ts' />
/// <reference path='../services/_services.d.ts' />

module Lobby {
  class Controller {

    public games: Games.Game[];
    public players: Players.Player[];

    public player: Players.Player; // the client player
    public playerName: string = ""; // for entering a player name

    public loading: boolean = false;
    public showCards: boolean = false;

    constructor(private $rootElement: ng.IRootElementService,private $scope: ng.IScope, $timeout: ng.ITimeoutService,
        private $state: ng.ui.IStateService, private $q: ng.IQService, private gamesService: Games.Service,
        private playersService: Players.Service) {
      this.games = [];
      this.loading = true
      const getClientPlayerPromise = this.playersService.getClientPlayer().then((player) => {
        this.player = player;
        this.playerName = this.player.name
      })
      const getGamesPromise = this.gamesService.getGames().then((games) => {
        this.games = games
      })
      this.$q.all([getClientPlayerPromise,getGamesPromise]).finally(() => {
        this.loading = false;
      })

      $rootElement.bind("keydown keypress", (event) => {
        if(event.which === 13) {
            $timeout(() => {
              if (!_.isUndefined(this.playerName)) {
                this.updatePlayerName(this.playerName)
              }
            });
            event.preventDefault();
        }
      });
      $scope.$on('$destroy', function() {
        $rootElement.unbind("keydown keypress")
      });

    }

    public createGame(game: Api.Game) {
      this.loading = true;
      this.gamesService.createGame(game).then((createdGame) => {
        this.games.push(createdGame)
      }).finally(() => {
        this.loading = false;
      })
    }

    private canCreateGame(game: Api.Game) {
      return !this.loading && !_.isEmpty(game) && !_.isEmpty(game.name)
    }

    public updateGameList() {
      this.loading = true;
      this.gamesService.getGames().then((games) => {
        this.games = games
      }).finally(() => {
        this.loading = false;
      })
    }

    public updatePlayerName(name: string) {
      this.loading = true
      this.player.name = name;
      this.playersService.updatePlayer(this.player).then((player) => {
        this.player = player;
      }).then(() => {
        this.showCards = true;
      }).finally(() => {
        this.loading = false
      })
    }

    public updatePlayersList() {
      this.playersService.getPlayers().then((players) => {
        this.players = players
      })
    }

    public doesGameAcceptPlayers(game: Games.Game) {
      return !Games.hasMatchInProgress(game)
    }

    public joinGame(game: Games.Game, player: Players.Player) {
      this.loading = true
      this.gamesService.getGameById(game.id).then((game) => {
        Games.addPlayer(game, player)
        this.gamesService.updateGame(game).then(() => {
          this.$state.go("game", {
            game: game,
            player: player,
          }, {relative: false})
        })
      }).finally(() => {
        this.loading = false
      })
    }
  }

  escobita.controller('LobbyController', ['$rootElement', '$scope', '$timeout','$state', '$q', 'GamesService', 'PlayersService', Controller]);
}