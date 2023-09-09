/// <reference path='../app.ts' />
/// <reference path='../services/_services.d.ts' />

module Lobby {
  interface ViewGamesMode {
      code: string;
      label: string;
  }

  class Controller {

    public games: Games.Game[];
    public players: Players.Player[];

    public player: Players.Player; // the client player
    public playerName: string = ""; // for entering a player name

    public loading: boolean = false;

    public playerGame: Games.Game; // dataholder for a current user's new game
    public canCreateNewGame: boolean;

    public viewGamesMode: ViewGamesMode;
    public viewGamesModes: ViewGamesMode[] = [{code: 'view-all',label: 'Ver todos'},{code: 'select',label: 'Lista Desplegable'}];

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
        this.viewGamesMode = _.size(games) <= 10 ? this.viewGamesModes[0] : this.viewGamesModes[1];
      })
      this.$q.all([getClientPlayerPromise,getGamesPromise]).finally(() => {
        this.loading = false;
      })

      $rootElement.bind("keydown keypress", (event) => {
        if(event.which === 13) {
            $timeout(() => {
              const dialog = document.getElementById('create-game-dialog') as HTMLDialogElement;
              if (dialog?.open) {
                if (!_.isEmpty(this.playerGame?.name)) {
                  this.hideCreateGameDialog();
                  this.createAndResetGame(this.playerGame)
                }
              } else {
                if (!_.isEmpty(this.playerName)) {
                  this.updatePlayerName(this.playerName)
                }
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
        if (Players.isPlayerRegistered(player)) {
          this.showDisplayPlayerName("none");
          this.rearrangeHeaderAfterRegistration("none");
        } else {
          this.rearrangeHeaderBeforeRegistration();
        }
      })
    }

    public isPlayerRegistered(player: Players.Player) {
      return Players.isPlayerRegistered(player)
    }

    private createGame(game: Api.Game) {
      this.loading = true;
      return this.gamesService.createGame(game).then((createdGame) => {
        Toastr.success("Juego creado")
        this.games.push(createdGame)
        return createdGame;
      }).catch((err) => {
        if ( (err?.data as string).includes("has reached the maximum game creation limit")) {
          Toastr.warn("No podés crear más juegos. Borra el que creaste anteriormente y luego crea uno nuevo.")
        }
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
      const playerIsRegistered = Players.isPlayerRegistered(this.player)
      this.player.name = name;
      this.playersService.updatePlayer(this.player).then((player) => {
        const msg = "Nombre de jugador " + ((playerIsRegistered) ? "actualizado" : "registrado")
        Toastr.success(msg)
        this.player = player;
      }).then(() => {
        this.showDisplayPlayerName("transform 1s ease");
        this.rearrangeHeaderAfterRegistration("opacity 1s ease");
      }).finally(() => {
        this.loading = false
      })
    }

    private setCssTransition(transtion: string, ...elements: HTMLElement[]) {
      _.forEach(elements, (element) => {
        element.style.transition = transtion
      })
    }

    public showDisplayPlayerName(transition: string) {
      const enter = document.getElementById('enter-player-name-section');
      const display = document.getElementById('display-player-name-section');

      this.setCssTransition(transition, enter, display)
      enter.style.transform = 'translateX(-101%)';
      display.style.transform = 'translateX(0%)';
    }

    public showEnterPlayerName(transition: string) {
      const enter = document.getElementById('enter-player-name-section');
      const display = document.getElementById('display-player-name-section');

      this.setCssTransition(transition, enter, display)
      enter.style.transform = 'translateX(0)';
      display.style.transform = 'translateX(101%)';
    }

    public rearrangeHeaderAfterRegistration(transition: string) {
      const navPanel = document.getElementById("nav-panel")
      this.setCssTransition(transition, navPanel)
      navPanel.className = 'visible'

      const longHeader = document.getElementById("long-header")
      longHeader.style.display = "none"

      const shortHeader = document.getElementById("short-header")
      shortHeader.style.display = "flex"
    }

    public rearrangeHeaderBeforeRegistration() {
      const longHeader = document.getElementById("long-header")
      longHeader.style.display = "flex"
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
      return !Games.isStarted(game)
    }

    public joinGame(game: Games.Game, player: Players.Player) {
      this.loading = true
      this.gamesService.getGameById(game.id).then((game) => {
        if (Games.isStarted(game)) {
          Toastr.warn("La partida esta en progreso, no te podés unir.")
          return
        }
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

    public deleteGame(game: Games.Game, player: Players.Player) {
      this.loading = true
      this.gamesService.deleteGame(game, player).then(() => {
        Toastr.success("Juego eliminado")
        this.games = this.games.filter((g) => g.id !== game.id)
        return game;
      }).finally(() => {
        this.loading = false
      })
    }

    public canDeleteGame(game: Games.Game, player: Players.Player) {
      return Games.canDeleteGame(game,player)
    }

    public showCreateGameDialog() {
      const dialog = document.getElementById('create-game-dialog') as HTMLDialogElement;
      dialog.showModal()
    }

    public hideCreateGameDialog() {
      const dialog = document.getElementById('create-game-dialog') as HTMLDialogElement;
      dialog.close()
    }
  }

  escobita.controller('LobbyController', ['$rootElement', '$scope', '$timeout','$state', '$q', 'GamesService', 'PlayersService', Controller]);
}