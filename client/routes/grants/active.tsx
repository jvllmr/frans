import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/grants/active')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/grants/active"!</div>
}
