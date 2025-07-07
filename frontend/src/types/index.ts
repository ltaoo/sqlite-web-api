import { BizError } from "@/domains/error";
import { Result, UnpackedResult } from "@/domains/result";

export type Unpacked<T> = T extends (infer U)[]
  ? U
  : T extends (...args: any[]) => infer U
  ? U
  : T extends Promise<infer U>
  ? U
  : T extends Result<infer U>
  ? U
  : T;

export type MutableRecord<U> = {
  [SubType in keyof U]: {
    type: SubType;
    data: U[SubType];
  };
}[keyof U];
export type MutableRecord2<U> = {
  [SubType in keyof U]: {
    type: SubType;
  } & U[SubType];
}[keyof U];

export type BaseApiResp<T> = {
  code: number;
  msg: string;
  data: T;
};

export type ListResponse<T> = {
  total: number;
  page: number;
  page_size: number;
  no_more: boolean;
  list: T[];
};
export type ListResponseWithCursor<T> = {
  page_size: number;
  next_marker?: string;
  total: number;
  list: T[];
};

export type RequestedResource<T extends (...args: any[]) => any> = UnpackedResult<Unpacked<ReturnType<T>>>;
export type Shift<T extends any[]> = ((...args: T) => void) extends (arg1: any, ...rest: infer R) => void ? R : never;

export type Rect = {
  width: number;
  height: number;
  x: number;
  y: number;
  // scrollHeight: number;
  left: number;
  right: number;
  top: number;
  bottom: number;
};

export interface JSONArray extends Array<JSONValue> {}
export type JSONValue = string | number | boolean | JSONObject | JSONArray | null;
export type JSONObject = { [Key in string]?: JSONValue };
