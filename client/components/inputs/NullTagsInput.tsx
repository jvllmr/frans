import { TagsInput, TagsInputProps } from "@mantine/core";
import React, { useMemo } from "react";

export interface NullTagsInputProps extends Omit<
  TagsInputProps,
  "value" | "onChange"
> {
  value?: string[] | null;
  onChange?: (value: string[] | null) => void;
}

export const NullTagsInput = React.forwardRef<
  HTMLInputElement,
  NullTagsInputProps
>(function NullTagsInput({ value, onChange, ...props }, ref) {
  const fixedValue = useMemo(() => (value === null ? [] : value), [value]);
  const onChangeWrapper: TagsInputProps["onChange"] | undefined = useMemo(
    () =>
      onChange
        ? (value) => {
            onChange(value.length > 0 ? value : null);
          }
        : undefined,
    [onChange],
  );

  return (
    <TagsInput
      {...props}
      ref={ref}
      value={fixedValue}
      onChange={onChangeWrapper}
    />
  );
});
