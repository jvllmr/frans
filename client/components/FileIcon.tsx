import { Loader } from "@mantine/core";

import {
  IconFile,
  IconFileDatabase,
  IconFileSettings,
  IconFileSpreadsheet,
  IconFileTypeBmp,
  IconFileTypeCss,
  IconFileTypeCsv,
  IconFileTypeDoc,
  IconFileTypeDocx,
  IconFileTypeHtml,
  IconFileTypeJpg,
  IconFileTypeJs,
  IconFileTypeJsx,
  IconFileTypePdf,
  IconFileTypePhp,
  IconFileTypePng,
  IconFileTypePpt,
  IconFileTypeRs,
  IconFileTypeSql,
  IconFileTypeSvg,
  IconFileTypeTs,
  IconFileTypeTsx,
  IconFileTypeTxt,
  IconFileTypeVue,
  IconFileTypeXls,
  IconFileTypeXml,
  IconFileZip,
} from "@tabler/icons-react";
import React, { Suspense } from "react";
const iconsMap: Record<string, React.FC<{ size?: number }>> = {
  bmp: IconFileTypeBmp,
  css: IconFileTypeCss,
  csv: IconFileTypeCsv,
  db: IconFileDatabase,
  doc: IconFileTypeDoc,
  docx: IconFileTypeDocx,
  gz: IconFileZip,
  html: IconFileTypeHtml,
  jpg: IconFileTypeJpg,
  js: IconFileTypeJs,
  json: IconFileSettings,
  jsx: IconFileTypeJsx,
  pdf: IconFileTypePdf,
  php: IconFileTypePhp,
  png: IconFileTypePng,
  ppt: IconFileTypePpt,
  rs: IconFileTypeRs,
  sql: IconFileTypeSql,
  svg: IconFileTypeSvg,
  tar: IconFileZip,
  ts: IconFileTypeTs,
  tsx: IconFileTypeTsx,
  txt: IconFileTypeTxt,
  vue: IconFileTypeVue,
  xls: IconFileTypeXls,
  xml: IconFileTypeXml,
  xlsx: IconFileSpreadsheet,
  yaml: IconFileSettings,
  zip: IconFileZip,
  "7zip": IconFileZip,
};

const extensionRegex = /.+\.(.+)$/;
export function FileIcon({
  filename,
  size,
}: {
  filename: string;
  size?: number;
}) {
  const lowered = filename.toLocaleLowerCase();
  extensionRegex.exec(lowered);
  const IconComponent =
    iconsMap[extensionRegex.exec(lowered)?.[1] ?? ""] ?? IconFile;

  return (
    <Suspense fallback={<Loader size="sm" />}>
      <IconComponent size={size} />
    </Suspense>
  );
}
