import { TextInput, TextInputProps } from "@mantine/core";
import React, { useMemo } from "react";

export interface NullTextInputProps
  extends Omit<TextInputProps, "value" | "onChange"> {
  value?: string | null;
  onChange?: (value: string | null) => void;
}

export const NullTextInput = React.forwardRef<
  HTMLInputElement,
  NullTextInputProps
>(function NullTextInput({ value, onChange, ...props }, ref) {
  const fixedValue = useMemo(() => (value === null ? "" : value), [value]);
  const onChangeWrapper:
    | React.ChangeEventHandler<HTMLInputElement>
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
    <TextInput
      {...props}
      ref={ref}
      value={fixedValue}
      onChange={onChangeWrapper}
    />
  );
});
