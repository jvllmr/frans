import { Textarea, TextareaProps } from "@mantine/core";
import React, { useMemo } from "react";

export interface NullTextareaProps
  extends Omit<TextareaProps, "value" | "onChange"> {
  value?: string | null;
  onChange?: (value: string | null) => void;
}

export const NullTextarea = React.forwardRef<
  HTMLTextAreaElement,
  NullTextareaProps
>(function NullTextarea({ value, onChange, ...props }, ref) {
  const fixedValue = useMemo(() => (value === null ? "" : value), [value]);
  const onChangeWrapper:
    | React.ChangeEventHandler<HTMLTextAreaElement>
    | undefined = useMemo(
    () =>
      onChange
        ? (e) => {
            const newValue = e.target.value;
            onChange(newValue === "" ? null : newValue);
          }
        : undefined,
    [onChange],
  );

  return (
    <Textarea
      {...props}
      ref={ref}
      value={fixedValue}
      onChange={onChangeWrapper}
    />
  );
});
