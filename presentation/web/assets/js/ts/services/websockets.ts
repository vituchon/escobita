


namespace WebSockets {
  interface CloseEvent extends Event {
    code: number;
    reason: string;
    wasClean: boolean;
  }

  interface ServerMessage {
    kind: string;
    message: string;
  }

  const normalCloseEventCode = 1000
  const closeDescriptionByCode: {[code:number]: string} = { // taken from https://github.com/Luka967/websocket-close-codes
    [normalCloseEventCode]:	"CLOSE_NORMAL",
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

  function resolveProtocol() {
    const isSecure = window.location.protocol.indexOf("https") != -1
    return (isSecure) ? "wss" : "ws"
  }

  function resolveHost() {
    return window.location.host;
  }


  export class Service {

    private webSocket: WebSocket

    constructor(private $http: ng.IHttpService, private $q: ng.IQService, private $window: ng.IWindowService) {
    }

    public retrieve() {
      if (_.isUndefined(this.webSocket)) {
        return this.adquire()
      } else {
        if (this.webSocket.readyState === WebSocket.OPEN) {
          console.debug("Web socket already adquired and open")
          return this.$q.when(this.webSocket) // in the case is called more than once, return the ws already established
        } else {
          return this.$q.reject("already binded to another tab")
        }
      }
    }

    private adquire() {
      const deffered = this.$q.defer<WebSocket>();
      try {
        const protocol = resolveProtocol();
        const host = resolveHost();
        this.webSocket = new WebSocket(`${protocol}://${host}/adquire-ws`);

        this.webSocket.onopen = (event : Event) => {
          this.$window.addEventListener("beforeunload",(event : Event) => {
            cleanup() // try to clean up before exiting
          });
          deffered.resolve(this.webSocket)
        }
        this.webSocket.onerror = (event : Event) => {
          console.warn("Web socket error, event is: ", event)
          cleanup()
        }
        this.webSocket.onclose = (event: CloseEvent) => { // can be called directly (after the server closes conn) or after an error ocurrs that result in the invokaton after onerror callback)
          const reason = event.reason || closeDescriptionByCode[event.code]
          console.debug("Closing web socket. Was clean:", event.wasClean,", code:", event.code, ", reason:", reason)
          const hasNotCloseNormal = event.code !== normalCloseEventCode
          if (hasNotCloseNormal || !event.wasClean) {
            deffered.reject(reason) // if was already resolved then this reject has no effect
          }
          this.webSocket?.removeEventListener?.("message", handleServerMessage)
          this.webSocket = undefined
        }
        const handleServerMessage =  (event: MessageEvent<any>) => {
          const notification: ServerMessage = JSON.parse(event.data)
          if (notification.kind === "debug") {
            Toastr.info(notification.message)
          }
        }
        const cleanup = () => {
          this.webSocket?.removeEventListener?.("message", handleServerMessage)
          this.release();
        }
        this.webSocket.addEventListener("message", handleServerMessage)
      } catch(error) {
        deffered.reject(error)
        throw error
      }
      return deffered.promise
    }

    public release(code?: number, reason?: string) {
      if (_.isUndefined(this.webSocket)) {
        return this.$q.reject("No web socket to release")
      }
      this.webSocket.close(code, reason) // honoring convention (but it is not necesary)
      return this.$http.get("/release-ws").then(() => { // this will make the server to send a close message to this client websocket that will trigger the websocket's onclose method invocation
      })
    }

    public pingme(message?: string) {
      const config: ng.IRequestShortcutConfig = {
        params: {
          message: message
        }
      };
      return this.$http.get("/send-message-ws",config)
    }

    public pingall(message?: string) {
      const config: ng.IRequestShortcutConfig = {
        params: {
          message: message
        }
      };
      return this.$http.get("/send-message-all-ws",config)
    }
  }

  escobita.service('WebSocketsService', ['$http', '$q', '$window', Service]);
}


var wss:any;
escobita.run(['WebSocketsService', (_wss: any) => {
  wss = _wss;
}])
