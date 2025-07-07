/**
 * @file 控制内容显隐的组件
 */
import { createSignal, JSX } from "solid-js";

import { Show } from "@/packages/ui/show";
import { PresenceCore } from "@/domains/ui/presence";
import { cn } from "@/utils/index";

export const Presence = (
  props: {
    store: PresenceCore;
    enterClassName?: string;
    exitClassName?: string;
  } & JSX.HTMLAttributes<HTMLElement>
) => {
  const { store, enterClassName, exitClassName } = props;

  const [state, setState] = createSignal(store.state);

  store.onStateChange((v) => setState(v));

  // const { visible, mounted, text } = state;

  return (
    <Show when={state().mounted}>
      <div
        class={cn(
          "presence",
          state().enter && enterClassName ? enterClassName : "",
          state().exit && exitClassName ? exitClassName : "",
          props.class
        )}
        role="presentation"
        style={{ display: state().visible ? "block" : "none" }}
        // data-state={visible ? "open" : "closed"}
        // onAnimationEnd={() => {
        //   store.unmount();
        // }}
      >
        {props.children}
      </div>
    </Show>
  );
};
