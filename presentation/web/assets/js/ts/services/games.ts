/// <reference path='./cards.ts' />

namespace Games {

  export interface Game extends Api.Game {  // TODO : analyse if this approach is worty...
  }

  export function hasMatchInProgress(game :Game) : boolean {
    return !_.isUndefined(game.currentMatch)
  }

  export function isStarted(game: Game): boolean {
    return _.size(game.matchs) > 0 || hasMatchInProgress(game)
  }

  export function canDeleteGame(game :Game, player: Players.Player) {
    return isPlayerOwner(game,player)
  }

  export function isPlayerOwner(game :Game, player: Players.Player) {
    if (Util.isDefined(player.id)) {
      return player.id === game.owner.id
    } else {
      return player.name === game.owner.name // Dev notes (Loop hole here!) : recall that in a game players MUST have different names...
    }
  }

  export function hasPlayerJoin(game :Game, player: Players.Player) {
    return _.findIndex(game.players,(gamePlayer) => gamePlayer.id === player.id) !== -1
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


  export interface VolatileMessage {
    gameId: number;
    player: Api.Player;
    text: string;
  }

  export function isVolatile(message: any): message is VolatileMessage {
    return Util.isDefined(message.player)
  }

  export function newMessage(gameId: number, player: Api.Player, text: string): VolatileMessage {
    return {
      gameId: gameId,
      player: player,
      text: text,
    }
  }

  function enhanceMaps(game: Api.Game) {
    if (Util.isDefined(game.currentMatch)) {
      game.currentMatch.actionsByPlayerName = new Players.MapByPlayer<Api.PlayerAction>(<any>game.currentMatch.actionsByPlayerName)
      game.currentMatch.matchCards.byPlayerName = new Players.MapByPlayer<Api.PlayerMatchCards>(<any>game.currentMatch.matchCards.byPlayerName)
    }

    if (_.size(game.matchs) > 0) {
      for (let i = 0; i < game.matchs.length; i++) {
        game.matchs[i].actionsByPlayerName = new Players.MapByPlayer<Api.PlayerAction>(<any>game.matchs[i].actionsByPlayerName)
        game.matchs[i].matchCards.byPlayerName = new Players.MapByPlayer<Api.PlayerMatchCards>(<any>game.matchs[i].matchCards.byPlayerName)
      }
    }

    return game
  }
  export class Service {
    constructor(private $http: ng.IHttpService, private $q: ng.IQService) {
    }

    // TODO: remove Game or Games letters from these methods, as Games. provides already context (avoid slutering)
    getGames(): ng.IPromise<Game[]> {
      return this.$http.get<Game[]>(`/api/v1/games`).then((response) => {
        const games = response.data
        //return games
        return games.map(game => enhanceMaps(game));
      });
    }

    getGameById(gamesId: number): ng.IPromise<Game> {
      return this.$http.get<Game>(`/api/v1/games/${gamesId}`).then((response) => {
        const game = response.data
        //return game
        return enhanceMaps(game);
      });
    }

    createGame(game: Game): ng.IPromise<Game> {
      return this.$http.post<Game>(`/api/v1/games`,game).then((response) => {
        const game = response.data
        return enhanceMaps(game);
      })
    }

    /*updateGame(game: Game): ng.IPromise<Game> {
      return this.$http.put<Game>(`/api/v1/games/${game.id}`,game).then((response) => {
        const game = response.data
        return enhanceMaps(game);
      })
    }*/

    deleteGame(game: Game, player: Players.Player): ng.IPromise<any> {
      const config: ng.IRequestShortcutConfig = {
        data: player
      };
      return this.$http.delete<any>(`/api/v1/games/${game.id}`, config).then((_) => {
        return null;
      })
    }

    startGame(game: Game): ng.IPromise<Game> {
      return this.$http.post<Game>(`/api/v1/games/${game.id}/resume`,game).then((response) => {
        const game = response.data
        return enhanceMaps(game);
      })
    }

    joinGame(game: Game): ng.IPromise<Game> {
      return this.$http.post<Game>(`/api/v1/games/${game.id}/join`,undefined).then((response) => {
        const game = response.data
        return enhanceMaps(game);
      })
    }

    quitGame(game: Game): ng.IPromise<Game> {
      return this.$http.post<Game>(`/api/v1/games/${game.id}/quit`,undefined).then((response) => {
        const game = response.data
        return enhanceMaps(game);
      })
    }

    addComputerPlayer(game: Game): ng.IPromise<Game> {
      return this.$http.post<Game>(`/api/v1/games/${game.id}/add-computer`,undefined).then((response) => {
        const game = response.data
        return enhanceMaps(game);
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
      return this.$http.get<Api.ScoreSummaryByPlayerUniqueKey>(`/api/v1/games/${id}/calculate-stats`,config).then((response) => {
        return response.data
      })
    }

    bindWebSocket(id: number) {
      return this.$http.get<Game>(`/api/v1/games/${id}/bind-ws`).then((response) => {
        const game = response.data
        return enhanceMaps(game);
      });
    }

    unbindWebSocket(id: number) {
      return this.$http.get<Game>(`/api/v1/games/${id}/unbind-ws`).then((response) => {
        const game = response.data
        return enhanceMaps(game);
      });
    }

    sendMessage(msg: VolatileMessage) {
      return this.$http.post<Game>(`/api/v1/games/${msg.gameId}/message`, msg).then((response) => {
        const game = response.data
        return enhanceMaps(game);
      })
    }
  }

  escobita.service('GamesService', ['$http', '$q', Service]);
}