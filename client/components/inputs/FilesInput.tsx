import { Center, Grid, Paper, Stack, Text } from "@mantine/core";
import { Dropzone, FileWithPath } from "@mantine/dropzone";
import { FileIcon } from "../FileIcon";

export interface FilesInputProps {
  value?: FileWithPath[];
  onChange: (value: FileWithPath[]) => void;
}

export function FilesInput({ onChange, value }: FilesInputProps) {
  return (
    <Dropzone
      maxSize={window.fransMaxSizes}
      maxFiles={window.fransMaxFiles}
      onDrop={onChange}
    >
      <Paper withBorder mih={150} p="xl">
        <Dropzone.Accept>Yay</Dropzone.Accept>
        <Grid gutter="xl">
          {value?.map((fileWithPath) => (
            <Grid.Col span={3} key={fileWithPath.path}>
              <Paper withBorder p="sm">
                <Stack gap={5}>
                  <Center>
                    <FileIcon size={32} filename={fileWithPath.name} />
                  </Center>
                  <Center>
                    <Text size="xs">{fileWithPath.name}</Text>
                  </Center>
                </Stack>
              </Paper>
            </Grid.Col>
          ))}
        </Grid>
      </Paper>
    </Dropzone>
  );
}
