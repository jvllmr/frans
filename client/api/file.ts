import { z } from "zod/v4";

export const fileSchema = z.object({
  sha512: z.string(),
  name: z.string(),
  size: z.int(),
});
