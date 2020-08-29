/// <reference path='../app.ts' />
/// <reference path='../services/services.ts' />

module Game {

  // BEGIN : move to common.js
  function isNumeric(s: any): boolean {
    return !isNaN(parseFloat(s)) && isFinite(s); // based from: http://stackoverflow.com/a/6449623
  }

  function unixToReadableClock(unix: number): string {
    return formatUnixTimestamp(unix,"HH:mm")
  }

  function unixToReadableDay(unix: number): string {
    return formatUnixTimestamp(unix,"DD/MM/YYYY")
  }

  function unixToReadableDate(unix: number): string {
    return formatUnixTimestamp(unix,"DD/MM/YYYY HH:mm")
  }

  function unixToReadableDateVerbose(unix: number): string {
    return formatUnixTimestamp(unix,"dddd DD/MM/YYYY [a las] HH:mm")
  }


  function formatUnixTimestamp(unix: number, layout:string) {
    if (isNumeric(unix)) {
      return moment.unix(unix).format(layout);
    } else {
      console.warn(`Unix timestamp value(=${unix}) is not a number`)
      return '';
    }
  }
  // END :  move to util.js

  class Controller {

    public game: Games.Game; // the current game
    public player: Players.Player; // the client player
    public isPlayerTurn: boolean;
    public currentTurnPlayer: Players.Player; // the player that acts in the current turn
    public messages: Messages.Message[]; // all from the server related to this game
    public isBoardCardSelectedById: _.Dictionary<boolean>;
    public selectedHandCard: Api.Card;

    public message: Messages.Message; // buffer for user input
    public disableSendMessageBtn: boolean = false; // avoids multiples clicks!
    public hideChat: boolean = true;
    public matchInProgress: boolean = false;

    public players: Players.Player[];
    public playersById: Util.EntityById<Players.Player>;

    private lastUpdateUnixTimestamp: number = undefined;

    public formatUnixTimestamp = unixToReadableClock

    constructor(private $scope: ng.IScope, private $state: ng.ui.IStateService, private gamesService: Games.Service, private playersService: Players.Service,
      private messagesService: Messages.Service, private $interval: ng.IIntervalService, private $timeout: ng.ITimeoutService,
      private $q: ng.IQService) {
      this.game = $state.params["game"]
      this.player = $state.params["player"]

      /*this.$interval(() => {
        this.updatePlayers()
          .then(() => this.updateGameMessages ())
          .then(() => {
            console.log("Updated players and messages OK!")
            if (!_.isUndefined(this.lastUpdateUnixTimestamp)) {
              const now = 	Math.floor(new Date().getTime()/1000.0)
              console.log("demora aproximada ", now - this.lastUpdateUnixTimestamp)
            }
            this.lastUpdateUnixTimestamp = 	Math.floor(new Date().getTime()/1000.0)
          })
      },2000)*/

      this.$scope.$watch(() => {
        if (_.isUndefined(this.game.currentMatch)) {
          return undefined
        }
        return this.game.currentMatch.currentRound.currentTurnPlayer
      }, (currentTurnPlayer,previousTurnPlayer) => {
        if (!_.isUndefined(currentTurnPlayer)) {
          this.currentTurnPlayer = currentTurnPlayer;
          this.isPlayerTurn = Rounds.isPlayerTurn(this.game.currentMatch.currentRound,this.currentTurnPlayer)
        }
      })
    }

    public updateGameMessages() {
      return this.messagesService.getMessagesByGame(this.game.id).then((messages) => {
        this.messages = messages;
        return messages;
      })
    }

    private updatePlayers() {
      return this.playersService.getPlayers().then((players) => {
        this.playersById = Util.toMapById(players)
        this.players = players;
        return players
      })
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
      // TODO : Use loading flag for UI
      this.gamesService.startGame(game).then((game) => {
        console.log(game);
        this.game = game;
        this.matchInProgress = true;
      })
    }

  }

  escobita.controller('GameController', ['$scope','$state', 'GamesService', 'PlayersService', 'MessagesService', '$interval', '$timeout', '$q', Controller]);
}

// TODO: place in separate file

namespace Util {


  export interface EntityById<T> extends _.Dictionary<T>  {
    [id: number] : T
  }

  /** An identifiable entity has an numeric id that identifies unequivocally within a context. */
  export interface Identificable {
    id?: number; // it is optional due to have api model objects with nullable id, as they first are created in the client and then saved on the server granting an id, basically for avoiding this -> `Property 'id' is optional in type 'Player' but required in type 'Identificable'`
  }

  /** Generates a map by id of the given collection of identificables. */
  export function toMapById<T extends Identificable>(entites: T[]) : EntityById<T> {
    return _.indexBy(entites, 'id');
  }

  /** Generates a map by id of the given collection of elements whose id is extracted using the correspondant method */
  export function toMapByIdUsingGetter<T>(list: T[], idGetterFunc: (elem:T) => number) : EntityById<T> {
    return list.reduce((map: any, elem: T) => {
      map[idGetterFunc(elem)] = elem
      return map;
    }, {});
  }
}