import { BaseDomain, Handler } from "@/domains/base";

enum Events {
  StateChange,
  Change,
}
type TheTypesOfEvents = {
  [Events.StateChange]: DragZoneState;
  [Events.Change]: File[];
};

type DragZoneProps = {
  onChange?: (files: File[]) => void;
};
type DragZoneState = {
  hovering: boolean;
};

export class DragZoneCore extends BaseDomain<TheTypesOfEvents> {
  hovering: boolean = false;
  files: File[] = [];

  get state(): DragZoneState {
    return {
      hovering: this.hovering,
    };
  }

  constructor(props: Partial<{ _name: string }> & DragZoneProps = {}) {
    super(props);

    const { onChange } = props;
    if (onChange) {
      this.onChange(onChange);
    }
  }

  handleDragover() {
    this.hovering = true;
    this.emit(Events.StateChange, { ...this.state });
  }
  handleDragleave() {
    this.hovering = false;
    this.emit(Events.StateChange, { ...this.state });
  }
  handleDrop(files: File[]) {
    this.hovering = false;
    if (!files || files.length === 0) {
      return;
    }
    this.files = files;
    this.emit(Events.Change, [...files]);
    this.emit(Events.StateChange, { ...this.state });
  }
  getFileByName(name: string) {
    return this.files.find((f) => f.name === name);
  }

  onStateChange(handler: Handler<TheTypesOfEvents[Events.StateChange]>) {
    return this.on(Events.StateChange, handler);
  }
  onChange(handler: Handler<TheTypesOfEvents[Events.Change]>) {
    return this.on(Events.Change, handler);
  }
}
