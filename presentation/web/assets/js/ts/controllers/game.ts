/// <reference path='../app.ts' />
/// <reference path='../services/_services.d.ts' />
/// <reference path='../directives/_directives.d.ts' />

module Game {


  namespace UIMessages {

    export const baseFontSize = 12;
    const maxFontSize = 24;
    const fontSizeAmplitude = maxFontSize - baseFontSize;


    function determineFontSize(position: number) {
        const size = baseFontSize + (fontSizeAmplitude) / position; // B + M/X (Homográfica desplazada)
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
    public messages: Messages.Message[]; // all from the server related to this game
    public isBoardCardSelectedById: _.Dictionary<boolean>;
    public selectedHandCard: Api.Card;

    public messageText: string; // buffer for user input
    public disableSendMessageBtn: boolean = false; // avoids multiples clicks!
    public isChatEnabled: boolean = false;
    private currentFontSizeByPlayerName: UIMessages.FontSizeByPlayerName; // funny font size to use by player name
    private currentPositionByPlayerName: Matchs.Rules.PositionByPlayerName; // positions by player name

    public isMatchInProgress: boolean = false;
    public currentMatchStats: Api.ScoreSummaryByPlayerName;

    /* TODO : Implement displaying of previous matchs score, leaving code below as an first approach
    // the two below vars are a simple solution for displaying the last match scores! nothing more.. it may be polished a lot
    public currentMatchStatsCopy: Api.ScoreSummaryByPlayerName;
    public displayCurrentMatchStatsCopy: boolean = false;*/

    public players: Players.Player[]; // not sure if it will be use somewhere!
    public playersById: Util.EntityById<Players.Player>;

    private lastUpdateUnixTimestamp: number = undefined;

    public formatUnixTimestamp = Util.unixToReadableClock
    public translateSuit = Cards.Suits.translate

    public loading: boolean = false;
    public displayCardsAsSprites: boolean = true;

    private updaterInterval: ng.IPromise<any>; // "handler" to the one interval that updates the UI according to the controller's state

    constructor(private $rootScope: ng.IRootElementService, private $scope: ng.IScope, private $state: ng.ui.IStateService, private gamesService: Games.Service, private playersService: Players.Service,
      private messagesService: Messages.Service, private $interval: ng.IIntervalService, private $timeout: ng.ITimeoutService,
      private $q: ng.IQService) {
      this.game = $state.params["game"]
      this.player = $state.params["player"]
      this.isPlayerGameOwner = Games.isPlayerOwner(this.player,this.game)

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

      this.$scope.$on('$destroy', () => {
        this.$interval.cancel(this.updaterInterval)
      });

      this.updaterInterval = this.$interval(() => {
        const mustRefreshGame = (!this.isMatchInProgress) || (this.isMatchInProgress && !this.isPlayerTurn) // ~p v (p y q~) == ~p v ~q
        var refreshGamePromise = this.$q.when(this.game)
        if (mustRefreshGame) {
          refreshGamePromise = this.refreshGame()
        }
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
      }, 2000)

    /*
      this.$scope.$watch(() => {
        return this.refreshGameInterval
      }, (refreshGameInterval) => {
        const mustHaveARefresh = (!this.isMatchInProgress) || (this.isMatchInProgress && !this.isPlayerTurn) // ~p v (p y q~) == ~p v ~q
        if (mustHaveARefresh && _.isUndefined(refreshGameInterval)) {
          // console.warn("must have automatic refresh!")
          // in some cases when another player start the games, the watch this.isPlayerTurn executes before the watch over this.isMatchInProgress; thus disabling the match
          this.refreshGameInterval = this.$interval(() => {
            return this.refreshGame()
          },2000)
        }
      })*/

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

    public updateChat() {
      return this.playersService.getPlayers().then((players) => {
        this.players = players
        this.playersById = Util.toMapById(this.players);
        return this.players
      }).then((players) => {
        return this.messagesService.getMessagesByGame(this.game.id).then((messages) => {
          const incomingMessages = this.determineIncomingMessages(this.messages,messages)
          _.forEach(incomingMessages,(incomingMessage) => {
            const player = this.playersById[incomingMessage.playerId]
            const $elem = Toastr.chat(player.name,incomingMessage.text)
            const fontSize = this.getFontSize(player)
            $(".toasrt-chat-message",$elem).css("font-size", fontSize + "px");
          })
          this.messages = messages;
          return messages;
        })
      })
    }

    private determineIncomingMessages(received: Api.Message[], all: Api.Message[]) {
      if (_.isEmpty(received)) {
        return all
      } else {
        return _.filter(all,(message) => {
          const machtingMessage = _.find(received,(recievedMessage) => recievedMessage.id === message.id)
          const notReceived = _.isUndefined(machtingMessage)
          return notReceived
        })
      }
    }

    public sendMessage(text: string) {
      if (this.disableSendMessageBtn) {
        return
      }
      const message = Messages.newMessage(this.game.id, this.player.id, text)
      this.messagesService.createMessage(message)
      this.disableSendMessageBtn = true;
      this.$timeout(() => {
        this.disableSendMessageBtn = false;
      }, 2000)
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
        return this.setGame(data.game)
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
        return this.setGame(data.game)
      }).finally(() => {
        this.selectedHandCard = undefined
        this.loading = false;
      })
    }

    private setGame(game: Games.Game) {
      this.game = game;
      this.isMatchInProgress = Games.hasMatchInProgress(game)
      if (this.isMatchInProgress) {
        this.currentTurnPlayer = this.game.currentMatch.currentRound.currentTurnPlayer;
        this.isPlayerTurn = Rounds.isPlayerTurn(this.game.currentMatch.currentRound,this.player)
        return this.updateGameStats().then(() => {
          return this.game
        })
      } else {
        return this.$q.when(this.game)
      }
    }

    public refreshGame() {
      return this.gamesService.getGameById(this.game.id).then((game) => {
        return this.setGame(game)
      })
    }

    private updateGameStats() {
      //this.currentMatchStatsCopy = angular.copy(this.currentMatchStats)
      return this.gamesService.calculateStatsByGameId(this.game.id).then((stats) => {
        this.currentMatchStats = stats;
        this.currentPositionByPlayerName = Matchs.Rules.calculatePositionByPlayerName(stats)
        this.currentFontSizeByPlayerName = UIMessages.calculateFontSizeByPlayerName(this.currentPositionByPlayerName)
        return stats
      })
    }

    public getFontSize(player: Players.Player) {
      if (_.isEmpty(this.currentFontSizeByPlayerName)) {
        return UIMessages.baseFontSize;
      } else {
        return this.currentFontSizeByPlayerName[player.name]
      }
    }


  }

  escobita.controller('GameController', ['$rootScope','$scope','$state', 'GamesService', 'PlayersService', 'MessagesService', '$interval', '$timeout', '$q', Controller]);
}
