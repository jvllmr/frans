import { z } from "zod/v4";

export const fileSchema = z.object({
  sha256sum: z.string(),
  name: z.string(),
  size: z.int(),
});
