/// <reference path='../app.ts' />
/// <reference path='../services/services.ts' />
var Lobby;
(function (Lobby) {
    var Controller = (function () {
        function Controller($state, gamesService) {
            this.$state = $state;
            this.gamesService = gamesService;
            this.games = [];
        }
        Controller.prototype.createGame = function (game) {
            var _this = this;
            this.gamesService.createGame(game).then(function (createdGame) {
                _this.games.push(createdGame);
            });
        };
        Controller.prototype.updateGameList = function () {
            var _this = this;
            this.gamesService.getGames().then(function (games) {
                _this.games = games;
            });
        };
        return Controller;
    }());
    escobita.controller('LobbyController', ['$state', 'GamesService', Controller]);
})(Lobby || (Lobby = {}));
