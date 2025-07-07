import { JSX } from "solid-js/jsx-runtime";

import { RouteViewCore } from "@/domains/route_view";
import { ScrollViewCore } from "@/domains/ui";
import { StorageCore } from "@/domains/storage";
import { HttpClientCore } from "@/domains/http_client";

import { app, history } from ".";
import { PageKeys } from "./routes";
import { storage } from "./storage";

export type GlobalStorageValues = (typeof storage)["values"];
export type ViewComponentProps = {
  app: typeof app;
  history: typeof history;
  client: HttpClientCore;
  view: RouteViewCore;
  storage: StorageCore<GlobalStorageValues>;
  pages: Omit<Record<PageKeys, ViewComponent>, "root">;
  parent?: {
    view: RouteViewCore;
    scrollView?: ScrollViewCore;
  };
};
export type ViewComponent = (props: ViewComponentProps) => JSX.Element;
