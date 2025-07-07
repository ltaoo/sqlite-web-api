import { describe, it, expect } from "vitest";

import { InputCore, SelectCore, TabsCore } from "@/domains/ui";
import { TableCore } from "@/domains/ui/table/table";
import { TableColumnType } from "@/domains/ui/table/column";

import { FilterInput, PrefixTag } from "./index";
import { buildORMObject, buildQuerySQL } from "./utils";

describe("some special case", () => {
  it("only search by foreign key", () => {
    const rows = [
      [
        new FilterInput({
          type: "multiple",
          $input: new PrefixTag({ value: "WHERE" }),
        }),
        new FilterInput({
          type: "field",
          column: {
            name: "media_id",
            type: TableColumnType.Text,
            references: "Media",
          },
          $input: new SelectCore({ defaultValue: "media_id" }),
        }),
        new FilterInput({
          type: "condition",
          $input: new SelectCore({ defaultValue: "=" }),
        }),
        new FilterInput({
          type: "value",
          $input: new InputCore({ defaultValue: "mp3F4nNgc5vIdzx" }),
        }),
      ],
    ];
    const orm = buildORMObject(
      {
        name: "PlayHistoryV2",
        columns: [
          {
            name: "id",
            type: TableColumnType.Text,
          },
          {
            name: "media_id",
            type: TableColumnType.Text,
            references: "Media",
          },
        ],
      },
      rows,
      [
        {
          name: "Media",
          columns: [
            {
              name: "id",
              type: TableColumnType.Text,
              is_primary_key: 1,
            },
          ],
        },
      ]
    );
    //     expect(sql).toBe(
    //       "SELECT `PlayHistoryV2`.* FROM `PlayHistoryV2` WHERE `PlayHistoryV2`.`media_id` = 'mp3F4nNgc5vIdzx'"
    //     );
    expect(orm).toStrictEqual({
      name: "PlayHistoryV2",
      where: {
        media_id: {
          equals: "mp3F4nNgc5vIdzx",
        },
      },
    });
  });
});
