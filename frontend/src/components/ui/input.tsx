import { JSX, createSignal, onMount } from "solid-js";
import { Loader2 } from "lucide-solid";

import { InputCore, InputCoreProps } from "@/domains/ui/input";
import { connect } from "@/domains/ui/input/connect.web";
import { Input as InputPrimitive } from "@/packages/ui/input";
import { cn } from "@/utils/index";

const Input = (props: { store: InputCore<any>; prefix?: JSX.Element; class?: string }) => {
  const { store, prefix } = props;

  const [state, setState] = createSignal(store.state);
  store.onStateChange((nextState) => {
    setState(nextState);
  });

  onMount(() => {
    store.setMounted();
  });

  return (
    <div class="relative w-full">
      <div class="absolute left-3 top-[50%] translate-y-[-50%] text-slate-400 ">
        {(() => {
          if (!prefix) {
            return null;
          }
          if (state().loading) {
            return <Loader2 class="w-4 h-4 animate-spin" />;
          }
          return prefix;
        })()}
      </div>
      <InputPrimitive
        class={cn(
          "flex items-center w-full px-3",
          "focus:outline-none focus:ring-2 focus:ring-slate-400 focus:ring-offset-2",
          "disabled:cursor-not-allowed disabled:opacity-50",
          "placeholder:text-slate-400",
          prefix ? "pl-8" : "",
          props.class
        )}
        style={{
          "vertical-align": "bottom",
        }}
        store={store}
      />
    </div>
  );
};

export const FreeInput = (
  props: InputCoreProps<any> & { autoFocus?: boolean; prefix?: JSX.Element } & JSX.HTMLAttributes<HTMLInputElement>
) => {
  const { prefix, class: className, ...rest } = props;

  const store = new InputCore({
    ...rest,
  });
  const [state, setState] = createSignal(store.state);
  store.onStateChange((nextState) => {
    setState(nextState);
  });

  return (
    <div class="relative w-full">
      <div class="absolute left-3 top-[50%] translate-y-[-50%] text-slate-400 ">
        {(() => {
          if (!prefix) {
            return null;
          }
          if (state().loading) {
            return <Loader2 class="w-4 h-4 animate-spin" />;
          }
          return prefix;
        })()}
      </div>
      <InputPrimitive
        class={cn(
          "flex items-center w-full px-3",
          "focus:outline-none focus:ring-2 focus:ring-slate-400 focus:ring-offset-2",
          "disabled:cursor-not-allowed disabled:opacity-50",
          "placeholder:text-slate-400",
          prefix ? "pl-8" : "",
          props.class
        )}
        style={{
          "vertical-align": "bottom",
        }}
        store={store}
      />
    </div>
  );
};

export { Input };
