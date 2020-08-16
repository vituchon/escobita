/// <reference path='../app.ts' />
/// <reference path='../api-model.ts' />



namespace Games {

  export interface Game extends Api.Game {
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

    getGameById(gamesId: number): ng.IPromise<Game[]> {
      return this.$http.get<Game[]>(`/api/v1/games/${gamesId}`).then((response) => {
        return response.data;
      });
    }

    createGame(game: Game): ng.IPromise<Game> {
      return this.$http.post<Game>(`/api/v1/games`,game).then((response) => {
        return response.data
      }).catch((err) => {
        console.log(err)
        throw err
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

var gameService: Games.Service;
// TODO : remove when dev finish
escobita.run(['GamesService', (gs: Games.Service) => {
  gameService = gs
}])