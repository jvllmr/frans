/// <reference types="vite/client" />

interface Window {
  fransRootPath: string;
  __fransAssetUrl: (filename: string) => string;
}

declare module "virtual:*";
