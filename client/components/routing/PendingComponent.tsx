import { Flex, Loader } from "@mantine/core";

export function PendingComponent() {
  return (
    <Flex my="xl" w="100%" justify="center" align="center">
      <Loader type="dots" />
    </Flex>
  );
}
