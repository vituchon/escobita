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

    public playerGame: Games.Game; // dataholder for a current user's new game
    public canCreateNewGame: boolean;

    public viewGamesMode: string;

    constructor($rootElement: ng.IRootElementService, $scope: ng.IScope, $timeout: ng.ITimeoutService,
        private $state: ng.ui.IStateService, private $q: ng.IQService, private gamesService: Games.Service,
        private playersService: Players.Service) {
      this.games = [];
      this.loading = true
      const getClientPlayerPromise = this.playersService.getClientPlayer().then((player) => {
        this.player = player;
        this.playerName = this.player.name
        return player
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

      $scope.$watch(() => {
        return this.canCreateGame(this.playerGame)
      }, (can) => {
        this.canCreateNewGame = !!can;
      })

      getClientPlayerPromise.then((player) => {
        if (!_.isEmpty(player.name)) {
          this.showCards = true;
          this.showPlayerNameStatic(false);
        }
      })
    }

    private createGame(game: Api.Game) {
      this.loading = true;
      return this.gamesService.createGame(game).then((createdGame) => {
        this.games.push(createdGame)
        return createdGame;
      }).finally(() => {
        this.loading = false;
      })
    }

    public createAndResetGame(game: Api.Game) {
      this.createGame(game).then(() => {
        game.name = '';
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
        this.showPlayerNameStatic(true);
      }).finally(() => {
        this.loading = false
      })
    }

    public showPlayerNameStatic(animate: boolean) {
      const enter = document.getElementById('enter-player-name-section');
      const display = document.getElementById('display-player-name-section');
      if (animate) {
        enter.style.transition = "transform 1s ease";
        display.style.transition = "transform 1s ease";
      } else {
        enter.style.transition = "none";
        display.style.transition = "none";
      }

      enter.style.transform = 'translateX(-101%)';
      display.style.transform = 'translateX(0%)';
    }

    public showPlayerNameForm(animate: boolean) {
      const enter = document.getElementById('enter-player-name-section');
      const display = document.getElementById('display-player-name-section');
      if (animate) {
        enter.style.transition = "transform 1s ease";
        display.style.transition = "transform 1s ease";
      } else {
        enter.style.transition = "none";
        display.style.transition = "none";
      }

      enter.style.transform = 'translateX(0)';
      display.style.transform = 'translateX(101%)';
    }

    public canUpdatePlayerName(name: string) {
      return !this.loading && !_.isEmpty(name);
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
          })
        })
      }).finally(() => {
        this.loading = false
      })
    }
  }

  escobita.controller('LobbyController', ['$rootElement', '$scope', '$timeout','$state', '$q', 'GamesService', 'PlayersService', Controller]);
}