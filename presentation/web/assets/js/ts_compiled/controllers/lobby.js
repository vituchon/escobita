/// <reference path='../app.ts' />
/// <reference path='../services/services.ts' />
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
            this.gamesService.getGames().then(function (games) {
                _this.games = games;
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
        Controller.prototype.doesGameAcceptPlayers = function (game) {
            return !Games.isStarted(game);
        };
        Controller.prototype.joinGame = function (game, player) {
            var _this = this;
            Games.addPlayer(game, player);
            this.gamesService.updateGame(game).then(function () {
                _this.$state.go("game", { game: game }, { relative: false });
            });
        };
        return Controller;
    }());
    escobita.controller('LobbyController', ['$state', 'GamesService', 'PlayersService', Controller]);
})(Lobby || (Lobby = {}));
