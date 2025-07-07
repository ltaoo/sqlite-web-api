import { base, BaseDomain, Handler } from "@/domains/base";
import { SelectCore, InputCore, ButtonCore } from "@/domains/ui";
import { TableColumn, TableColumnCore } from "@/domains/ui/table/column";
import { TableWithColumns } from "@/domains/ui/table/table";

import { buildQuerySQL } from "./utils";

enum Events {
  Change,
}
type TheTypesOfBaseEvents = {
  [Events.Change]: {};
};
/** multiple 待定，就是最前面的 AND、OR 这个input的类型 */
type FilterInputType = "field" | "condition" | "value" | "multiple";
export class PrefixTag {
  value: string = "";
  constructor(props: { value: string }) {
    this.value = props.value;
  }
}
export class FilterInput extends BaseDomain<TheTypesOfBaseEvents> {
  type: FilterInputType;
  //   value: string | number | null = "";
  column: TableColumn | null;

  $input: SelectCore<any> | InputCore<any> | PrefixTag;

  constructor(props: {
    type: FilterInputType;
    column?: TableColumn;
    $input: SelectCore<any> | InputCore<any> | PrefixTag;
  }) {
    super(props);

    const { type, column = null, $input } = props;
    this.type = type;
    //     this.value = value;
    this.column = column;

    this.$input = $input;
  }
  get value() {
    return this.$input.value;
  }
}

export function TableFilterCore(props: {
  table: TableWithColumns;
  tables: TableWithColumns[];
  onSearch?: (values: FilterInput[][]) => void;
  onPreviewSQL?: (sql: string) => void;
}) {
  const { table, tables } = props;

  let _table = table;
  let _tables = tables;
  let _values: FilterInput[][] = [];
  let _baseColumns: TableColumn[] = [];

  const $submit = new ButtonCore({
    onClick() {
      if (props.onSearch) {
        props.onSearch(_values);
      }
    },
  });
  const $reset = new ButtonCore({
    onClick() {
      initValues();
      emitter.emit(Events.Change, [..._values]);
      if (props.onSearch) {
        props.onSearch(_values);
      }
    },
  });
  const $more = new ButtonCore({
    onClick() {
      const rowIndex = _values.length;
      const prefix = new FilterInput({
        type: "multiple",
        $input: new SelectCore<string>({
          defaultValue: "AND",
          options: [
            {
              label: "AND",
              value: "AND",
            },
            {
              label: "OR",
              value: "OR",
            },
          ],
        }),
      });
      const first = new FilterInput({
        type: "field",
        $input: new SelectCore<string>({
          defaultValue: "",
          options: buildOptions(_baseColumns),
          onChange(v) {
            const column = _table.columns.find((c) => c.name === v);
            first.column = column || null;
            // 这个 2 表示取前两个
            _values[rowIndex] = _values[rowIndex].slice(0, 2);
            // 这个 1 表示是一行 filter input 中的第二个
            handleColumnSelectValueChange(v, [rowIndex, 1], first, _baseColumns);
          },
        }),
      });
      _values = [..._values, [prefix, first]];
      emitter.emit(Events.Change, [..._values]);
    },
  });
  // const $preview = new ButtonCore({});

  /**
   * 这个方法仅提供给 Field 类型的 Input 使用
   * 选择不同列后，出现的 条件 和 值类型 就不同
   */
  function handleColumnSelectValueChange(
    v: string | null,
    index: [number, number],
    $input: FilterInput,
    columns: TableColumn[]
  ) {
    const [x, y] = index;
    (() => {
      // 第一步，先拿到选中值对应的 column。有一种情况，reference 不是真实存在的 column，是外键关联的表名
      // 如果选择了这种，就重复 handleColumnSelectValueChange 执行
      const column = columns.find((c) => c.name === v);
      const reference = columns.find((c) => c.references === v);
      console.log("[BIZ]filter/index handleValueChange", v, index, { column, reference }, columns);
      if (reference) {
        const table = _tables.find((t) => t.name === reference.references);
        $input.column = reference;
        console.log("[BIZ]filter/index - before if (table", table, reference.references, _tables);
        if (table) {
          const $sub = new FilterInput({
            type: "field",
            column,
            $input: new SelectCore({
              defaultValue: "",
              options: buildOptions(table.columns),
              onChange(v) {
                const column = table.columns.find((c) => c.name === v);
                $sub.column = column || null;
                handleColumnSelectValueChange(v, [x, y + 1], $sub, table.columns);
              },
            }),
          });
          _values[x] = [
            ..._values[x].slice(0, y + 1),
            // ...
            $sub,
          ];
          return;
        }
        return;
      }
      if (column) {
        if (column.type === "datetime") {
          _values[x] = [
            ..._values[x].slice(0, y + 1),
            new FilterInput({
              type: "condition",
              $input: new SelectCore({
                defaultValue: "<",
                options: [
                  {
                    label: "之前",
                    value: "<",
                  },
                  {
                    label: "之后",
                    value: ">",
                  },
                  {
                    label: "等于",
                    value: "=",
                  },
                  // {
                  //   label: "在这之间",
                  //   value: "BETWEEN",
                  // },
                ],
              }),
            }),
            new FilterInput({
              type: "value",
              $input: new InputCore({
                defaultValue: new Date(),
                type: "datetime-local",
                onEnter() {
                  $submit.click();
                },
              }),
            }),
          ];
          return;
        }
        if (column.type === "integer") {
          _values[x] = [
            ..._values[x].slice(0, y + 1),
            new FilterInput({
              type: "condition",
              $input: new SelectCore({
                defaultValue: "<",
                options: [
                  {
                    label: "小于",
                    value: "<",
                  },
                  {
                    label: "大于",
                    value: ">",
                  },
                  {
                    label: "等于",
                    value: "=",
                  },
                  // {
                  //   label: "IN",
                  //   value: "IN",
                  // },
                  // BETWEEN 5 AND 10;
                  // {
                  //   label: "之间",
                  //   value: "BETWEEN",
                  // },
                ],
              }),
            }),
            new FilterInput({
              type: "value",
              $input: new InputCore({
                defaultValue: 0,
                onEnter() {
                  $submit.click();
                },
              }),
            }),
          ];
          return;
        }
      }
      _values[x] = [
        ..._values[x].slice(0, y + 1),
        new FilterInput({
          type: "condition",
          $input: new SelectCore({
            defaultValue: "=",
            options: [
              {
                label: "等于",
                value: "=",
              },
              {
                label: "包含",
                value: "LIKE",
              },
              {
                label: "不等于",
                value: "!=",
              },
              // {
              //   label: "不包含",
              //   value: "NOT LIKE",
              // },
              // 这个应该判断字段是否支持为空
              //     {
              //       label: "IS NULL",
              //       value: "IS NULL",
              //     },
            ],
          }),
        }),
        new FilterInput({
          type: "value",
          $input: new InputCore({
            defaultValue: "",
            onEnter() {
              $submit.click();
            },
          }),
        }),
      ];
      return;
    })();
    emitter.emit(Events.Change, [..._values]);
  }
  function buildOptions(columns: TableColumn[]) {
    const options: { label: string; value: string }[] = [];
    for (let j = 0; j < columns.length; j += 1) {
      const column = columns[j];
      const { name, type, references } = column;
      (() => {
        if (type === "index") {
          return;
        }
        const base = {
          label: name,
          value: name,
        };
        if (references) {
          options.push(
            ...[
              base,
              {
                label: references,
                value: references,
              },
            ]
          );
          return;
        }
        options.push(base);
      })();
    }
    // console.log("build options", options);
    return options;
  }
  function initValues() {
    _values = [];
    const rowIndex = _values.length;
    const prefix = new FilterInput({
      type: "multiple",
      $input: new PrefixTag({
        value: "WHERE",
      }),
    });
    const first = new FilterInput({
      type: "field",
      $input: new SelectCore<string>({
        defaultValue: "",
        options: buildOptions(_baseColumns),
        onChange(v) {
          const column = _table.columns.find((c) => c.name === v);
          first.column = column || null;
          console.log(_values[rowIndex].map((input) => input.type));
          _values[rowIndex] = _values[rowIndex].slice(0, 2);
          handleColumnSelectValueChange(v, [rowIndex, 1], first, _baseColumns);
        },
      }),
    });
    _values = [[prefix, first]];
  }

  enum Events {
    Change,
  }
  type TheTypesOfEvents = {
    [Events.Change]: typeof _values;
  };

  const emitter = base<TheTypesOfEvents>();
  return {
    get values() {
      return _values;
    },
    $submit,
    $reset,
    $more,
    setTable(table: TableWithColumns) {
      _table = table;
    },
    setTables(tables: TableWithColumns[]) {
      _tables = tables;
    },
    setOptions(columns: TableColumnCore[]) {
      _baseColumns = columns;
      initValues();
      emitter.emit(Events.Change, [..._values]);
    },
    buildSQL(table: TableWithColumns, tables: TableWithColumns[], response: { pageSize: number; page: number }) {
      // console.log("[BIZ]filter/index - buildSQL", table, tables, _values);
      const partSQL = buildQuerySQL(table, _values, tables);
      const pagination = ` LIMIT ${response.pageSize} OFFSET ${(response.page - 1) * response.pageSize}`;
      const sql = partSQL + pagination + ";";
      return sql;
    },
    onChange(handler: Handler<TheTypesOfEvents[Events.Change]>) {
      emitter.on(Events.Change, handler);
    },
  };
}
