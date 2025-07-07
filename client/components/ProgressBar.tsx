import { Box, Center, Flex, Progress, Stack } from "@mantine/core";
import { useFileSizeFormatter } from "~/i18n";
import { ProgressState } from "~/util/progress";

function secondsToTime(seconds: number) {
  const hrs = Math.floor(seconds / 3600);
  const mins = Math.floor((seconds % 3600) / 60);
  const secs = Math.floor(seconds % 60);

  const formattedHrs = hrs.toString().padStart(2, "0");
  const formattedMins = mins.toString().padStart(2, "0");
  const formattedSecs = secs.toString().padStart(2, "0");

  return `${formattedHrs}:${formattedMins}:${formattedSecs}`;
}

export interface ProgressBarProps {
  state: ProgressState;
}

export function ProgressBar({
  state: { percentage, speed, estimatedSeconds },
}: ProgressBarProps) {
  const fileSizeFormatter = useFileSizeFormatter();
  if (!percentage && !speed && !estimatedSeconds) return null;
  return (
    <Stack>
      <Center>
        <Progress value={percentage * 100} animated w="70%" />
        <Box ml="sm">{Math.floor(percentage * 100)}%</Box>
      </Center>
      <Flex justify="space-evenly" align="center">
        <Box p="xs" w="15%">
          {secondsToTime(estimatedSeconds)}
        </Box>
        <Box p="xs" w="15%">
          {fileSizeFormatter(speed)}/s
        </Box>
      </Flex>
    </Stack>
  );
}
