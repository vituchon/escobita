
namespace WebSockets {

  export class Service {

    private webSocket: WebSocket

    constructor(private $http: ng.IHttpService, private $q: ng.IQService) {
    }


    public retrieve() {
      return this.adquire().then((ws) => {
        return ws
      }).catch(() => {
        return this.release().then(() => { // if fails then try to release and re create a new one!
          return this.adquire()
        })
      })
    }

    private adquire() {
      if (!_.isUndefined(this.webSocket)) {
        return this.$q.when(this.webSocket)
      }
      const deffered = this.$q.defer();
      try {
        this.webSocket = new WebSocket("ws://localhost:9090/adquire-ws");
        deffered.resolve(this.webSocket)
      } catch(error) {
        deffered.reject(error)
      }
      return deffered.promise
    }

    private release() {
      if (!_.isUndefined(this.webSocket)) {
        return this.$http.get("/release-ws").then(( ) => {
          this.webSocket.close()
          this.webSocket = undefined
        })
      } else {
        return this.$q.when(undefined)
      }
    }

  }

  escobita.service('WebSocketsService', ['$http', '$q', Service]);
}