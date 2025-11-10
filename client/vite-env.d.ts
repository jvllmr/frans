/// <reference types="vite/client" />

interface Window {
  fransRootPath: string;
  __fransAssetUrl: (filename: string) => string;
  fransMaxFiles: number;
  fransMaxSizes: number;
  fransDefaultExpiryTotalDays: number;
  fransDefaultExpiryTotalDownloads: number;
  fransDefaultExpiryDaysSinceLastDownload: number;
  fransGrantDefaultExpiryTotalDays: number;
  fransGrantDefaultExpiryTotalUploads: number;
  fransGrantDefaultExpiryDaysSinceLastUpload: number;
  fransColor: string;
  fransCustomColor: string;

  fransVersion: string;
  fransTitle: string;
}

declare module "virtual:*";
