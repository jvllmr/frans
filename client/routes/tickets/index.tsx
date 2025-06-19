import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/tickets/')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/ticket/"!</div>
}
