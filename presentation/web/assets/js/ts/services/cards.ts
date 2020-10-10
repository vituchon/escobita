/// <reference path='../app.ts' />
/// <reference path='../api-model.ts' />

namespace Cards {

  export namespace Suits {
    // dev notes: the values must match some at /local/escobita/model/card.go#Line:25
    export const sword = 0
    export const club = 1
    export const cup = 2
    export const gold = 3
    export const all = [sword, club, cup, gold]

    export const labels: _.Dictionary<string> = {
      [sword]: "espada",
      [club]: "palo",
      [cup]: "copa",
      [gold]: "oro",
    }

    export function translate(suit: number) {
      return labels[suit]
    }
  }
}

namespace Rounds {

  export interface Round extends Api.Round {
  }

  export function isPlayerTurn(round: Round, player: Players.Player): boolean {
    return !_.isEmpty(round.currentTurnPlayer) && round.currentTurnPlayer.name == player.name;
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
      })
      // TODO: below is not used code, server will return a 500 is a client seek for a non existing client, guess a 404 is more correct...
      /*.catch((err) => {
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
