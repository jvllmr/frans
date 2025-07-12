import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useEffect } from "react";

export const Route = createFileRoute("/s/$shareId")({
  component: RouteComponent,
});

function RouteComponent() {
  const navigate = useNavigate();
  const shareId = Route.useParams({ select: (p) => p.shareId });
  useEffect(() => {
    fetch(`${window.fransRootPath}/s/${shareId}`).then((resp) => {
      const nextLocation = new URL(resp.url).pathname;

      navigate({ to: nextLocation });
    });
  }, [navigate, shareId]);

  return "Redirecting...";
}
