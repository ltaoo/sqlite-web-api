import { request } from "@/biz/requests";
import { TmpRequestResp } from "@/domains/request/utils";
import { Result } from "@/domains/result";
import { TableColumnType } from "@/domains/ui/table/column";

export function fetchTableList() {
  return request.post<
    {
      name: string;
      columns: {
        index: number;
        name: string;
        type: TableColumnType;
        not_null: boolean;
        primary_key: boolean;
        foreign_key: boolean;
        auto_increment: boolean;
        unique: boolean;
        references: string;
      }[];
    }[]
  >("/api/v1/database/tables");
}
export function fetchTableListProcess(r: TmpRequestResp<typeof fetchTableList>) {
  if (r.error) {
    return Result.Err(r.error.message);
  }
  const data: {
    name: string;
    columns: {
      name: string;
      type: TableColumnType;
      not_null: number;
      is_primary_key: number;
      width: number;
      references?: string;
    }[];
  }[] = [];
  for (let i = 0; i < r.data.length; i += 1) {
    const table = r.data[i];
    const columns: (typeof data)[number]["columns"] = [];
    for (let j = 0; j < table.columns.length; j += 1) {
      const column = table.columns[j];
      const { name, type, not_null, primary_key, references } = column;
      (() => {
        const base = {
          name,
          type,
          not_null: Number(not_null),
          is_primary_key: Number(primary_key),
          width: 200,
          references,
        };
        if (references) {
          columns.push(
            ...[
              base,
              //       {
              //         name: references,
              //         type: "table" as TableColumnType,
              //         not_null: base.not_null,
              //         is_primary_key: 0,
              //         width: 200,
              //       },
            ]
          );
          return;
        }
        columns.push(base);
      })();
    }
    data.push({
      name: table.name,
      columns,
    });
  }
  return Result.Ok(
    data.sort((a, b) => {
      return a.name.toLowerCase().charCodeAt(0) - b.name.toLowerCase().charCodeAt(0);
    })
  );
}
export function execQueryRaw(params: { query: string }) {
  return request.post<null | string[][]>("/api/v1/database/exec", params);
}
