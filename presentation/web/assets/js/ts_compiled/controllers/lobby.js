/// <reference path='../app.ts' />
var Lobby;
(function (Lobby) {
    var Controller = (function () {
        function Controller($state) {
            this.$state = $state;
        }
        return Controller;
    }());
    escobita.controller('LobbyController', ['$state', Controller]);
})(Lobby || (Lobby = {}));
