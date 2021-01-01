/// <reference path='../app.ts' />
/// <reference path='../services/_services.d.ts' />

namespace Cards {

  function generateImageTag (suit: number, rank: number) {
    return Cards.Sprites.generateImageTag(suit,rank)
  }
  /** Angular directive 'card': Render a card.
  *
  * Example usage (in HTML template):
  *  <card ng-model="identifier_1"/>
  *
  * @ng-model: Reference to an object of type "Api.Card"
  */
  escobita.directive('card', ['$sce', ($sce: ng.ISCEService) => {
    return {
      restrict: 'E',
      require: 'ngModel',
      scope: {
        card: "=ngModel",
        displayMode: "=?"
      },
      template: `
        <div class="ng-card-container" ng-if="displayMode === 'text'">
          <span>Palo:&nbsp;</span>
          <span class="value">{{::translateSuit(card.suit)}}&nbsp;</span>
          <span>Número:&nbsp;</span>
          <span class="value">{{::card.rank}}</div>
        </div>
        <div class="ng-card-container" ng-if="displayMode === 'sprite'">
          <div ng-bind-html="$sce.trustAsHtml(generateImageTag(card.suit,card.rank))"></div>
        </div>`,
      link: function ($scope: ng.IScope, $element: JQuery, attrs: ng.IAttributes, ngModel: ng.INgModelController) {
        $scope.translateSuit = Cards.Suits.translate
        $scope.displayMode = $scope.displayMode || 'sprite';
        $scope.generateImageTag = generateImageTag
        $scope.$sce = $sce;
      }
    }

  }]);

};