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
        <div class="ng-card-container">
          <span>Palo:&nbsp;</span>
          <span class="value">{{::translateSuit(card.suit)}}&nbsp;</span>
          <span>Número:&nbsp;</span>
          <span class="value">{{::card.rank}}</div>
        </div>`,
      link: function ($scope: ng.IScope, $element: JQuery, attrs: ng.IAttributes, ngModel: ng.INgModelController) {
        $scope.translateSuit = Cards.Suits.translate
      }
    }

  });

};