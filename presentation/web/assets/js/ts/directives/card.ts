/// <reference path='../app.ts' />
/// <reference path='../api-model.ts' />

namespace Cards {

  /** Angular directive 'card': Render a card.
  *
  * Example usage (in HTML template):
  *  <card ng-model="identifier_1"/>
  *
  * @ng-model: Reference to an object of type "Api.Card"
  */
  escobita.directive('card', () => {
    return {
      restrict: 'E',
      require: 'ngModel',
      scope: {
        card: "=ngModel",
      },
      template: `
        <div class="card-container">
          <div class="lightText" style="font-size:12px">Palo:&nbsp;</div>
          <div>{{::translateSuit(card.suit)}}&nbsp;</div>
          <div class="lightText" style="font-size:12px">Número:&nbsp;</div>
          <div>{{::card.rank}}</div>
        </div>`,
      link: function ($scope: ng.IScope, $element: JQuery, attrs: ng.IAttributes, ngModel: ng.INgModelController) {
        $scope.translateSuit = Cards.Suits.translate
      }
    }

  });

};