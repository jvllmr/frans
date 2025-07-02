/* eslint-disable @typescript-eslint/no-explicit-any */
import { Code } from "@mantine/core";
import { UseFormReturnType } from "@mantine/form";

function DebugInfo({ values }: { values: any }) {
  return import.meta.env.DEV ? (
    <Code mt="lg" block>
      {JSON.stringify(values, undefined, 2)}
    </Code>
  ) : null;
}

export function FormDebugInfo({ form }: { form: UseFormReturnType<any> }) {
  return <DebugInfo values={form.values} />;
}
