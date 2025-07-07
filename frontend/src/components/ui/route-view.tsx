/**
 * @file 会销毁页面的视图（如果希望不销毁可以使用 keep-alive-route-view
 */
import { Show, createSignal, JSX, onCleanup } from "solid-js";

import { RouteViewCore } from "@/domains/route_view";

export function RouteView(props: { store: RouteViewCore; index: number } & JSX.HTMLAttributes<HTMLDivElement>) {
  const { store, index } = props;

  const [state, setState] = createSignal(store.state);

  const visible = () => state().visible;
  const mounted = () => state().mounted;

  store.onStateChange((nextState) => {
    setState(nextState);
  });

  return (
    <Show when={mounted()}>
      <div
        class={props.class}
        style={{
          "z-index": index,
        }}
        data-state={visible() ? "open" : "closed"}
      >
        {props.children}
      </div>
    </Show>
  );
}
