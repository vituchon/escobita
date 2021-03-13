/// <reference path='../app.ts' />
/// <reference path='../services/_services.d.ts' />

namespace Cards { // Yup... inside the same Cards namespace as they don't collide and basically could live in the same "file" perhaps..

  export const changeDisplayModeEventName = 'change-display-mode';

  interface Scope extends ng.IScope {
    translateSuit: Function;
    displayMode: string;
    generateImageTag: Function;
    $sce: ng.ISCEService;
  }

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
        <div class="ng-card-container card-text" ng-if="displayMode === 'text'">
          <span>Palo:&nbsp;</span>
          <span class="value">{{::translateSuit(card.suit)}}&nbsp;</span>
          <span>Número:&nbsp;</span>
          <span class="value">{{::card.rank}}</div>
        </div>
        <div class="ng-card-container card-image" ng-if="displayMode === 'sprite'">
          <div ng-bind-html="$sce.trustAsHtml(generateImageTag(card.suit,card.rank))"></div>
        </div>`,
      link: function ($scope: Scope, $element: JQuery, attrs: ng.IAttributes, ngModel: ng.INgModelController) {
        $scope.translateSuit = Cards.Suits.translate
        $scope.displayMode = $scope.displayMode || 'sprite';
        $scope.generateImageTag = generateImageTag
        $scope.$sce = $sce;

        $scope.$on(changeDisplayModeEventName, (event: ng.IAngularEvent, displayMode: string) => {
          $scope.displayMode = displayMode
        });
      }
    }

  }]);

};