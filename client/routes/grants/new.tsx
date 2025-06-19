import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/grants/new")({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/grants/new"!</div>;
}
