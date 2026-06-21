export type ID = string | number

export interface Item<T extends ID> {
  id: T
  name: string
}

export type FormValue = Record<string, unknown>
export type FormShape = Record<
  string,
  FieldString | FieldNumber | FieldItems<any>
>

export interface Field {
  readonly type: string
  name: string
  info?: string
}

export interface FieldString extends Field {
  readonly type: "string"
}

export interface FieldNumber extends Field {
  readonly type: "number"
  min?: number
  max?: number
}

export interface FieldItems<T extends ID> extends Field {
  readonly type: "items"
  items: Item<T>[]
}
