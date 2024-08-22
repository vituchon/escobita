/// <reference path='../app.ts' />
/// <reference path='../services/_services.d.ts' />

namespace CountdownClock {

  /**
   * A countdown clock widget that displays a countdown timer. It can be customized by providing
   * the total number of seconds to count down from, the refresh rate of the display, and a callback
   * function that is invoked when the countdown reaches zero.
   *
   * @example
   * // Basic usage in HTML:
   * // <countdown-clock total-seconds="120" refresh-rate-in-milis="1000" on-end="callbackFunction()"></countdown-clock>
   *
   * @param {number} [totalSeconds=60] - Optional. The total number of seconds to count down from.
   *                                      Defaults to 30 seconds if not provided.
   * @param {number} [refreshRateInMilis=1000] - Optional. The interval in milliseconds to refresh the display.
   *                                             Defaults to 100ms (0.1 second) if not provided.
   * @param {function} [onEnd] - Optional. A callback function to be invoked when the countdown finishes.
   *                             No arguments are passed to this function.
   * @param {Handler} [handler] - Optional. The countdown clock handler used for precise control and manipulation.
   * @returns {void}
   */
  escobita.directive('countdownClock', () => {
    return <ng.IDirective>{
      restrict: 'E',
      scope: {
        onEnd: "&?",
        totalSeconds: "=?",
        refeshRateInMilis : "=?",
        handler: "=?",
      },
      bindToController: true,
      controller: 'CountdownClockCtr',
      controllerAs: 'ctr',
      template: `
        <!-- <button ng-click="ctr.start()">Empezar contador</button>
        <button ng-click="ctr.cancel()">Detener contador</button> -->
        <div class="countdownContainer horizontalLayout" ng-style="{'background-image': ctr.backgroundRotation()}" ng-if="ctr.countdownInterval != null">
          <div class="countdownCover horizontalLayout">
            <span>{{ctr.minutesStr}}:{{ctr.secondsStr}}</span>
          </div>
        </div>`
    }
  });

  type WrappedOnEndFunc = () => void;

  export interface Handler {
    start(): void;
    cancel(): void;
  }

  class Controller {

    public onEnd: WrappedOnEndFunc;

    private totalSeconds: number;
    private elapsedSeconds: number;

    public countdownInterval: ng.IPromise<any> = null;
    private refreshRateInMilis: number;

    public minutesStr: string;
    public secondsStr: string;

    public handler: Handler;

    constructor(private $scope: ng.IScope, private $interval: ng.IIntervalService, private $attrs: ng.IAttributes) {
    }

    public $onInit() {
      this.handler = this
      this.totalSeconds = this.totalSeconds || 30;
      this.refreshRateInMilis = this.refreshRateInMilis || 100
      this.onEnd = this.onEnd || ( () => {} )
      this.$scope.$on('$destroy', () => {
        this.cancel();
      });

      if (Util.isDefined(this.$attrs["startInmediate"])) {
        this.start()
      }
    }

    public start () {
      if (this.countdownInterval == null) {
        this.elapsedSeconds = 0
        this.count()
        this.countdownInterval = this.$interval(() => {
          this.count();
        }, 100);
      }
    }

    public cancel () {
      if (this.countdownInterval != null) {
        this.$interval.cancel(this.countdownInterval);
        this.countdownInterval = null
      }
    }

    public end () {
      this.onEnd();
    }

    public count() {
      const increment = this.refreshRateInMilis / 1000;
      this.elapsedSeconds += increment
      if (this.elapsedSeconds >= this.totalSeconds) {
        this.end();
        this.cancel()
      } else {
        const remaingSeconds = Math.floor(this.totalSeconds - this.elapsedSeconds)

        const minutes = Math.floor(remaingSeconds / 60)
        this.minutesStr = minutes.toString()
        if (this.minutesStr.length < 2) {
          this.minutesStr = '0' + this.minutesStr;
        }

        const seconds = remaingSeconds % 60
        this.secondsStr = seconds.toString()
        if (this.secondsStr.length < 2) {
          this.secondsStr = '0' + this.secondsStr;
        }
      }
    }

    public backgroundRotation () {
      const remaingSeconds = this.totalSeconds - this.elapsedSeconds
      const red =  Math.round(255 *  (1 - (remaingSeconds / this.totalSeconds)))
      var hexRed = red.toString(16);
      if (hexRed.length == 1) {
        hexRed = '0' + hexRed
      }
      var green = red ^ 0xFF;
      var hexGreen = green.toString(16);
      if (hexGreen.length == 1) {
        hexGreen = '0' + hexGreen
      }
      let color = '#'+hexRed+hexGreen+'00';
      //let color = '#ff9800';
      //console.log(color)
      let background = '#ebebeb';
      if (remaingSeconds > (this.totalSeconds / 2)) {
        const angle = ((remaingSeconds - (this.totalSeconds / 2)) * 180 / (this.totalSeconds / 2)) - 90;
        return 'linear-gradient(-90deg, transparent 50%, ' + color + ' 50%), linear-gradient(' + -angle + 'deg, ' + color + ' 50%, ' + background +  ' 50%)';
      } else {
        const angle = (-(remaingSeconds - (this.totalSeconds / 2)) * 180 / (this.totalSeconds / 2)) - 270;
        return 'linear-gradient(90deg, transparent 50%, ' + background +  ' 50%), linear-gradient(' + angle + 'deg, ' + color + ' 50%, ' + background +  ' 50%)'
      }
    }
  }

  escobita.controller("CountdownClockCtr", ["$scope", "$interval", "$attrs", Controller]);
}
