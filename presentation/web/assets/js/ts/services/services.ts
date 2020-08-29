/// <reference path='../app.ts' />
/// <reference path='../api-model.ts' />



namespace Games {

  export interface Game extends Api.Game {
  }

  export function hasMatchInProgress(game :Game) : boolean{
    return !_.isUndefined(game.currentMatch)
  }

  export function addPlayer(game :Game, player: Players.Player) {
    if (_.isEmpty(game.players)) {
      game.players = [player]
    } else {
      // it allow to share screen for those users with same name... could result in unexpected behaviour, although it may be very funny!
      const gamePlayer = _.find(game.players,(gamePlayer) => gamePlayer.name == player.name)
      const playerNotJoined = _.isUndefined(gamePlayer)
      if (playerNotJoined) {
        game.players.push(player)
      }
    }
  }

  export namespace Periods {
    const dayInSeconds = 24 * 60 * 60;
    const weekInSeconds = dayInSeconds * 7;

    //  Well know values
    export enum Values {
      Daily = dayInSeconds,
      Weekly = weekInSeconds,
     };

     // common legend to display for each value, it shall be accessed via the well know values
    export const labels = Object.freeze({
      [Values.Daily] : "Diaria",
      [Values.Weekly] : "Semanal"
    });

    export function getPeriodDescription(period: Periods.Values): string {
      return (<any>Periods.labels)[period]
    }
  }

  export class Service {
    constructor(private $http: ng.IHttpService, private $q: ng.IQService) {
    }

    getGames(): ng.IPromise<Game[]> {
      return this.$http.get<Game[]>(`/api/v1/games`).then((response) => {
        return response.data;
      });
    }

    getGameById(gamesId: number): ng.IPromise<Game> {
      return this.$http.get<Game>(`/api/v1/games/${gamesId}`).then((response) => {
        return response.data;
      });
    }

    createGame(game: Game): ng.IPromise<Game> {
      return this.$http.post<Game>(`/api/v1/games`,game).then((response) => {
        return response.data
      })
    }

    updateGame(game: Game): ng.IPromise<Game> {
      return this.$http.put<Game>(`/api/v1/games/${game.id}`,game).then((response) => {
        return response.data
      })
    }

    deleteGame(game: Game): ng.IPromise<any> {
      return this.$http.delete<any>(`/api/v1/games/${game.id}`).then((_) => {
        return null;
      })
    }

    startGame(game: Game): ng.IPromise<Game> {
      return this.$http.post<Game>(`/api/v1/games/${game.id}/start`,game).then((response) => {
        return response.data
      })
    }

  }

  escobita.service('GamesService', ['$http', '$q', Service]);
}

namespace Rounds {

  export interface Round extends Api.Round {
  }

  export function isPlayerTurn(round: Round, player: Players.Player): boolean {
    return round.currentTurnPlayer.name == player.name;
  }
}

namespace Players {

  export interface Player extends Api.Player {
  }

  export class Service {
    constructor(private $http: ng.IHttpService, private $q: ng.IQService) {
    }

    getPlayers(): ng.IPromise<Player[]> {
      return this.$http.get<Player[]>(`/api/v1/players`).then((response) => {
        return response.data;
      });
    }

    getPlayerById(playersId: number): ng.IPromise<Player> {
      return this.$http.get<Player>(`/api/v1/players/${playersId}`).then((response) => {
        return response.data;
      })/*.catch((err) => {
        if (err.status === 404) { // not found: there is no player with the given id
          return null
        } else {
          return err
        }
      })*/
    }

    getClientPlayer(): ng.IPromise<Player> {
      return this.$http.get<Player>(`/api/v1/player`).then((response) => {
        return response.data;
      })
    }

    createPlayer(player: Player): ng.IPromise<Player> {
      return this.$http.post<Player>(`/api/v1/players`,player).then((response) => {
        return response.data
      })
    }

    updatePlayer(player: Player): ng.IPromise<Player> {
      return this.$http.put<Player>(`/api/v1/players/${player.id}`,player).then((response) => {
        return response.data
      })
    }

    deletePlayer(player: Player): ng.IPromise<any> {
      return this.$http.delete<any>(`/api/v1/players/${player.id}`).then((_) => {
        return null;
      })
    }

  }

  escobita.service('PlayersService', ['$http', '$q', Service]);
}


namespace Messages {

  export interface Message extends Api.Message {
  }

  export function newMessage(gameId: number, playerId: number, text: string): Message {
    return {
      gameId: gameId,
      playerId: playerId,
      text: text,
    }

  }

  export class Service {
    constructor(private $http: ng.IHttpService, private $q: ng.IQService) {
    }

    getMessages(): ng.IPromise<Message[]> {
      return this.$http.get<Message[]>(`/api/v1/messages`).then((response) => {
        return response.data;
      });
    }

    getMessageById(messagesId: number): ng.IPromise<Message> {
      return this.$http.get<Message>(`/api/v1/messages/${messagesId}`).then((response) => {
        return response.data;
      })/*.catch((err) => {
        if (err.status === 404) { // not found: there is no message with the given id
          return null
        } else {
          return err
        }
      })*/
    }

    getMessagesByGame(gameId: number): ng.IPromise<Message[]> {
      return this.$http.get<Message[]>(`/api/v1/messages/get-by-game/${gameId}`).then((response) => {
        return response.data;
      })
    }


    getClientMessage(): ng.IPromise<Message> {
      return this.$http.get<Message>(`/api/v1/message`).then((response) => {
        return response.data;
      })
    }

    createMessage(message: Message): ng.IPromise<Message> {
      return this.$http.post<Message>(`/api/v1/messages`,message).then((response) => {
        return response.data
      })
    }

    updateMessage(message: Message): ng.IPromise<Message> {
      return this.$http.put<Message>(`/api/v1/messages/${message.id}`,message).then((response) => {
        return response.data
      })
    }

    deleteMessage(message: Message): ng.IPromise<any> {
      return this.$http.delete<any>(`/api/v1/messages/${message.id}`).then((_) => {
        return null;
      })
    }

  }

  escobita.service('MessagesService', ['$http', '$q', Service]);
}