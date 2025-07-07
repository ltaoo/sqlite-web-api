import { BaseDomain, Handler } from "@/domains/base";
import { debounce } from "@/utils/lodash/debounce";

enum Events {
  StateChange,
}
type TheTypesOfEvents = {
  [Events.StateChange]: StorageCoreState;
};

type StorageCoreProps<T> = {
  key: string;
  values: T;
  client: {
    setItem: (key: string, value: string) => void;
    getItem: (key: string) => void;
  };
};
type StorageCoreState = {};

export class StorageCore<T extends Record<string, unknown>> extends BaseDomain<TheTypesOfEvents> {
  key: string;
  values: T;
  client: StorageCoreProps<T>["client"];

  constructor(props: Partial<{ _name: string }> & StorageCoreProps<T>) {
    super(props);

    const { key, client, values } = props;
    this.key = key;
    this.values = values;
    this.client = client;
  }

  get<K extends keyof T>(key: K, defaultValue?: T[K]) {
    const v = this.values[key];
    if (v === undefined && defaultValue) {
      return defaultValue as T[K];
    }
    return v as T[K];
  }
  set = debounce(100, <K extends keyof T>(key: K, values: unknown) => {
    // console.log("cache set", key, values);
    const nextValues = {
      ...this.values,
      [key]: values,
    };
    this.values = nextValues;
    this.client.setItem(this.key, JSON.stringify(this.values));
  }) as (key: keyof T, value: unknown) => void;

  setWithPrev = debounce(100, <K extends keyof T>(key: K, handler: (v: T[K]) => T[K]) => {
    // console.log("cache set", key, values);
    const v = this.get(key);
    const nextValue = handler(v);
    const nextValues = {
      ...this.values,
      [key]: nextValue,
    };
    this.values = nextValues;
    this.client.setItem(this.key, JSON.stringify(this.values));
  });

  merge = <K extends keyof T>(key: K, values: unknown) => {
    console.log("[]merge", key, values);
    const prevValues = this.get(key) || {};
    if (typeof prevValues === "object" && typeof values === "object") {
      const nextValues = {
        ...prevValues,
        ...values,
      };
      this.set(key, nextValues);
      return nextValues;
    }
    console.warn("the params of merge must be object");
    return prevValues;
  };
  clear<K extends keyof T>(key: K) {
    const v = this.values[key];
    if (v === undefined) {
      return null;
    }
    const nextValues = {
      ...this.values,
    };
    delete nextValues[key];
    this.values = { ...nextValues };
    this.client.setItem(this.key, JSON.stringify(this.values));
  }

  onStateChange(handler: Handler<TheTypesOfEvents[Events.StateChange]>) {
    return this.on(Events.StateChange, handler);
  }
}
