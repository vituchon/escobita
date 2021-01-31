/// <reference path='./cards.ts' />

namespace Games {

  export interface Game extends Api.Game {
  }

  export function hasMatchInProgress(game :Game) : boolean {
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

  export function isPlayerOwner(player: Players.Player, game :Game) {
    if (_.isUndefined(game.playerId)) {
      console.warn("suspicious things, programmer must check something...")
    }
    return player.id === game.playerId
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

  export interface GameActionResponse {
    game: Api.Game;
    action: Api.PlayerTakeAction;
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
      return this.$http.post<Game>(`/api/v1/games/${game.id}/resume`,game).then((response) => {
        return response.data
      })
    }

    performTakeAction(game: Game,takeAction: Api.PlayerTakeAction) {
      return this.$http.post<GameActionResponse>(`/api/v1/games/${game.id}/perform-take-action`,takeAction).then((response) => {
        return response.data
      })
    }

    performDropAction(game: Game,dropAction: Api.PlayerDropAction) {
      return this.$http.post<GameActionResponse>(`/api/v1/games/${game.id}/perform-drop-action`,dropAction).then((response) => {
        return response.data
      })
    }

    calculateStatsByGameId(id: number, matchIndex: number) {
      const config: ng.IRequestShortcutConfig = {
        params: {
          matchIndex: matchIndex
        }
      };
      return this.$http.get<Api.ScoreSummaryByPlayerName>(`/api/v1/games/${id}/calculate-stats`,config).then((response) => {
        return response.data
      })
    }

    bindWebSocket(id: number, webSocket: WebSocket) { // TODO : not sure why the second parameter..
      return this.$http.get<Game>(`/api/v1/games/${id}/bind-ws`).then((response) => {
        return response.data;
      });
    }


  }

  escobita.service('GamesService', ['$http', '$q', Service]);
}