import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/grants/")({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/grants/"!</div>;
}
