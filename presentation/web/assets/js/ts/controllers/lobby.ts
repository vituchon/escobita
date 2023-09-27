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

    public isPlayerRegistered: boolean = false;
    public player: Players.Player; // the client player

    public loading: boolean = false;

    public playerGame: Games.Game; // dataholder for a current user's new game
    public canCreateNewGame: boolean;

    public viewGamesMode: ViewGamesMode;
    public viewGamesModes: ViewGamesMode[] = [{code: 'view-all',label: 'Ver todos'},{code: 'select',label: 'Lista Desplegable'}];

    constructor(private $rootElement: ng.IRootElementService,private $scope: ng.IScope,private $timeout: ng.ITimeoutService,
        private $state: ng.ui.IStateService, private $q: ng.IQService, private gamesService: Games.Service,
        private playersService: Players.Service, private $window: ng.IWindowService, private appStateService: AppState.Service) {
      this.games = [];
      this.loading = true
      this.init().then(() => {
        try { // if something unexpected happens angularjs wraps the error within the promise and shallow its
          this.setupWatchs();
          this.setupUI();
        } catch(err) {
          console.error(err) // ... so this way at least I log something here
          Toastr.error("Hubo un error con el navegador 😿")
          throw err
        }
      }).catch((err) => {
        console.error(err)
        Toastr.error("Hubo un error con el navegador 😿")
      }).finally(() => {
        this.loading = false
      })
    }

    private setupWatchs() {
      this.$scope.$watch(() => {
        return this.canCreateGame(this.playerGame)
      }, (can) => {
        this.canCreateNewGame = !!can;
      })
    }

    private setupUI() {
      this.$rootElement.bind("keydown", (event) => {
        if(event.which === 13) {
            this.$timeout(() => {
              const dialog = document.getElementById('create-game-dialog') as HTMLDialogElement;
              if (dialog?.open) {
                if (!_.isEmpty(this.playerGame?.name)) {
                  this.hideCreateGameDialog();
                  this.createAndResetGame(this.playerGame)
                }
              } else {
                if (!_.isEmpty(this.player.name)) {
                  this.updatePlayerName(this.player.name)
                }
              }
            });
            event.preventDefault();
        }
      });
      this.$scope.$on('$destroy', () => {
        this.$rootElement.unbind("keydown")
      });

      if (this.isPlayerRegistered) {
        this.showDisplayPlayerName("none");
        this.rearrangeHeaderAfterRegistration("none");
      } else {
        this.rearrangeHeaderBeforeRegistration();
        const nameInput = document.getElementById("name-input")
        nameInput.focus()
      }
    }

    private init() {
      const getClientPlayerPromise = this.playersService.getClientPlayer().then((player) => {
        this.player = player;
        this.appStateService.set("clientPlayer", player)
        this.isPlayerRegistered = AppState.isPlayerRegistered(player)
        return player
      })
      const getGamesPromise = this.gamesService.getGames().then((games) => {
        this.games = games
        this.viewGamesMode = _.size(games) <= 10 ? this.viewGamesModes[0] : this.viewGamesModes[1];
      })

      return this.$q.all([getClientPlayerPromise,getGamesPromise]).finally(() => {

      })
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
      this.player.name = name;
      const wasPlayerRegistered = this.isPlayerRegistered
      this.playersService.updatePlayer(this.player).then((player) => {
        this.player = player;
        this.isPlayerRegistered = true;
        this.appStateService.set("clientPlayer", player)
        const msg = "Nombre de jugador " + ((wasPlayerRegistered) ? "actualizado" : "registrado")
        Toastr.success(msg)
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

    public animateCardOnClick(id: number) {
      const card = document.getElementById(id.toString());
      var x = 200 + Math.random() * 200
      var y = 200 + Math.random() * 200
      if (Math.random() < 0.5) { // thanks https://stackoverflow.com/a/36756480/903998
        x = -x
      }
      if (Math.random() < 0.5) { // thanks again :D
        y = -y
      }
      card.style.transform = `translate(${x}%,${y}%)`
    }
  }

  escobita.controller('LobbyController', ['$rootElement', '$scope', '$timeout','$state', '$q', 'GamesService', 'PlayersService', '$window', 'AppStateService', Controller]);
}