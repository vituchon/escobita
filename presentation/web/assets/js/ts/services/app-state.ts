/// <reference path='../app.ts' />
/// <reference path='../api-model.ts' />


namespace AppState {

  export function isPlayerRegistered(player: Api.Player) {
    return Util.isDefined(player?.id) && !_.isEmpty(player?.name)
  }

  export class Service {
    private state: {
      [key:string]: any;
    }
    constructor() {
      this.clear();
    }

    public set<T>(key: string, value:T) {
      this.state[key] = value;
    };

    public get<T>(key: string): T {
      return this.state[key]
    };

    public clear() {
      this.state = {};
    };

    public isClientPlayerRegistered(): boolean {
      const clientPlayer = this.get<Api.Player>("clientPlayer")
      return isPlayerRegistered(clientPlayer)
    }
  }

  escobita.service('AppStateService', [Service]);
}
