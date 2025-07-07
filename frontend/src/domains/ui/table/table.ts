import { BaseDomain, Handler } from "@/domains/base";

import { TableCellCore } from "./cell";
import { TableRowCore } from "./row";
import { TableColumn, TableColumnCore } from "./column";

export type TableWithColumns = {
  name: string;
  columns: TableColumn[];
};

enum Events {
  Change,
  SelectRow,
}
type TheTypesOfBaseEvents = {
  [Events.SelectRow]: TableCore;
  [Events.Change]: TableCoreState;
};
type TableCellCoreProps = {
  name?: string;
  disabled?: boolean;
  onSelect?: (ins: TableCore) => void;
};
type TableCoreState = {
  name: string;
  width: number;
  height: number;
  offsetTop: number;
  columns: TableColumnCore[];
  rows: TableRowCore[];
  visibleRows: TableRowCore[];
  pendingUpdate: TableCellCore[];
  selectedRows: number[];
};

export class TableCore extends BaseDomain<TheTypesOfBaseEvents> {
  name = "";
  disabled = false;
  selected = false;
  editing = false;

  columns: TableColumnCore[] = [];
  rows: TableRowCore[] = [];
  visibleRows: TableRowCore[] = [];
  cells: TableCellCore[] = [];
  data: string[][] = [];
  selectedRows: number[] = [];

  width = 0;
  height = 0;
  offsetTop = 0;
  size = 30;
  BUFFER_ITEM_count = 4;
  range: { start: number; end: number } = { start: 0, end: this.size };

  get state(): TableCoreState {
    return {
      name: this.name,
      width: this.width,
      height: this.height,
      offsetTop: this.offsetTop,
      columns: this.columns,
      rows: this.rows,
      visibleRows: this.visibleRows,
      selectedRows: this.selectedRows,
      pendingUpdate: this.cells.filter((c) => c.waitUpdate),
    };
  }

  constructor(props: TableCellCoreProps) {
    super(props);

    const { name = "", disabled = false, onSelect } = props;
    this.name = name;
    this.disabled = disabled;

    if (onSelect) {
      this.onSelect(onSelect);
    }
  }
  appendCell(cell: TableCellCore) {
    cell.onChange(() => {
      this.emit(Events.Change, { ...this.state });
    });
    this.cells.push(cell);
  }
  selectRow(y: number) {
    console.log("[DOMAIN]ui/table - select row", y, this.selectedRows);
    const existing = this.selectedRows.find((n) => n === y);
    if (existing !== undefined) {
      this.selectedRows = this.selectedRows.filter((n) => n !== y);
      this.emit(Events.SelectRow);
      this.emit(Events.Change, { ...this.state });
      return;
    }
    this.selectedRows = [...this.selectedRows, y];
    this.emit(Events.SelectRow);
    this.emit(Events.Change, { ...this.state });
  }
  clearSelectedRows() {
    this.selectedRows = [];
    this.emit(Events.Change, { ...this.state });
  }
  getVisibleRange(scrollTop: number) {
    if (scrollTop === 0) {
      return { start: 0, end: this.size };
    }
    const items = this.data;
    let start = 0;
    let end = 0;
    let accumulatedHeight = 0;
    (() => {
      for (let i = 0; i < items.length; i++) {
        accumulatedHeight += 42;
        // accumulatedHeight += items[i].height;
        if (accumulatedHeight >= scrollTop && start === 0) {
          start = i;
          end = start + this.size;
        }
        // if (accumulatedHeight >= scrollTop + containerHeight) {
        //     end = i;
        //     return;
        // }
      }
    })();
    const BUFFER_ITEM_count = this.BUFFER_ITEM_count;
    const startIndex = Math.max(0, start - BUFFER_ITEM_count);
    const endIndex = Math.min(end + BUFFER_ITEM_count, items.length);
    // console.log("range compute", startIndex, endIndex, this.rows.length);
    return { start: startIndex, end: endIndex };
  }
  handleScroll = (values: { scrollTop: number }) => {
    const { scrollTop } = values;
    const range = this.getVisibleRange(scrollTop);
    if (range.start === this.range.start && range.end === this.range.end) {
      return;
    }
    this.update(range);
  };
  update = (range: { start: number; end: number }) => {
    const visibleItems = this.rows.slice(range.start, range.end);
    const item = visibleItems[0];
    if (!item) {
      return;
    }
    let offset = 0;
    for (let j = 0; j < item.index; j++) {
      // const child = this.state.items[j];
      // offset += child.height;
      offset += 42;
    }
    this.offsetTop = offset;
    this.range = range;
    this.visibleRows = visibleItems;
    this.emit(Events.Change, { ...this.state });
  };
  setColumns(columns: TableColumnCore[]) {
    this.columns = columns;
    this.width = columns.reduce((total, a) => {
      // console.log("column width", a.title, a.width);
      return total + a.width;
    }, 0);
  }
  setWidth() {
    this.width = this.columns.reduce((total, a) => {
      // console.log("column width", a.title, a.width);
      return total + a.width;
    }, 0);
  }
  setData(data: string[][]) {
    this.data = data;
    this.height = (this.data.length + 2) * 42;
    // console.log("data length", this.height);
    this.offsetTop = 0;
    this.emit(Events.Change, { ...this.state });
  }
  appendData(data: string[][]) {
    this.data = [...this.data, ...data];
    this.height = (this.data.length + 2) * 42;
    // console.log("data length", this.height);
    // this.emit(Events.Change, { ...this.state });
  }
  appendRows(rows: TableRowCore[]) {
    this.rows = [...this.rows, ...rows];
  }
  setRows(rows: TableRowCore[]) {
    this.rows = rows;
  }
  setVisibleRows(rows: TableRowCore[]) {
    this.visibleRows = rows;
    this.emit(Events.Change, { ...this.state });
  }
  refresh() {
    console.log("[DOMAIN]ui/table - refresh", this.rows.length);
    this.visibleRows = this.rows.slice(this.range.start, this.range.end);
    this.emit(Events.Change, { ...this.state });
  }

  onChange(handler: Handler<TheTypesOfBaseEvents[Events.Change]>) {
    return this.on(Events.Change, handler);
  }
  onSelect(handler: Handler<TheTypesOfBaseEvents[Events.SelectRow]>) {
    return this.on(Events.SelectRow, handler);
  }
}
