/// <reference path='./third_party_definitions/_definitions.ts' />

namespace Util {

  export function isDefined(value:any):boolean {
    return (!_.isNull(value) && !_.isUndefined(value) );
  }

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

namespace Arrays {

    /** Determines if two arrays has the same values regardless their order and duplicates.
        Dev notes: It doesn't support arrays with duplicated values with same quantity of elements, e.g: [1,2,2] would be considered as equal to [1,1,2]
    */
    export function hasSameValues<T>(a1: T[],a2: T[], areEquals: (e1:T,e2:T) => boolean): boolean {
      if (_.size(a1) !== _.size(a2)) {
        return false
      }
      const a1ConstainsA2 = _.reduce(a1,(acc,e1) => {
        return acc && ((a2).findIndex((e2) => areEquals(e1,e2)) !== -1)
      },true)
      const a2ConstainsA1 = _.reduce(a2,(acc,e2) => {
        return acc && ((a1).findIndex((e1) => areEquals(e1,e2)) !== -1)
      },true)

      return a1ConstainsA2 && a2ConstainsA1
    }

  /**
   * Generates an array containing all possible combinations (in form of array).
   * So it returns an array of n! elements where each element is an array holding a possible combination.
   * @param array
   * @returns An array of arrays that containst all the possible combinations.
   */
  export function generatePermutations<T>(array: T[]): T[][] {
    if (array.length == 1) {
      return [array]
    } else {
      var combinations: T[][] = []
      array.forEach( (value,index) => {
        const others = array.slice(0,index).concat(array.slice(index+1,array.length))
        const subcombinations = generatePermutations(others)
        subcombinations.forEach((subcombination) => {
          //const combination = [value].concat(subcombination).flat(1)
          const combination = [value].concat(...subcombination)
          combinations.push(combination)
        })
      });
      return combinations
    }
  }

  export namespace Tests {

    interface CombineTestRun<T> {
      title: string,
      input: T[],
      expected: T[][];
    }

    export function CombineWorks() {
      const testRuns = [
        <CombineTestRun<number>>{
          title: "Empty array",
          input: [],
          expected: [],
        },
        <CombineTestRun<number>>{
          title: "2 element array",
          input: [1,2],
          expected: [[1,2],[2,1]],
        },
        <CombineTestRun<number>>{
          title: "3 element array",
          input: [1,2,3],
          expected: [[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]],
        },
        <CombineTestRun<number>>{
          title: "4 element array",
          input: [1,2,3,4],
          expected: [[1,2,3,4],[1,2,4,3],[1,3,2,4],[1,3,4,2],[1,4,2,3],[1,4,3,2],[2,1,3,4],[2,1,4,3],[2,3,1,4],[2,3,4,1],[2,4,1,3],[2,4,3,1],[3,1,2,4],[3,1,4,2],[3,2,1,4],[3,2,4,1],[3,4,1,2],[3,4,2,1],[4,1,2,3],[4,1,3,2],[4,2,1,3],[4,2,3,1],[4,3,1,2],[4,3,2,1]],
        },
      ]

      _.forEach(testRuns,(testRun) => {
        const computed = Arrays.generatePermutations(testRun.input);
        if (_.size(computed) != _.size(testRun.expected)) {
          console.error("No tiene la misma dimension, computada es: ",_.size(computed), " y esperada es: ", _.size(testRun.expected))
        }
        for (var i = 0; i < _.size(computed); i++) {
          const strComputedElem = JSON.stringify(computed[i])
          const strExpectedElem = JSON.stringify(testRun.expected[i])
          if (strComputedElem !== strExpectedElem) {
            console.error(computed[i], " <> ", testRun.expected[i]);
          }
        }
      })
    }
  }
}
