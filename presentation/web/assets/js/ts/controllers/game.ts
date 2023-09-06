/// <reference path='../app.ts' />
/// <reference path='../services/_services.d.ts' />
/// <reference path='../directives/_directives.d.ts' />

module Game {


  namespace UIMessages {

    export const minFontSize = 10;
    const maxFontSize = 34;
    const fontSizeAmplitude = maxFontSize - minFontSize;


    function determineFontSize(position: number) {
        const size = minFontSize + (fontSizeAmplitude) / position; // B + M/X (Homográfica desplazada)
        return size;
    }

    export interface FontSizeByPlayerName extends _.Dictionary<number> {
      [name:string]: number;
    }

    export function calculateFontSizeByPlayerName(positionsByPlayerName: Matchs.Rules.PositionByPlayerName) : FontSizeByPlayerName {
      return _.reduce(positionsByPlayerName,(acc, position, name) => {
        acc[name] = determineFontSize(position + 1) // position starts at 0
        return acc
      },<FontSizeByPlayerName>{})
    }
  }

  // END :  move to util.js
  class Controller {

    public game: Games.Game; // the current game
    public player: Players.Player; // the client player
    public isPlayerTurn: boolean = undefined; // initial value because match didn't start, on start a true/false value is assigned
    public isPlayerGameOwner: boolean;
    public currentTurnPlayer: Players.Player; // the player that acts in the current turn
    public messages: Messages.Message[] = []; // all from the server related to this game
    public isBoardCardSelectedById: _.Dictionary<boolean>;
    public selectedHandCard: Api.Card;

    public playerMessage: Api.Message; // current player message
    private sendingMessage: boolean = false;
    private allowSendMessage: boolean = true; // avoid message spawn
    public isChatEnabled: boolean = false;
    private lastChatUpdateUnixTimestamp: number = 0;
    private currentFontSizeByPlayerName: UIMessages.FontSizeByPlayerName; // funny font size to use by player name
    private currentPositionByPlayerName: Matchs.Rules.PositionByPlayerName; // positions by player name

    public isMatchInProgress: boolean = false;
    public currentMatchStats: Api.ScoreSummaryByPlayerName;

    public players: Players.Player[]; // not sure if it will be use somewhere!
    public playersById: Util.EntityById<Players.Player>;

    private lastUpdateUnixTimestamp: number = undefined;

    public formatUnixTimestamp = Util.unixToReadableClock
    public translateSuit = Cards.Suits.translate

    public loading: boolean = false;
    public displayCardsAsSprites: boolean = true;

    private updaterInterval: ng.IPromise<any>; // "handler" to the one interval that updates the UI according to the controller's state

    constructor(private $rootElement: ng.IRootElementService, private $rootScope: ng.IRootScopeService, private $scope: ng.IScope, private $state: ng.ui.IStateService,
      private gamesService: Games.Service, private playersService: Players.Service,  private messagesService: Messages.Service, private webSocketsService: WebSockets.Service,
      private $interval: ng.IIntervalService, private $timeout: ng.ITimeoutService, private $q: ng.IQService, private $window: ng.IWindowService) {
      this.game = $state.params["game"]
      this.player = $state.params["player"]
      this.isPlayerGameOwner = Games.isPlayerOwner(this.player,this.game)
      this.playerMessage = Messages.newMessage(this.game.id,this.player.id,"");

      this.$scope.$watch(() => {
        return this.isMatchInProgress
      },(isMatchInProgress,wasMatchInProgress) => {
        if (isMatchInProgress && !wasMatchInProgress) {
          Toastr.info("¡La partida ha comenzado!")
        }
        if (!isMatchInProgress && wasMatchInProgress) {
          Toastr.success("¡La partida ha terminado!")
          //this.displayCurrentMatchStatsCopy = true
        }
      })

      this.setupPullRefresh(2000) // experimenting see // (*)
      this.webSocketsService.retrieve().then((ws) => {
        this.gamesService.bindWebSocket(this.game.id).then(() => {
          this.setupPushRefresh(ws)
        })
      }).catch((reason) => {
        console.warn("could not adquire web socket: ", reason);
        Toastr.error(`No se pudo establecer conexión con el servidor, motivo: ${reason}`)
      })

      this.$scope.$watch(() => {
        return this.displayCardsAsSprites;
      }, (displayCardsAsSprites) => {
        if (_.isUndefined(displayCardsAsSprites)) {
          return
        }
        const displayMode = displayCardsAsSprites ? 'sprite' : 'text'
        this.$rootScope.$broadcast(Cards.changeDisplayModeEventName, displayMode);
      })

      $rootElement.bind("keydown keypress", (event) => {
        if(event.which === 13) {
            $("#chat-press-enter-hint").hide();
            $timeout(() => {
              if (this.isChatEnabled && this.canSendMessage(this.playerMessage)) {
                this.sendAndCleanMessage(this.playerMessage);
              }
            });
            event.preventDefault();
        }
      });
      $scope.$on('$destroy', function() {
        $rootElement.unbind("keydown keypress")
      });

      this.$scope.$watch(() => {
        if (!this.isMatchInProgress) {
          return undefined
        }
        return this.game.currentMatch.currentRound.currentTurnPlayer.name;
      }, (currentTurnPlayerName, previousTurnPlayerName) => {
        if (_.isUndefined(previousTurnPlayerName)) {
          return
        }
        if (previousTurnPlayerName !== this.player.name) {
          this.displayLastAction();
        }
      })

      // event binding on dynamically created, "live" watching  https://stackoverflow.com/a/1207393/903998
      $("div.game-match-section").on("mouseover mouseout","div.play-section .card-image",(event: Event) => {
        if (event.type === "mouseover") { // taken inspiration from https://stackoverflow.com/a/13504775/903998
          var rotateDegress = Math.random() * 10 - Math.random() * 5;
          $((<any>(event.target)).parentElement.parentElement).css('transform', 'rotate(' + rotateDegress + 'deg) scale(1.25)');
        } else {
          $((<any>(event.target)).parentElement.parentElement).css('transform', 'none');
        }
      })
    }

    private setupPushRefresh(webSocket: WebSocket) {
      webSocket.onmessage = (event) => {
        const notification : {
          kind: string,
          data: {
            game: Api.Game,
            action?: Api.PlayerAction
          };
        } = JSON.parse(event.data)
        console.log("llega una notificación", notification);

        switch (notification.kind) {
          case "drop":
          case "take":
          case "updated":
          case "resume":
            this.setGame(notification.data.game)
            break;
          default:
            break;
        }
      }
      const onUnload = (event: any):any => {
        //this.gamesService.unbindWebSocket(this.game.id)
        event.preventDefault();
        return event.returnValue = null;
      }
      this.$window.addEventListener("beforeunload",onUnload)
      this.$scope.$on('$destroy', () => {
        this.gamesService.unbindWebSocket(this.game.id)
        this.$window.removeEventListener("beforeunload",onUnload)
      })

      this.$scope.$on('$stateChangeStart', (event, toState, toParams, fromState, fromParams) => {
        if (this.isMatchInProgress) {
          alert("Si te vas en plena partida, vas a cagarle la partida a los demás")
          event.preventDefault(); // Prevent the state change for now
        }
      });
    }

    private setupPullRefresh(delay: number = 2000) {
      this.$scope.$on('$destroy', () => {
        this.$interval.cancel(this.updaterInterval)
      });

      this.updaterInterval = this.$interval(() => {
        const mustRefreshGame = (!this.isMatchInProgress) || (this.isMatchInProgress && !this.isPlayerTurn) // ~p v (p y q~) == ~p v ~q
        var refreshGamePromise = this.$q.when(this.game)
        /*if (mustRefreshGame) { // (*) intenttionally commented so push notifications update the game but pull notification the chat
          refreshGamePromise = this.refreshGame()
        }*/
        var updateChatPromise = this.$q.when(this.messages)
        if (this.isChatEnabled) {
          updateChatPromise = this.updateChat();
        }

        return this.$q.all([refreshGamePromise,updateChatPromise]).then((response) => {
          if (!_.isUndefined(this.lastUpdateUnixTimestamp)) {
            const nowUnixTimestamp = moment().unix()
            const seconds = nowUnixTimestamp - this.lastUpdateUnixTimestamp
            //console.log("demora aproximada ",seconds)
            if (seconds > 5) {
              Toastr.warn("Puede que haya algunos problemas de conexión!")
            }
          }
          this.lastUpdateUnixTimestamp = 	moment().unix()
          return response
        })
      }, delay)
    }

    private displayLastAction() {
      const options: ToastrOptions = {
        timeOut: 5000,
        toastClass: "toastr-info-action-class",
        closeButton: true,
      }
      Toastr.info(`${this.generateLastActionDescription()}`, options)
    }

    private generateLastActionDescription() {
      const lastActionIndex = _.size(this.game.currentMatch.playerActions) - 1
      const lastAction = this.game.currentMatch.playerActions[lastActionIndex];
      return this.generateActionDescription(lastAction);
    }

    private generateActionDescription(action: Api.PlayerAction) {
      if (Matchs.isTakeAction(action)) {
        const cardsText = _.map(action.boardCards,(card) => this.cardToText(card))
        var description = `levantó ${cardsText.join(",")} usando ${this.cardToText(action.handCard)}`
        if (action.isEscobita) {
          description += "<b>e hizó escobita!</b>"
        }
        return `<b>${action.player.name}</b>: ${description}`
      } else {
        return `<b>${action.player.name}</b>: descartó ${this.cardToText(action.handCard)}`
      }
    }

    private cardToText(card: Api.Card) {
      return `${card.rank} de ${Cards.Suits.translate(card.suit)}`
    }

    public updateChat() {
      return this.playersService.getPlayers().then((players) => {
        this.players = players
        this.playersById = Util.toMapById(this.players);
        return this.players
      }).then((players) => {
        return this.messagesService.getMessagesByGame(this.game.id, this.lastChatUpdateUnixTimestamp).then((incomingMessages) => {
          console.log("incoming are: ", incomingMessages)
          this.lastChatUpdateUnixTimestamp = moment().unix();
          _.forEach(incomingMessages,(incomingMessage) => {
            const player = this.playersById[incomingMessage.playerId]
            const $elem = Toastr.chat(player.name,incomingMessage.text)
            const fontSize = this.getFontSize(player)
            $(".toasrt-chat-message",$elem).css("font-size", fontSize + "px");
          })
          this.messages.push(...incomingMessages)
          return undefined; // it is the default return value, see https://plnkr.co/edit/ZdXQymjYFON0VIcD
        })
      })
    }

    public sendAndCleanMessage(msg: Api.Message) {
      this.sendingMessage = true;
      this.messagesService.sendMessage(msg).then(() => {
        this.cleanMessage(msg);
      }).finally(() => {
        this.sendingMessage = false;
      })
      this.allowSendMessage = false;
      this.$timeout(() => {
        this.allowSendMessage = true;
      }, 2000)
    }

    public canSendMessage(msg: Api.Message) {
      return !this.loading && this.allowSendMessage && !this.sendingMessage && !_.isEmpty(msg.text);
    }

    private cleanMessage(msg: Api.Message) {
      msg.text = '';
    }

    public startGame(game: Games.Game, players: Players.Player[]) {
      this.loading = true;
      return this.gamesService.getGameById(this.game.id).then((game) => { // update the game (in case another player join on other client AND this client is outdated)
        return this.gamesService.startGame(game).then((game) => {
          return this.setGame(game)
        })
      }).finally(() => {
        this.loading = false
      })
    }

    public hasValidTakeAction() {
      if (_.isEmpty(this.selectedHandCard)) {
        return false
      }
      const selectedBoardCards = this.getSelectedBoardCards()
      if (_.isEmpty(selectedBoardCards)) {
        return false
      }
      return Matchs.Rules.canTakeCards(this.selectedHandCard,selectedBoardCards)
    }

    public performTakeAction() {
      const selectedBoardCards = this.getSelectedBoardCards()
      const takeAction = Matchs.createTakeAction(this.player,selectedBoardCards,this.selectedHandCard)
      this.loading = true;
      this.gamesService.performTakeAction(this.game,takeAction).then((data) => {
        if (data.action.isEscobita) {
          Toastr.success("Has hecho escoba!")
        }
        //return this.setGame(data.game) // don't need to refresh as this clients gets notified via ws
      }).finally(() => {
        this.isBoardCardSelectedById = {}
        this.loading = false;
      })
    }

    private getSelectedBoardCards() {
       return _.reduce(this.isBoardCardSelectedById,(acc,selected,id) => {
        if (selected) {
          const card = _.find(this.game.currentMatch.matchCards.board,(boardCard) => boardCard.id === +id)
          const isNotInBoard = _.isUndefined(card)
          if (isNotInBoard) {
            console.warn("suspicious things, programmer must check something...")
          } else {
            acc.push(card)
          }
        }
        return acc
      },<Api.Card[]>[])
    }

    public hasValidDropAction() {
      return !_.isEmpty(this.selectedHandCard)
    }

    public performDropAction() {
      const selectedBoardCards = this.getSelectedBoardCards()
      const dropAction = Matchs.createDropAction(this.player,this.selectedHandCard)
      this.loading = true;
      this.gamesService.performDropAction(this.game,dropAction).then((data) => {
        //return this.setGame(data.game) // don't need to refresh as this clients gets notified via ws
      }).finally(() => {
        this.selectedHandCard = undefined
        this.loading = false;
      })
    }

    private setGame(game: Games.Game) {
      this.game = game;
      this.isMatchInProgress = Games.hasMatchInProgress(game)
      const currentMatchIndex = _.size(this.game.matchs)
      if (this.isMatchInProgress) {
        this.currentTurnPlayer = this.game.currentMatch.currentRound.currentTurnPlayer;
        this.isPlayerTurn = Rounds.isPlayerTurn(this.game.currentMatch.currentRound,this.player)
        return this.updateGameStats(currentMatchIndex).then(() => {
          return this.game
        })
      } else {
        const hasPreviousMatchs =  (currentMatchIndex > 0)
        if (hasPreviousMatchs) {
          return this.updateGameStats(currentMatchIndex-1).then(() => {
            return this.game
          })
        } else {
          return this.$q.when(this.game)
        }
      }
    }

    public refreshGame() {
      return this.gamesService.getGameById(this.game.id).then((game) => {
        return this.setGame(game)
      })
    }

    private updateGameStats(matchIndex: number) {
      return this.gamesService.calculateStatsByGameId(this.game.id,matchIndex).then((stats) => {
        this.currentMatchStats = stats;
        this.currentPositionByPlayerName = Matchs.Rules.calculatePositionByPlayerName(stats)
        this.currentFontSizeByPlayerName = UIMessages.calculateFontSizeByPlayerName(this.currentPositionByPlayerName)
        return stats
      })
    }

    public getFontSize(player: Players.Player) {
      if (_.isEmpty(this.currentFontSizeByPlayerName)) {
        return UIMessages.minFontSize;
      } else {
        return this.currentFontSizeByPlayerName[player.name]
      }
    }

  }

  escobita.controller('GameController', ['$rootElement','$rootScope','$scope','$state', 'GamesService', 'PlayersService',
    'MessagesService', 'WebSocketsService', '$interval', '$timeout', '$q','$window', Controller]);
}
