/**
 * @file 数据库管理
 */
import { createSignal, For, onMount } from "solid-js";
import { Plus, X } from "lucide-solid";

import { ViewComponent } from "@/store/types";
import { RequestCore } from "@/domains/request";
import { DialogCore, PopoverCore } from "@/domains/ui";
import { Dialog } from "@/components/ui";
import { TableWithColumns } from "@/domains/ui/table/table";
import { fetchTableList, fetchTableListProcess } from "@/biz/services";
import { TablePanel, TablePanelCore } from "./table";

export const SqliteDatabasePage: ViewComponent = (props) => {
  const { app, storage } = props;
  // const $request = $page.$request;
  const $request = {
    tableList: new RequestCore(fetchTableList, {
      process: fetchTableListProcess,
      delay: null,
      // onSuccess(v) {
      //   $filter.setTables(v);
      // },
    }),
  };
  const $settings = new PopoverCore({});
  const $tables = new DialogCore({
    title: "打开表",
    footer: false,
  });
  const $record = new PopoverCore({});

  const [tables, setTables] = createSignal($request.tableList.response);
  const [cur, setPanel] = createSignal<ReturnType<typeof TablePanelCore> | null>(null);
  // const [filters, setFilters] = createSignal($page.ui.$filter.values);
  // const [record, setRecord] = createSignal<{ field: string; value: string }[]>([]);
  const [panels, setPanels] = createSignal<ReturnType<typeof TablePanelCore>[]>([]);

  function handleClickTable(table: TableWithColumns) {
    const existing = panels().find((p) => p.name === table.name);
    if (existing) {
      setPanel(existing);
      $tables.hide();
      return;
    }
    const $panel = TablePanelCore({ ...props, table, tables: $request.tableList.response || [] });
    $panel.onMounted(() => {
      setPanel($panel);
      $panel.loadDataSource();
    });
    setPanels((prev) => {
      return [...prev, $panel];
    });
    $tables.hide();
  }
  function handleClickTableTab(panel: ReturnType<typeof TablePanelCore>) {
    setPanel(panel);
  }
  function handleRemovePanel(panel: ReturnType<typeof TablePanelCore>) {
    const panelsCopy = [...panels()];
    const existingIndex = panelsCopy.findIndex((p) => p.name === panel.name);
    if (existingIndex == -1) {
      return;
    }
    setPanels((prev) => {
      const nextPanels = prev.filter((p) => p.name !== panel.name);
      return nextPanels;
    });
    if (cur()?.name !== panel.name) {
      return;
    }
    const nextPanel = (() => {
      const prev = panelsCopy[existingIndex - 1];
      if (prev) {
        return prev;
      }
      const next = panelsCopy[existingIndex + 1];
      if (next) {
        return next;
      }
      return null;
    })();
    if (nextPanel === null) {
      $tables.show();
      return;
    }
    setPanel(nextPanel);
  }

  $request.tableList.onResponseChange((v) => {
    if (v === null) {
      return;
    }
    // console.log("[PAGE]home/index - the table list is", v);
    // const prev = storage.get("column_widths", {});
    // const widths = prev[$page.ui.$table.name];
    // if (widths) {
    //   for (let i = 0; i < v.length; i += 1) {
    //     const t = v[i];
    //     for (let j = 0; j < t.columns.length; j += 1) {
    //       const column = t.columns[j];
    //       const w = widths[column.name] || 200;
    //       column.width = w;
    //     }
    //   }
    // }
    setTables(v);
    if (panels().length === 0) {
      $tables.show();
    }
  });

  onMount(() => {
    const keycodes: Record<string, boolean> = {};
    document.addEventListener("keydown", (event) => {
      // console.log(event.code);
      keycodes[event.code] = true;
      (() => {
        if (event.code === "KeyC" && (keycodes.ControlLeft || keycodes.ControlRight)) {
          const panel = cur();
          if (panel === null) {
            return;
          }
          const curCell = panel.ui.$table.cells.find((cell) => cell.selected);
          if (curCell) {
            app.copy(curCell.value);
            app.tip({
              text: ["已复制到剪贴板"],
            });
          }
        }
      })();
      for (let i = 0; i < panels().length; i += 1) {
        const panel = panels()[i];
        panel.keycodes[event.code] = true;
      }
    });
    document.addEventListener("keyup", (event) => {
      keycodes[event.code] = false;
      for (let i = 0; i < panels().length; i += 1) {
        const panel = panels()[i];
        panel.keycodes[event.code] = false;
      }
    });
    $request.tableList.run();
  });

  return (
    <div class="h-screen flex flex-col">
      <div class="flex items-center h-[48px] border border-t-0 border-l-0 border-r-0 bg-w-bg-0">
        <div class="placeholder px-4 w-[32px] h-full border border-t-0 border-l-0 border-b-0 border-r-1"></div>
        <For each={panels()}>
          {(panel, i) => {
            return (
              <div
                classList={{
                  "flex items-center px-4 h-full border border-t-0 border-l-0 border-b-0 border-r-1 cursor-pointer":
                    true,
                  "bg-w-bg-0": cur()?.name === panel.name,
                }}
                onClick={() => {
                  handleClickTableTab(panel);
                }}
              >
                <div>{panel.name}</div>
                <div
                  class="ml-4"
                  onClick={(event) => {
                    event.stopPropagation();
                    handleRemovePanel(panel);
                  }}
                >
                  <X class="w-4 h-4 cursor-pointer" />
                </div>
              </div>
            );
          }}
        </For>
        <div
          class="placeholder flex items-center justify-center w-[48px] h-full border border-t-0 border-l-0 border-b-0 border-r-1  cursor-pointer"
          onClick={() => {
            $tables.show();
          }}
        >
          <div>
            <Plus class="w-4 h-4" />
          </div>
        </div>
      </div>
      <div class="panels relative flex-1">
        <For each={panels()}>
          {(panel, i) => {
            return (
              <div classList={{ "absolute inset-0 bg-w-bg-0": true, hidden: cur()?.name !== panel.name }}>
                <div></div>
                <TablePanel store={panel} />
              </div>
            );
          }}
        </For>
      </div>
      {/* <Popover
        store={$settings}
        content={
          <div class="space-y-1">
            <div class="">
              <For each={table().columns}>
                {(column) => {
                  return (
                    <div class="flex items-center justify-between w-full">
                      <div>{column.name}</div>
                      <input
                        class="ml-4"
                        placeholder="输入宽度"
                        onBlur={(event) => {
                          const { value } = event.currentTarget;
                          $page.setWidth(column, value);
                        }}
                      />
                    </div>
                  );
                }}
              </For>
            </div>
          </div>
        }
      ></Popover> */}
      {/* <Popover
        store={$record}
        content={
          <div class="space-y-2">
            <For each={record()}>
              {(field) => {
                return (
                  <div class="flex items-center">
                    <div class="w-[88px]">{field.field}:</div>
                    <div>{field.value}</div>
                  </div>
                );
              }}
            </For>
          </div>
        }
      ></Popover> */}
      <Dialog store={$tables}>
        <div class="w-[520px]">
          <div class="mt-4 overflow-y-auto max-h-[680px]">
            <div class="space-y-2 ">
              <For each={tables()}>
                {(t) => {
                  return (
                    <div
                      classList={{
                        "flex items-center justify-between px-4 py-2 rounded-md bg-w-bg-0 cursor-pointer truncate": true,
                        "hover:bg-w-bg-3": true,
                      }}
                      onClick={() => {
                        handleClickTable(t);
                      }}
                    >
                      <div class="overflow-hidden truncate">{t.name}</div>
                      {/* <div
                      class="ml-2"
                      onClick={(event) => {
                        event.stopPropagation();
                        const { pageX, pageY } = event;
                        $settings.show({
                          x: pageX,
                          y: pageY,
                        });
                      }}
                    >
                      <Settings class="w-4 h-4 cursor-pointer" />
                    </div> */}
                    </div>
                  );
                }}
              </For>
            </div>
          </div>
        </div>
      </Dialog>
    </div>
  );
};
