/// <reference path='../app.ts' />
/// <reference path='../services/_services.d.ts' />

namespace CountdownClock {

  /** A countdown clock widget.
  *
  * Usage: <countdown-clock></countdown-clock>
  *
  * Parameters:
  * @seconds: (Optional) Total seconds to count, defaults to 60.
  * @onEnd: (Optional) Function to be invoked when user changes the selection. The actual function must include '(value)' literally ir order to receive the new selected value as input.
      So if the  function is identified in the current scope by "myOnChangeFunc" then the right way to reference the function is "myOnChangeFunc(value)".
  */
  escobita.directive('countdownClock', () => {
    return <ng.IDirective>{
      restrict: 'E',
      scope: {
        onEnd: "&?",
        totalSeconds: "=?seconds",
      },
      bindToController: true,
      controller: 'CountdownClockCtr',
      controllerAs: 'ctr',
      template: `
        <button ng-click="ctr.start()">Empezar contador</button>
        <button ng-click="ctr.cancel()">Detener contador</button>
        <div class="countdownContainer horizontalLayout" ng-style="{'background-image': ctr.backgroundRotation()}" ng-if="ctr.countdownInterval != null">
          <div class="countdownCover horizontalLayout">
            <span>{{ctr.minutesStr}}:{{ctr.secondsStr}}</span>
          </div>
        </div>`
    }
  });

  type WrappedOnEndFunc = (params: {seconds: number}) => void;

  interface Contdown {
    total?: number;
    sec?: number;
    min?: number;
  }

  class Controller {

    public onEnd: WrappedOnEndFunc;

    public totalSeconds: number;
    public elapsedSeconds: number;

    public countdownInterval: ng.IPromise<any> = null;

    public minutesStr: string;
    public secondsStr: string;

    constructor(private $scope: ng.IScope, private $interval: ng.IIntervalService) {
    }

    public $onInit() {
      this.totalSeconds = this.totalSeconds || 12;
      this.onEnd = this.onEnd || ((params: {seconds: number}) => {})
      this.$scope.$on('$destroy', () => {
        this.cancel();
      });
    }

    public start () {
      this.elapsedSeconds = 0
      this.count()
      this.countdownInterval = this.$interval(() => {
        this.count();
      }, 1000);
    }

    public cancel () {
      if (this.countdownInterval != null) {
        this.$interval.cancel(this.countdownInterval);
        this.countdownInterval = null
      }
    }

    public end () {
      // Dev notes: It apperas that the countdown finish a little after the server already cleans up the seized resource
      Toastr.info('FIN DE TIEMPO');
      this.onEnd({seconds: this.totalSeconds});
    }

    public count() {
      this.elapsedSeconds += 1
      if (this.elapsedSeconds == this.totalSeconds) {
        this.end();
        this.cancel()
      } else {
        const remaingSeconds = this.totalSeconds - this.elapsedSeconds

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
      //let color = '#ff9800';
      let color = '#'+hexRed+hexGreen+'00';
      console.log(color)
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

  escobita.controller("CountdownClockCtr", ['$scope', '$interval', Controller]);
}
