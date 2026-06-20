export type ID = string | number

export interface Item<T extends ID> {
  id: T
  name: string
}
