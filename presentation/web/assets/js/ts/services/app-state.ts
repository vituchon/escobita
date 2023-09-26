/// <reference path='../app.ts' />
/// <reference path='../api-model.ts' />


namespace AppState {

  export class Service {
    private state: {
      [key:string]: any;
    }
    constructor() {
      console.log("App state set")
      this.clear();
    }

    public incByOne(key: string, startFrom: number = 1) {
      var value =  this.state[key]
      if (_.isUndefined(value)) {
        value = startFrom
      }
      return this.state[key] = value
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
  }

  escobita.service('AppStateService', [Service]);
}
