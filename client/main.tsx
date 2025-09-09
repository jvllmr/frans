import { RouterProvider, createRouter } from "@tanstack/react-router";
import { StrictMode } from "react";
import ReactDOM from "react-dom/client";
import { queryClient } from "~/api";
import { PendingComponent } from "./components/routing/PendingComponent";
import i18n from "./i18n";
import { routeTree } from "./routeTree.gen";

const router = createRouter({
  routeTree,
  basepath: window.fransRootPath,
  defaultPendingComponent: PendingComponent,
  defaultPreloadStaleTime: 0,
  context: {
    queryClient,
    i18n,
  },
});

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

const rootElement = document.getElementById("root")!;
if (!rootElement.innerHTML) {
  const root = ReactDOM.createRoot(rootElement);
  root.render(
    <StrictMode>
      <RouterProvider router={router} />
    </StrictMode>,
  );
}
