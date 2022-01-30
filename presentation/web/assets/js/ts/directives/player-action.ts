/// <reference path='../app.ts' />
/// <reference path='../api-model.ts' />
/// <reference path='../services/_services.d.ts' />

namespace Cards {

  /** Angular directive 'playerAction': Render a player action
  *
  * Example usage (in HTML template):
  *  <player-action ng-model="identifier_1"/>
  *
  * @ng-model: Reference to an object of type " Api.PlayerTakeAction | Api.PlayerDropAction"
  */

  interface Scope extends ng.IScope {
    isTakenAction: boolean;
    translateSuit: (suit: number) => string;
    action: Api.PlayerTakeAction | Api.PlayerDropAction;
  }

  escobita.directive('playerAction', () => {
    return {
      restrict: 'E',
      require: 'ngModel',
      scope: {
        action: "=ngModel",
        order: "="
      },
      template: `
        <div class="ng-player-action-container">
          <div class="title">
            <strong>{{:: order+1}}.&nbsp;</strong>
            <i>{{:: isTakenAction ? "Toma" : "Descarte"}}&nbsp;</i>
          </div>
          <div class="hand-card">
            <span>Carta de mano usada</span>
            <card ng-model="action.handCard"></card>
          </div>
          <div class="boards-card" ng-if="isTakenAction">
            <span>Cartas de mesa usadas</span>
            <div ng-repeat="card in action.boardCards">
              <card ng-model="card"></card>
            </div>
            <strong ng-show="action.isEscobita">¡Fue escobita!</strong>
          </div>
        </div>`,
      link: function ($scope: Scope, $element: JQuery, attrs: ng.IAttributes, ngModel: ng.INgModelController) {
        $scope.translateSuit = Cards.Suits.translate
        $scope.isTakenAction = Matchs.isTakeAction($scope.action)
      }
    }

  });

};