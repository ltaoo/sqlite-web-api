import { ViewComponent } from "@/store/types";
import { NotFoundPage } from "@/pages/notfound";
import { SqliteDatabasePage } from "@/pages/home";

import { PageKeys } from "./routes";

export const pages: Omit<Record<PageKeys, ViewComponent>, "root"> = {
  "root.home": SqliteDatabasePage,
  "root.notfound": NotFoundPage,
};
