/// <reference path='./third_party_definitions/_definitions.ts' />

namespace Util {

  export function isNumeric(s: any): boolean {
    return !isNaN(parseFloat(s)) && isFinite(s); // based from: http://stackoverflow.com/a/6449623
  }

  export function unixToReadableClock(unix: number): string {
    return formatUnixTimestamp(unix,"HH:mm")
  }

  export function unixToReadableDay(unix: number): string {
    return formatUnixTimestamp(unix,"DD/MM/YYYY")
  }

  export function unixToReadableDate(unix: number): string {
    return formatUnixTimestamp(unix,"DD/MM/YYYY HH:mm")
  }

  export function unixToReadableDateVerbose(unix: number): string {
    return formatUnixTimestamp(unix,"dddd DD/MM/YYYY [a las] HH:mm")
  }

  export function formatUnixTimestamp(unix: number, layout:string) {
    if (isNumeric(unix)) {
      return moment.unix(unix).format(layout);
    } else {
      console.warn(`Unix timestamp value(=${unix}) is not a number`)
      return '';
    }
  }

  export interface EntityById<T> extends _.Dictionary<T>  {
    [id: number] : T
  }

  /** An identifiable entity has an numeric id that identifies unequivocally within a context. */
  export interface Identificable {
    id?: number; // it is optional due to have api model objects with nullable id, as they first are created in the client and then saved on the server granting an id, basically for avoiding this -> `Property 'id' is optional in type 'Player' but required in type 'Identificable'`
  }

  /** Generates a map by id of the given collection of identificables. */
  export function toMapById<T extends Identificable>(entites: T[]) : EntityById<T> {
    return _.indexBy(entites, 'id');
  }

  /** Generates a map by id of the given collection of elements whose id is extracted using the correspondant method */
  export function toMapByIdUsingGetter<T>(list: T[], idGetterFunc: (elem:T) => number) : EntityById<T> {
    return list.reduce((map: any, elem: T) => {
      map[idGetterFunc(elem)] = elem
      return map;
    }, {});
  }
}