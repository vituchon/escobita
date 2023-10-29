/// <reference path='../app.ts' />
/// <reference path='../services/_services.d.ts' />
/// <reference path='../directives/_directives.d.ts' />

namespace Game {

  namespace UIMessages {

    export const minFontSize = 10;
    const maxFontSize = 34;
    const fontSizeAmplitude = maxFontSize - minFontSize;


    function determineFontSize(position: number) {
        const size = minFontSize + (fontSizeAmplitude) / position; // B + M/X (Homográfica desplazada)
        return size;
    }

    export interface FontSizeByPlayerUniqueKey extends _.Dictionary<number> {
      [uniqueKey:string]: number;
    }

    export function calculateFontSizeByPlayerUniqueKey(positionsByPlayerUniqueKey: Matchs.Rules.PositionByPlayerUniqueKey) : FontSizeByPlayerUniqueKey {
      return _.reduce(positionsByPlayerUniqueKey,(acc, position, name) => {
        acc[name] = determineFontSize(position + 1) // position starts at 0
        return acc
      },<FontSizeByPlayerUniqueKey>{})
    }
  }

  const game: Games.Game = {
    "players": [
      {
        "name": "Betoven",
        "id" : 1,
      }
    ],
    "id": 1,
    "name": "1",
    "owner": {
      id: 1,
      name: "Betoven"
    }
  }

  const onGoingGame: Games.Game = <any>{
    "players": [
      {
        "name": "Betoven",
        "id": 1,
      }
    ],
    "currentMatch": {
      "players": [
        {
          "name": "Betoven",
          "id": 1,
        }
      ],
      "actionsByPlayerUniqueKey": {
        "1|Betoven": []
      },
      "playerActions": [],
      "matchCards": {
        "board": [
          {
            "id": 27,
            "suit": 2,
            "rank": 7
          },
          {
            "id": 7,
            "suit": 0,
            "rank": 7
          },
          {
            "id": 39,
            "suit": 3,
            "rank": 11
          },
          {
            "id": 30,
            "suit": 2,
            "rank": 12
          }
        ],
        "left": [
          {
            "id": 23,
            "suit": 2,
            "rank": 3
          },
          {
            "id": 17,
            "suit": 1,
            "rank": 7
          },
          {
            "id": 40,
            "suit": 3,
            "rank": 12
          },
          {
            "id": 8,
            "suit": 0,
            "rank": 10
          },
          {
            "id": 33,
            "suit": 3,
            "rank": 3
          },
          {
            "id": 20,
            "suit": 1,
            "rank": 12
          },
          {
            "id": 31,
            "suit": 3,
            "rank": 1
          },
          {
            "id": 38,
            "suit": 3,
            "rank": 10
          },
          {
            "id": 13,
            "suit": 1,
            "rank": 3
          },
          {
            "id": 15,
            "suit": 1,
            "rank": 5
          },
          {
            "id": 18,
            "suit": 1,
            "rank": 10
          },
          {
            "id": 12,
            "suit": 1,
            "rank": 2
          },
          {
            "id": 37,
            "suit": 3,
            "rank": 7
          },
          {
            "id": 2,
            "suit": 0,
            "rank": 2
          },
          {
            "id": 5,
            "suit": 0,
            "rank": 5
          },
          {
            "id": 10,
            "suit": 0,
            "rank": 12
          },
          {
            "id": 26,
            "suit": 2,
            "rank": 6
          },
          {
            "id": 25,
            "suit": 2,
            "rank": 5
          },
          {
            "id": 19,
            "suit": 1,
            "rank": 11
          },
          {
            "id": 32,
            "suit": 3,
            "rank": 2
          },
          {
            "id": 1,
            "suit": 0,
            "rank": 1
          },
          {
            "id": 9,
            "suit": 0,
            "rank": 11
          },
          {
            "id": 11,
            "suit": 1,
            "rank": 1
          },
          {
            "id": 28,
            "suit": 2,
            "rank": 10
          },
          {
            "id": 3,
            "suit": 0,
            "rank": 3
          },
          {
            "id": 6,
            "suit": 0,
            "rank": 6
          },
          {
            "id": 21,
            "suit": 2,
            "rank": 1
          },
          {
            "id": 22,
            "suit": 2,
            "rank": 2
          },
          {
            "id": 4,
            "suit": 0,
            "rank": 4
          },
          {
            "id": 36,
            "suit": 3,
            "rank": 6
          },
          {
            "id": 34,
            "suit": 3,
            "rank": 4
          },
          {
            "id": 14,
            "suit": 1,
            "rank": 4
          },
          {
            "id": 35,
            "suit": 3,
            "rank": 5
          }
        ],
        "byPlayer": {
          "Betoven": {
            "taken": null,
            "hand": [
              {
                "id": 16,
                "suit": 1,
                "rank": 6
              },
              {
                "id": 29,
                "suit": 2,
                "rank": 11
              },
              {
                "id": 24,
                "suit": 2,
                "rank": 4
              }
            ]
          }
        }
      },
      "firstPlayerIndex": 0,
      "roundNumber": 1,
      "currentRound": {
        "currentTurnPlayer": {
          "name": "Betoven"
        },
        "consumedTurns": 1,
        "number": 1
      }
    },
    "id": 1,
    "owner": {
      "name": "Betoven",
      "id": 1
    }
  }

  const player: Players.Player = {
    name: "Betoven",
    id: 1
  }

  // END :  move to util.js
  class Controller {

    public game: Games.Game; // the current game
    public player: Players.Player; // the client player
    public isPlayerTurn: boolean = undefined; // initial value because match didn't start, on start a true/false value is assigned
    public isClientPlayerGameOwner: boolean;
    public currentTurnPlayer: Players.Player; // the player that acts in the current turn
    public messages: Messages.Message[] = []; // persistent messages of this game (retrieved from the server using the previous message Api that works with persistent messages)
    public isBoardCardSelectedById: _.Dictionary<boolean>;
    public selectedHandCard: Api.Card;

    public playerMessage: Games.VolatileMessage; // last player message
    private sendingMessage: boolean = false;
    private allowSendMessage: boolean = true; // avoid message spawn
    public isChatEnabled: boolean = true;
    private currentFontSizeByPlayerUniqueKey: UIMessages.FontSizeByPlayerUniqueKey; // funny font size to use by player name
    private currentPositionByPlayerUniqueKey: Matchs.Rules.PositionByPlayerUniqueKey; // positions by player name

    public isMatchInProgress: boolean = false;
    public currentMatchStatsByPlayerUniqueKey: Api.ScoreSummaryByPlayerUniqueKey;

    public formatUnixTimestamp = Util.unixToReadableClock
    public translateSuit = Cards.Suits.translate

    public loading: boolean = false;
    public displayCardsAsSprites: boolean = true;

    constructor(private $rootElement: ng.IRootElementService, private $rootScope: ng.IRootScopeService, private $scope: ng.IScope, $state: ng.ui.IStateService,
        private gamesService: Games.Service, private webSocketsService: WebSockets.Service, private $timeout: ng.ITimeoutService,
        private $q: ng.IQService, private $window: ng.IWindowService, private appStateService: AppState.Service) {
      this.game = $state.params["game"]// || onGoingGame
      this.player = $state.params["player"]// || player
      this.setGame(this.game)

        /*const navPanel = document.getElementById("nav-panel")
        navPanel.className = 'visible'

        const shortHeader = document.getElementById("short-header")
        shortHeader.style.display = "flex"*/

      this.isClientPlayerGameOwner = Games.isPlayerOwner(this.game, this.player)
      this.playerMessage = Games.newMessage(this.game.id,this.player,""); // dev notes: the gameId and playerId are constants but the text (last arg) is set from the UI using ng-model="ctr.playerMessage.text"

      this.loading = true;
      this.init().then((ws) => {
        try { // if something unexpected happens angularjs wraps the error within the promise and shallow its
          this.setupWatchs();
          this.setupUI();
          this.setupWebsockt(ws)
        } catch(err) {
          console.error(err) // ... so this way at least I log something here
          Toastr.error("Hubo un error con el navegador 😿")
          throw err
        }
        this.appStateService.set("currentGame", this.game)
      }).catch((err) => {
        console.error(err)
        this.appStateService.set("isMatchInProgress", false)
        this.appStateService.set("currentGame", null)
        $state.go("lobby")
      }).finally(() => {
        this.loading = false
      })
    }

    private init() {
      return this.webSocketsService.retrieve().then((ws) => {
        return this.gamesService.bindWebSocket(this.game.id).then(() => {
          return ws
        })
      }).catch((retriveAndBindErr) => {
        console.warn("could not adquire web socket: ", retriveAndBindErr);
        Toastr.error(`No se pudo establecer conexión con el servidor 😢`)
        return this.webSocketsService.release().then(() => {
          Toastr.info(`Se reseteo la conexión, probá ingresar nuevamente`)
          throw retriveAndBindErr
        }).catch((releaseErr) => {
          Toastr.warn(`Asegurate de tener solo una pestaña en ${window.location.origin} y probá recargar la página`)
          if (releaseErr != retriveAndBindErr) { // if they are equals for avoiding loggin twice err (as throw err above is actually catched by ths catch block)
            console.warn("could not release web socket: ", releaseErr);
          }
          throw retriveAndBindErr
        })
      })
    }

    private setupWatchs() {
      this.$scope.$watch(() => {
        return this.isMatchInProgress
      },(isMatchInProgress,wasMatchInProgress) => {
        if (isMatchInProgress && !wasMatchInProgress) {
          this.suggestionRequestCount = 0; // reset "take action" suggestions request counter
          Toastr.info("¡La partida ha comenzado! ¡Las puertas se han cerrado!")
          this.appStateService.set("isMatchInProgress", true)
        }
        if (!isMatchInProgress && wasMatchInProgress) {
          Toastr.success("¡La partida ha terminado!")
          this.appStateService.set("isMatchInProgress", false)
        }
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
    }

    private setupUI() {
      const keyHandler = (event: JQueryEventObject) => {
        if(event.key === 'Enter') { // following convention at https://developer.mozilla.org/en-US/docs/Web/API/KeyboardEvent/key#examples
            $("#chat-press-enter-hint").hide();
            this.$timeout(() => {
              if (this.isChatEnabled && this.canSendMessage(this.playerMessage)) {
                this.sendAndCleanMessage(this.playerMessage);
              }
            });
            event.preventDefault();
        }

        // TODO (check): if I leave pressed down the 'x' then at some time errors ocurrs in  func (match *Match) Drop(action PlayerDropAction) PlayerAction at referre.go!, there are concurrent map writes....
        if (event.key === 'x') { // helper code for dev purposes
          if (this.isMatchInProgress) {
            const handCards = this.game.currentMatch.matchCards.byPlayer[Players.generateUniqueKey(this.player)].hand
            if (handCards.length > 0) {
              this.selectedHandCard = handCards[0]
              this.performDropAction();
            }
          }
        }
      };
      this.$rootElement.bind("keydown", keyHandler)
      this.$scope.$on('$destroy', () => {
        this.$rootElement.unbind("keydown", keyHandler)
        // TODO: Implement leave game API CALL in order to remove player from the player list
      });
      const chatMessageInput = document.getElementById("chat-message-input")
      chatMessageInput.focus();

      // event binding on dynamically created, "live" watching  https://stackoverflow.com/a/1207393/903998
      $("div.game-match-section").on("mouseover mouseout","div.play-section .card-image",(event: Event) => {
        if (event.type === "mouseover") { // taken inspiration from https://stackoverflow.com/a/13504775/903998
          var rotateDegress = Math.random() * 10 - Math.random() * 5;
          $((<any>(event.target)).parentElement.parentElement).css('transform', 'rotate(' + rotateDegress + 'deg) scale(1.25)');
        } else {
          $((<any>(event.target)).parentElement.parentElement).css('transform', 'none');
        }
      })

      // UI/UX fine tune of vertical-menu: stop click event propagation so the parents doesn't get the event
      $("div.header ul li ul li").click((e) => {
        e.stopPropagation();
      });

      // dev notes: above handler is not needed as there is no UI control to transition to another ui state when match is in progress
      /*this.$scope.$on('$stateChangeStart', (event, toState, toParams, fromState, fromParams) => {
        if (this.isMatchInProgress) {
          const quit = confirm("Si te vas en plena partida, vas a cagarle la partida a los demás. ¿Estás seguro de irte?")
          if (!quit) {
            event.preventDefault();
          } else {
            this.appStateService.set("isMatchInProgress", false)
            this.appStateService.set("currentGame", null)
          }
        }
      });*/

      // TODO : this hook doesn't works well as the other beforeunload in websockets.ts release the web socket! each beforeunload execute independently... so... a sync mechanish must be employed... perhaps adding to much complexity to something not required (disposing resource on exit tab, could be done by the browser I guess...)
      /*this.$window.addEventListener("beforeunload",(event : any) => {
        if (this.isMatchInProgress) {
          event.preventDefault();
          return event.returnValue = "Si te vas en plena partida, vas a cagarle la partida a los demás. ¿Estás seguro de irte?";
        }
        return undefined
      });*/
    }

    private setupWebsockt(webSocket: WebSocket) {
      const onmessage = (event: MessageEvent<any>) => {
        const notification : {
          kind: string,
          data: {
            game: Api.Game,
            action?: Api.PlayerAction // drop/quit message
            message?: Games.VolatileMessage // chat message
            player?: Api.Player // join/quit message
          };
        } = JSON.parse(event.data)

        switch (notification.kind) {
          case "drop":
          case "take":
            const isClientOtherPlayer = notification.data.action.player.id !== this.player.id
            if (isClientOtherPlayer) {
              this.displayAction(notification.data.action) // is designed to display actions performed by other users
            }
            this.setGame(notification.data.game)
            break
          case "start":
            this.setGame(notification.data.game) // no need to update UI managed by angular
            break;
          case "join":
          case "quit":
            this.$timeout(() => {
              this.setGame(notification.data.game)
            }) // need to update  UI managed by angular as this.game.players IS updated!
            break;
          case "game-chat":
            this.displayMessage(notification.data.message)
            break;
          default:
            console.log("not processing ", notification?.kind)
            break;
        }
      }
      webSocket.addEventListener("message", onmessage)
      this.$scope.$on('$destroy', () => {
        this.gamesService.unbindWebSocket(this.game.id)
        webSocket.removeEventListener("message", onmessage)
      })
    }

    private finalize() {
      // implementent code regarding exiting this screen
    }

    private displayMessage(message: Games.VolatileMessage) {
      const player = message.player
      const $elem = Toastr.chat(player.name,message.text)
      const fontSize = this.getFontSize(player)
      $(".toasrt-chat-message",$elem).css("font-size", fontSize + "px")
    }

    private displayAction(action: Api.PlayerAction) {
      const options: ToastrOptions = {
        timeOut: 5000,
        toastClass: "toastr-info-action-class",
        closeButton: true,
      }
      Toastr.info(`${this.generateActionDescription(action)}`, options)
    }

    private generateActionDescription(action: Api.PlayerAction) {
      if (Matchs.isTakeAction(action)) {
        const cardsText = _.map(action.boardCards,(card) => this.cardToText(card))
        var description = `<span>levantó ${cardsText.join(",")} usando ${this.cardToText(action.handCard)}`
        if (action.isEscobita) {
          const emogis = ["😩","😫","😵","😖"]
          const emogi = emogis[Math.round(Math.random() * emogis.length)]
          description += `&nbsp;<b>e hizó escobita!</b>&nbsp;<span style="font-size: 2rem; vertical-align: middle;">${emogi}</span>`
        }
        description += "</span>"
        return `<b class="player-name">${action.player.name}</b> ${description}`
      } else {
        return `<b>${action.player.name}</b> descartó ${this.cardToText(action.handCard)}`
      }
    }

    private cardToText(card: Api.Card) {
      return `${card.rank} de ${Cards.Suits.translate(card.suit)}`
    }

    private currentSendAndCleanMessage: (msg: Api.Message | Games.VolatileMessage) => void = this.hidePromoteChatUsageElementAndSendAndCleanMessage.bind(this)
    public sendAndCleanMessage(msg: Api.Message | Games.VolatileMessage) {
      this.currentSendAndCleanMessage(msg)
    }

    private hidePromoteChatUsageElementAndSendAndCleanMessage(msg: Api.Message | Games.VolatileMessage): void {
      $(".promote-chat-usage").toggleClass("not-visible");
      this.doSendAndCleanMessage(msg);
      this.currentSendAndCleanMessage = this.doSendAndCleanMessage
    }


    private doSendAndCleanMessage(msg: Api.Message | Games.VolatileMessage) {
      this.sendingMessage = true;
      const sendMessagePromise = this.gamesService.sendMessage(msg as Games.VolatileMessage).then(() => {
        this.cleanMessage(msg);
      }).finally(() => {
        this.sendingMessage = false;
      })
      this.allowSendMessage = false;
      this.$timeout(() => {
        this.allowSendMessage = true;
      }, 2000)
      return sendMessagePromise
    }

    public canSendMessage(msg: Api.Message | Games.VolatileMessage) {
      return !this.loading && this.allowSendMessage && !this.sendingMessage && !_.isEmpty(msg.text);
    }

    private cleanMessage(msg: Api.Message | Games.VolatileMessage) {
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
      return Matchs.Rules.isValidTakeAction(this.selectedHandCard,selectedBoardCards)
    }

    public performTakeAction() {
      const selectedBoardCards = this.getSelectedBoardCards()
      const takeAction = Matchs.createTakeAction(this.player,selectedBoardCards,this.selectedHandCard)
      this.loading = true;
      this.gamesService.performTakeAction(this.game,takeAction).then((data) => {
        if (data.action.isEscobita) {
          Toastr.success("Has hecho escoba! 🥳")
        }
        //this.setGame(data.game) // don't need to refresh as this client gets notified via ws
      }).finally(() => {
        this.selectedHandCard = undefined
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
      const dropAction = Matchs.createDropAction(this.player,this.selectedHandCard)
      this.loading = true;
      this.gamesService.performDropAction(this.game,dropAction).then((data) => {
        //this.setGame(data.game) // don't need to refresh as this client gets notified via ws
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
        const hasPreviousMatchs = (currentMatchIndex > 0)
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
        this.currentMatchStatsByPlayerUniqueKey = stats;
        this.currentPositionByPlayerUniqueKey = Matchs.Rules.calculatePositionByPlayerUniqueKey(stats)
        this.currentFontSizeByPlayerUniqueKey = UIMessages.calculateFontSizeByPlayerUniqueKey(this.currentPositionByPlayerUniqueKey)
        return stats
      })
    }

    public getFontSize(player: Players.Player) {
      if (_.isEmpty(this.currentFontSizeByPlayerUniqueKey)) {
        return UIMessages.minFontSize;
      } else {
        return this.currentFontSizeByPlayerUniqueKey[Players.generateUniqueKey(player)]
      }
    }

    private suggestionRequestCount: number = 0;
    public static maxSuggestionRequestCount = 3;
    public canRequestTakeActionsSuggestion() {
      return this.calculateRemainderTakeActionSuggestions() > 0
    }

    public calculateRemainderTakeActionSuggestions() {
      return Controller.maxSuggestionRequestCount - this.suggestionRequestCount
    }

    public possibleTakeActions: Matchs.Engine.SuggestedTakeAction[];
    public recomendedTakeActionTakeAction: Matchs.Engine.SuggestedTakeAction;
    public requestTakeActionsSuggestion() {
      var boardCards = this.game.currentMatch.matchCards.board
      if (_.size(boardCards) > 8) {
        boardCards = Matchs.Engine.getMostImportantCards(boardCards, 8)
      }

      this.suggestionRequestCount++
      const handCards = this.game.currentMatch.matchCards.byPlayer[Players.generateUniqueKey(this.player)].hand
      const possibleTakeActions = Matchs.Engine.calculatePossibleTakeActions(boardCards, handCards, this.player)
      const analizedActions = Matchs.Engine.analizeActions(possibleTakeActions, this.game.currentMatch)
      this.possibleTakeActions = analizedActions.possibleActions;
      this.recomendedTakeActionTakeAction = analizedActions.recomendedAction;

      this.playerMessage.text = "Soy 💩 y pido sugerencias al escoba master";
      this.loading = true;
      this.doSendAndCleanMessage(this.playerMessage).finally(() => {
        this.loading = false;
      })
    }

    public openPossibleTakeActionsDialog() {
      const dialog = document.getElementById('suggested-take-actions-dialog') as HTMLDialogElement;
      dialog.showModal()
    }

    public closePossibleTakeActionsDialog() {
      const dialog = document.getElementById('suggested-take-actions-dialog') as HTMLDialogElement;
      dialog.close();
    }

    public botCount: number = 0
    public addComputerPlayer(game: Games.Game) {
      this.loading = true;
      this.gamesService.addComputerPlayer(game).then(() => {
        Toastr.info("Computer added")
        this.botCount++
      }).finally(() => {
        this.loading = false;
      })
    }

    public isPlayerGameOwner = Games.isPlayerOwner
    public playerToUniqueKey = Players.generateUniqueKey // interesting case! both names are pretty same. the player context is provivded by the leading namespace or the inclusion in the name's identifier
    public extractPlayerName = Players.extractName
    public hasGameStarted = Games.isStarted
  }

  escobita.controller('GameController', ['$rootElement','$rootScope','$scope','$state', 'GamesService', 'WebSocketsService', '$timeout', '$q', '$window', 'AppStateService', Controller]);
}
