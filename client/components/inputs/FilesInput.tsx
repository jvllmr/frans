import { Center, Flex, Grid, Paper, Stack, Text } from "@mantine/core";
import { Dropzone, FileWithPath } from "@mantine/dropzone";
import { IconPlus } from "@tabler/icons-react";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useFileSizeFormatter } from "~/i18n";
import { FileIcon } from "../FileIcon";

function FileCard({ children }: { children: React.ReactNode }) {
  return (
    <Paper withBorder p="sm" w="100%" h="100%">
      {children}
    </Paper>
  );
}

export interface FilesInputProps {
  value?: FileWithPath[];
  onChange: (value: FileWithPath[]) => void;
}

export function FilesInput({ onChange, value }: FilesInputProps) {
  const { t } = useTranslation("file_input");
  const fileSizeFormatter = useFileSizeFormatter();
  const [filesCache, setFilesCache] = useState<FileWithPath[]>([]);
  useEffect(() => {
    if (value && value?.length < filesCache.length) {
      setFilesCache(value);
    }
  }, [value, filesCache]);

  return (
    <Dropzone
      maxSize={window.fransMaxSizes}
      maxFiles={window.fransMaxFiles}
      onDrop={(newFiles) => {
        const files = [...filesCache, ...newFiles];
        onChange(files);
        setFilesCache(files);
      }}
    >
      <Paper withBorder mih={150} p="xl">
        <Dropzone.Idle>
          {value?.length === 0 ? (
            <Flex w="100%" h="100%" justify="center" align="center">
              <b>{t("drop_file")}</b>
            </Flex>
          ) : null}
        </Dropzone.Idle>
        <Grid gutter="xl">
          {value?.map((fileWithPath) => (
            <Grid.Col span={3} key={fileWithPath.path}>
              <FileCard>
                <Stack gap={5}>
                  <Center>
                    <FileIcon size={32} filename={fileWithPath.name} />
                  </Center>
                  <Center>
                    <Text size="xs">
                      {fileWithPath.name} (
                      {fileSizeFormatter(fileWithPath.size)})
                    </Text>
                  </Center>
                </Stack>
              </FileCard>
            </Grid.Col>
          ))}
          {(value?.length || 0) > 0 &&
          (value?.length || 0) < window.fransMaxSizes ? (
            <Grid.Col span={3}>
              <FileCard>
                <Flex w="100%" h="100%" justify="center" align="center">
                  <IconPlus />
                </Flex>
              </FileCard>
            </Grid.Col>
          ) : null}
        </Grid>
      </Paper>
    </Dropzone>
  );
}
