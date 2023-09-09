/// <reference path='../app.ts' />
/// <reference path='../services/_services.d.ts' />

module About {
  class Controller {
    constructor() {

    }

    public $onInit() {
      this.playFadeIn()
    }

    public playFadeIn() {
      // dev notes: using HOT tip for running an animation again "n" times, see https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_animations/Tips#run_an_animation_again
      document.getElementById("credits").className = "not-visible";
      requestAnimationFrame((time) => {
        requestAnimationFrame((time) => {
          document.getElementById("credits").className = "visible";
        });
      });
    }
  }

  escobita.controller('AboutController', [Controller]);
}