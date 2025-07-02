/// <reference types="vite/client" />

interface Window {
  fransRootPath: string;
  __fransAssetUrl: (filename: string) => string;
  fransMaxFiles: number;
  fransMaxSizes: number;
  fransDefaultExpiryTotalDays: number;
  fransDefaultExpiryTotalDownloads: number;
  fransDefaultExpiryDaysSinceLastDownload: number;
}

declare module "virtual:*";
