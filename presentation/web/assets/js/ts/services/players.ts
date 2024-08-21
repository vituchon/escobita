/// <reference path='../app.ts' />
/// <reference path='../api-model.ts' />

namespace Players {

  // Dev notes: Would be the client side function of this back side function => model/player.go#MarshalText
  const playerFieldSeparator = "|"
  export function generateUniqueKey(player: Api.Player) {
    return player.id + playerFieldSeparator + player.name
  }

  export function extractName(playerKey: string): string {
    return playerKey?.split?.(playerFieldSeparator)?.[1]
  }

  export interface Player extends Api.Player {  // TODO : analyse if this approach is worty...
  }

  // a decorated object that facilitates access an object used as map, transforim the player into a key using the transformation function "generateUniqueKey"
  export class MapByPlayer<T> implements Api.MapByPlayer<T>, Object {
    [key: string]: T | any;

    constructor(object: any) {
      if (Util.isDefined(object)) {
        const props = Object.getOwnPropertyNames(object)
        _.forEach(props,(prop) => {
          this[prop] = object[prop]
        })
      }
    }

    public set(player: Player, data: T) {
      this[generateUniqueKey(player)] = data
    }

    public get(player: Player): T {
      return this[generateUniqueKey(player)]
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
