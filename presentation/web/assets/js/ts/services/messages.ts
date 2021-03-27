
namespace Messages {

  export interface Message extends Api.Message {
  }

  export function newMessage(gameId: number, playerId: number, text: string): Message {
    return {
      gameId: gameId,
      playerId: playerId,
      text: text,
    }

  }

  export class Service {
    constructor(private $http: ng.IHttpService, private $q: ng.IQService) {
    }

    getMessages(): ng.IPromise<Message[]> {
      return this.$http.get<Message[]>(`/api/v1/messages`).then((response) => {
        return response.data;
      });
    }

    getMessageById(messagesId: number): ng.IPromise<Message> {
      return this.$http.get<Message>(`/api/v1/messages/${messagesId}`).then((response) => {
        return response.data;
      })
      // TODO: below is not used code, server will return a 500 is a client seek for a non existing message, guess a 404 is more correct...
      /*.catch((err) => {
        if (err.status === 404) { // not found: there is no message with the given id
          return null
        } else {
          return err
        }
      })*/
    }

    getMessagesByGame(gameId: number, since?: number): ng.IPromise<Message[]> {
      const parameters = {
        'since': since,
      }
      const config:ng.IRequestShortcutConfig = {
        params: parameters
      }
      return this.$http.get<Message[]>(`/api/v1/messages/get-by-game/${gameId}`, config).then((response) => {
        return response.data;
      })
    }


    getClientMessage(): ng.IPromise<Message> {
      return this.$http.get<Message>(`/api/v1/message`).then((response) => {
        return response.data;
      })
    }

    createMessage(message: Message): ng.IPromise<Message> {
      return this.$http.post<Message>(`/api/v1/messages`,message).then((response) => {
        return response.data
      })
    }

    updateMessage(message: Message): ng.IPromise<Message> {
      return this.$http.put<Message>(`/api/v1/messages/${message.id}`,message).then((response) => {
        return response.data
      })
    }

    deleteMessage(message: Message): ng.IPromise<any> {
      return this.$http.delete<any>(`/api/v1/messages/${message.id}`).then((_) => {
        return null;
      })
    }

  }

  escobita.service('MessagesService', ['$http', '$q', Service]);
}