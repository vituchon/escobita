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

    public set(key: string, value:any) {
      this.state[key] = value;
    };

    public get(key: string) {
      return this.state[key]
    };

    public clear() {
      this.state = {};
    };
  }

  escobita.service('AppState', [Service]);
}
