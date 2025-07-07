import { BaseDomain, Handler } from "@/domains/base";

export enum TableColumnType {
  Index = "index",
  DateTime = "datetime",
  Text = "text",
  Integer = "integer",
  Table = "table",
}
export type TableColumn = {
  name: string;
  type: TableColumnType;
  width?: number;
  references?: string;
  is_primary_key?: number;
  options?: {
    format?: string;
  };
};

enum Events {
  Change,
}
type TheTypesOfBaseEvents = {
  [Events.Change]: TableColumnCoreState;
};
type TableColumnCoreProps = {
  name: string;
  type: TableColumn["type"];
  width?: number;
  references?: string;
  is_primary_key: number;
  x: number;
  options?: {
    format?: string;
  };
};
type TableColumnCoreState = {
  width: number;
  name: string;
  type: TableColumn["type"];
};

export class TableColumnCore extends BaseDomain<TheTypesOfBaseEvents> {
  name: string;
  type: TableColumn["type"] = TableColumnType.Text;
  primary_key: TableColumn["is_primary_key"] = 0;
  references?: string;
  width = 0;
  x = 0;
  options: TableColumn["options"] = {};

  get state(): TableColumnCoreState {
    return {
      name: this.name,
      type: this.type,
      width: this.width,
    };
  }

  constructor(props: TableColumnCoreProps) {
    super(props);

    const { name, type, is_primary_key, references, x, width = 200, options = {} } = props;

    this.name = name;
    this.type = type;
    this.references = references;
    this.primary_key = is_primary_key;
    this.x = x;
    this.width = width;
    this.options = options;
  }
  setWidth(width: number) {
    this.width = width;
    this.emit(Events.Change, { ...this.state });
  }

  onChange(handler: Handler<TheTypesOfBaseEvents[Events.Change]>) {
    return this.on(Events.Change, handler);
  }
}
