import { For, createSignal } from "solid-js";

import { TabHeaderCore } from "@/domains/ui/tab-header";
import { cn } from "@/utils";

export const TabHeader = (props: { store: TabHeaderCore<any> }) => {
  const { store } = props;

  const [state, setState] = createSignal(store.state);
  const [left, setLeft] = createSignal<null | number>(null);

  const { tabs: options, current } = state();

  store.onStateChange((v) => {
    setState(v);
  });
  store.onLinePositionChange((v) => {
    setLeft(v.left);
  });

  return (
    <div
      class={cn("__a tabs w-full overflow-x-auto scroll--hidden")}
      //       style="{{style}}"
      onAnimationStart={(event) => {
        const { width, height, left } = event.currentTarget.getBoundingClientRect();
        store.updateContainerClient({ width, height, left });
      }}
    >
      <div
        class="tabs-wrapper relative"
        // scroll-with-animation="{{scrollWithAnimation}}"
        // scroll-left="{{scrollLeftInset}}"
        // scroll-x
      >
        <div id="tabs-wrapper" class="flex">
          <For each={options}>
            {(tab, index) => {
              return (
                <div
                  classList={{
                    "__a p-4 break-keep cursor-pointer": true,
                  }}
                  // style="{{current === index ? activeItemStyle : itemStyle}}"
                  onClick={() => {
                    store.select(index());
                  }}
                  onAnimationStart={(event) => {
                    event.stopPropagation();
                    const { width, height, left } = event.currentTarget.getBoundingClientRect();
                    store.updateTabClient(index(), { width, height, left });
                  }}
                >
                  {tab.text}
                </div>
              );
            }}
          </For>
          {left() !== null ? (
            <div
              class="absolute bottom-0 w-4 bg-slate-900 transition-all"
              style={{
                left: `${left()}px`,
                height: "4px",
                transform: "translateX(-50%)",
              }}
            />
          ) : null}
        </div>
      </div>
    </div>
  );
};
