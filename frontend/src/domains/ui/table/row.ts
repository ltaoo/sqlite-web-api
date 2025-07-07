import { BaseDomain, Handler } from "@/domains/base";
import { TableCellCore } from "./cell";

enum Events {
  Change,
  Select,
}
type TheTypesOfBaseEvents = {
  [Events.Select]: TableRowCore;
  [Events.Change]: TableRowCoreState;
};
type TableCellCoreProps = {
  index: number;
  cells: TableCellCore[];
  onSelect?: (ins: TableRowCore) => void;
};
type TableRowCoreState = {
  cells: TableCellCore[];
};

export class TableRowCore extends BaseDomain<TheTypesOfBaseEvents> {
  index = 0;

  cells: TableCellCore[] = [];

  get state(): TableRowCoreState {
    return {
      cells: this.cells,
    };
  }

  constructor(props: TableCellCoreProps) {
    super(props);

    const { index, cells, onSelect } = props;
    this.index = index;
    this.cells = cells;

    if (onSelect) {
      this.onSelect(onSelect);
    }
  }

  select() {
    this.emit(Events.Change, { ...this.state });
  }
  unselect() {
    this.emit(Events.Change, { ...this.state });
  }

  onChange(handler: Handler<TheTypesOfBaseEvents[Events.Change]>) {
    return this.on(Events.Change, handler);
  }
  onSelect(handler: Handler<TheTypesOfBaseEvents[Events.Select]>) {
    return this.on(Events.Select, handler);
  }
}
