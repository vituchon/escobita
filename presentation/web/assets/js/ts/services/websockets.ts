
// TODO : IMPLEMENT on tab close to close the ws if opened
namespace WebSockets {

  export class Service {

    private webSocket: WebSocket

    constructor(private $http: ng.IHttpService, private $q: ng.IQService) {
    }


    public retrieve() {
      if (_.isUndefined(this.webSocket)) {
        return this.adquire();
      } else {
        if (this.webSocket.readyState === WebSocket.OPEN) {
          return this.$q.when(this.webSocket) // in the case is called more than once, return the ws already established
        } else {
          return this.$q.reject("already binded to another tab")
        }
      }
    }

    private adquire() {
      const deffered = this.$q.defer();
      try {
        this.webSocket = new WebSocket("ws://localhost:9090/adquire-ws"); // TODO : set domain dynamically (NICE: set protocol ws or wss accordingly)

        this.webSocket.onopen = (event : Event) => {
          console.log("Web socket opened, event is: ", event)
          deffered.resolve(this.webSocket)
        }
        this.webSocket.onerror = (event : Event) => {
          console.warn("Web socket error, event is: ", event)
          deffered.reject(event)
        }

      } catch(error) {
        deffered.reject(error)
      }
      return deffered.promise
    }

    private release() {
      return this.$http.get("/release-ws").then(( ) => {
        this.webSocket.close()
        this.webSocket = undefined
      })
    }

  }

  escobita.service('WebSocketsService', ['$http', '$q', Service]);
}