/// <reference path='../app.ts' />
/// <reference path='../api-model.ts' />



namespace Games {

  export interface Game extends Api.Game {
  }

  export function isStarted(game :Game) : boolean{
    return false
  }

  export function addPlayer(game :Game, player: Players.Player) {
    if (_.isEmpty(game.players)) {
      game.players = [player]
    } else {
      const gamePlayer = _.find(game.players,(gamePlayer) => gamePlayer.id == player.id)
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

  }

  escobita.service('GamesService', ['$http', '$q', Service]);
}


namespace Players {

  export interface Player extends Api.Player {
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