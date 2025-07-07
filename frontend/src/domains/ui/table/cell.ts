import { BaseDomain, Handler } from "@/domains/base";
import { TableCore } from "./table";

enum Events {
  Change,
  Select,
  Update,
}
type TheTypesOfBaseEvents = {
  [Events.Select]: TableCellCore;
  [Events.Change]: TableCellCoreState;
  [Events.Update]: { x: number; y: number; value: string };
};
type CellValueType = "index" | "text";
type TableCellCoreProps = {
  value: string;
  x: number;
  y: number;
  width?: number;
  type?: CellValueType;
  disabled?: boolean;
  $table: TableCore;
  onSelect?: (ins: TableCellCore) => void;
  onUpdate?: (opt: { x: number; y: number; value: string }) => void;
};
type TableCellCoreState = {
  type: CellValueType;
  value: string;
  selected: boolean;
  editing: boolean;
  disabled: boolean;
  width?: number;
  waitUpdate: boolean;
};

export class TableCellCore extends BaseDomain<TheTypesOfBaseEvents> {
  type: CellValueType = "text";
  value: string;
  x = 0;
  y = 0;

  width = 0;
  disabled = false;
  selected = false;
  editing = false;
  tmpValue = "";
  inputValue = "";
  waitUpdate = false;

  $table: TableCore;

  get state(): TableCellCoreState {
    return {
      type: this.type,
      value: this.value,
      selected: this.selected,
      editing: this.editing,
      disabled: this.disabled,
      width: this.width,
      waitUpdate: this.waitUpdate,
    };
  }

  constructor(props: TableCellCoreProps) {
    super(props);

    const { value, type = "text", x, y, width = 200, disabled = false, $table, onSelect, onUpdate } = props;
    this.type = type;
    this.value = value;
    this.inputValue = value;
    this.disabled = disabled;
    this.x = x;
    this.y = y;
    this.width = width;

    this.$table = $table;
    this.$table.appendCell(this);

    if (onSelect) {
      this.onSelect(onSelect);
    }
    if (onUpdate) {
      this.onUpdate(onUpdate);
    }
  }

  select() {
    if (this.disabled) {
      return;
    }
    if (this.selected === true) {
      return;
    }
    this.selected = true;
    this.emit(Events.Select, this);
    this.emit(Events.Change, { ...this.state });
  }
  unselect(options: Partial<{ silence: boolean }> = {}) {
    if (this.selected === false) {
      return;
    }
    this.selected = false;
    if (!options.silence) {
      this.emit(Events.Change, { ...this.state });
    }
  }
  edit(options: Partial<{ silence: boolean }> = {}) {
    if (this.disabled) {
      return;
    }
    if (this.editing === true) {
      return;
    }
    this.editing = true;
    if (!options.silence) {
      this.emit(Events.Change, { ...this.state });
    }
  }
  unedit(options: Partial<{ silence: boolean }> = {}) {
    if (this.editing === false) {
      return;
    }
    this.editing = false;
    if (!options.silence) {
      this.emit(Events.Change, { ...this.state });
    }
  }
  updateValue() {
    if (this.value === this.inputValue) {
      return;
    }
    this.tmpValue = this.value;
    this.value = this.inputValue;
    this.waitUpdate = true;
    this.unedit({ silence: true });
    this.unselect({ silence: true });
    this.emit(Events.Update, { x: this.x, y: this.y, value: this.inputValue });
    this.emit(Events.Change, { ...this.state });
  }
  completeUpdate() {
    this.waitUpdate = false;
    this.tmpValue = this.value;
    this.unedit({ silence: true });
    this.unselect({ silence: true });
    this.emit(Events.Change, { ...this.state });
  }
  cancelUpdate() {
    console.log("[COMPONENT]table/index - cancelUpdate", this.tmpValue, this.value, this.inputValue);
    this.value = this.tmpValue;
    this.inputValue = this.value;
    this.waitUpdate = false;
    this.unedit({ silence: true });
    this.unselect({ silence: true });
    this.emit(Events.Change, { ...this.state });
  }
  setWidth(width: number, options: Partial<{ silence: boolean }> = {}) {
    this.width = width;
    console.log("set width", width);
    if (!options.silence) {
      this.emit(Events.Change, { ...this.state });
    }
  }

  onChange(handler: Handler<TheTypesOfBaseEvents[Events.Change]>) {
    return this.on(Events.Change, handler);
  }
  onSelect(handler: Handler<TheTypesOfBaseEvents[Events.Select]>) {
    return this.on(Events.Select, handler);
  }
  onUpdate(handler: Handler<TheTypesOfBaseEvents[Events.Update]>) {
    return this.on(Events.Update, handler);
  }
}
