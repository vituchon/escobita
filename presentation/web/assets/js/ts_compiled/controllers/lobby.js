/// <reference path='../app.ts' />
/// <reference path='../services/_services.d.ts' />
var Lobby;
(function (Lobby) {
    var Controller = /** @class */ (function () {
        function Controller($rootElement, $scope, $timeout, $state, $q, gamesService, playersService) {
            var _this = this;
            this.$rootElement = $rootElement;
            this.$scope = $scope;
            this.$state = $state;
            this.$q = $q;
            this.gamesService = gamesService;
            this.playersService = playersService;
            this.playerName = ""; // for entering a player name
            this.loading = false;
            this.showCards = false;
            this.games = [];
            this.loading = true;
            var getClientPlayerPromise = this.playersService.getClientPlayer().then(function (player) {
                _this.player = player;
                _this.playerName = _this.player.name;
            });
            var getGamesPromise = this.gamesService.getGames().then(function (games) {
                _this.games = games;
            });
            this.$q.all([getClientPlayerPromise, getGamesPromise])["finally"](function () {
                _this.loading = false;
            });
            $rootElement.bind("keydown keypress", function (event) {
                if (event.which === 13) {
                    $timeout(function () {
                        if (!_.isUndefined(_this.playerName)) {
                            _this.updatePlayerName(_this.playerName);
                        }
                    });
                    event.preventDefault();
                }
            });
            $scope.$on('$destroy', function () {
                $rootElement.unbind("keydown keypress");
            });
        }
        Controller.prototype.createGame = function (game) {
            var _this = this;
            this.loading = true;
            this.gamesService.createGame(game).then(function (createdGame) {
                _this.games.push(createdGame);
            })["finally"](function () {
                _this.loading = false;
            });
        };
        Controller.prototype.updateGameList = function () {
            var _this = this;
            this.loading = true;
            this.gamesService.getGames().then(function (games) {
                _this.games = games;
            })["finally"](function () {
                _this.loading = false;
            });
        };
        Controller.prototype.updatePlayerName = function (name) {
            var _this = this;
            this.loading = true;
            this.player.name = name;
            this.playersService.updatePlayer(this.player).then(function (player) {
                _this.player = player;
            }).then(function () {
                _this.showCards = true;
            })["finally"](function () {
                _this.loading = false;
            });
        };
        Controller.prototype.updatePlayersList = function () {
            var _this = this;
            this.playersService.getPlayers().then(function (players) {
                _this.players = players;
            });
        };
        Controller.prototype.doesGameAcceptPlayers = function (game) {
            return !Games.hasMatchInProgress(game);
        };
        Controller.prototype.joinGame = function (game, player) {
            var _this = this;
            this.loading = true;
            this.gamesService.getGameById(game.id).then(function (game) {
                Games.addPlayer(game, player);
                _this.gamesService.updateGame(game).then(function () {
                    _this.$state.go("game", {
                        game: game,
                        player: player
                    }, { relative: false });
                });
            })["finally"](function () {
                _this.loading = false;
            });
        };
        return Controller;
    }());
    escobita.controller('LobbyController', ['$rootElement', '$scope', '$timeout', '$state', '$q', 'GamesService', 'PlayersService', Controller]);
})(Lobby || (Lobby = {}));
