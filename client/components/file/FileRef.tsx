import { Anchor } from "@mantine/core";
import { useFileSizeFormatter } from "~/i18n";

interface FileRefFile {
  name: string;
}

interface FileRefTextWithoutSizeProps {
  file: FileRefFile;
  withoutSize: true;
}

function FileRefTextWithoutSize({
  file,
}: Omit<FileRefTextWithoutSizeProps, "withoutSize">) {
  return file.name;
}

interface FileRefFileWithSize {
  name: string;
  size: number;
}

interface FileRefTextWithSizeProps {
  file: FileRefFileWithSize;
  withoutSize?: false;
}

function FileRefFileWithSize({ file }: FileRefTextWithSizeProps) {
  const fileSizeFormatter = useFileSizeFormatter();
  return (
    <>
      <FileRefTextWithoutSize file={file} /> ({fileSizeFormatter(file.size)})
    </>
  );
}

export type FileRefTextProps =
  | FileRefTextWithSizeProps
  | FileRefTextWithoutSizeProps;

export function FileRefText({ file, withoutSize }: FileRefTextProps) {
  return withoutSize ? (
    <FileRefTextWithoutSize file={file} />
  ) : (
    <FileRefFileWithSize file={file} />
  );
}

export type FileRefProps = FileRefTextProps & {
  link: string;
  onClick?: () => void;
};

export function FileRef({ link, onClick, ...props }: FileRefProps) {
  return (
    <Anchor href={link} onClick={onClick}>
      <FileRefText {...props} />
    </Anchor>
  );
}
