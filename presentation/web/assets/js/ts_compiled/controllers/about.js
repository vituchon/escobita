/// <reference path='../app.ts' />
/// <reference path='../services/_services.d.ts' />
var About;
(function (About) {
    var Controller = (function () {
        function Controller($state) {
            this.$state = $state;
        }
        return Controller;
    }());
    escobita.controller('AboutController', ['$state', Controller]);
})(About || (About = {}));
