/// <reference path='../ts/app.ts' />
/// <reference path='../ts/services/services.ts' />
var Lobby;
(function (Lobby) {
    var Controller = (function () {
        function Controller($state, gamesService, playersService) {
            var _this = this;
            this.$state = $state;
            this.gamesService = gamesService;
            this.playersService = playersService;
            this.games = [];
            this.playersService.getClientPlayer().then(function (player) {
                _this.player = player;
            });
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
        Controller.prototype.updatePlayer = function (player) {
            var _this = this;
            this.playersService.updatePlayer(player).then(function (player) {
                _this.player = player;
            });
        };
        Controller.prototype.updatePlayersList = function () {
            var _this = this;
            this.playersService.getPlayers().then(function (players) {
                _this.players = players;
            });
        };
        return Controller;
    }());
    escobita.controller('LobbyController', ['$state', 'GamesService', 'PlayersService', Controller]);
})(Lobby || (Lobby = {}));
