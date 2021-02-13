


namespace WebSockets {
  interface CloseEvent extends Event {
    code: number;
  }

  const closeDescriptionByCode: {[code:number]: string} = { // taken from https://github.com/Luka967/websocket-close-codes
    1000:	"CLOSE_NORMAL",
    1001:	"CLOSE_GOING_AWAY",
    1002:	"CLOSE_PROTOCOL_ERROR",
    1003:	"CLOSE_UNSUPPORTED",
    1004:	"Reserved",
    1005:	"CLOSED_NO_STATUS",
    1006:	"CLOSE_ABNORMAL",
    1007:	"Unsupported payload",
    1008:	"Policy violation",
    1009:	"CLOSE_TOO_LARGE",
    1010:	"Mandatory extension",
    1011:	"Server error",
    1012:	"Service restart",
    1013:	"Try again later",
    1014:	"Bad gateway",
    1015:	"TLS handshake fail"
  }
  export class Service {

    private webSocket: WebSocket

    constructor(private $http: ng.IHttpService, private $q: ng.IQService, private $window: ng.IWindowService) {
    }


    public retrieve() {
      if (_.isUndefined(this.webSocket)) {
        return this.adquire().catch((err) => {
          console.warn("adquire fails err: ", err, " trying to release and readquire (discard the last and create a new upgraded conn)")
          return this.release().then(() => {
            return this.adquire();
          })
        })
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
          this.$window.addEventListener("beforeunload",(event : Event) => {
            this.release();
          });
          deffered.resolve(this.webSocket)
        }
        this.webSocket.onerror = (event : Event) => {
          console.warn("Web socket error, event is: ", event)
          deffered.reject(event)
        }
        this.webSocket.onclose = (event: CloseEvent) => {
          console.log(closeDescriptionByCode[event.code])
          this.webSocket = undefined
        }

      } catch(error) {
        deffered.reject(error)
      }
      return deffered.promise
    }

    private release() {
      return this.$http.get("/release-ws").then(( ) => {
        this.webSocket = undefined
      })
    }

  }

  escobita.service('WebSocketsService', ['$http', '$q', '$window', Service]);
}
/*var wss:any;
escobita.run(['WebSocketsService', (_wss: any) => {
  wss = _wss;
}])
*/