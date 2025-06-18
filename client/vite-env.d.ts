/// <reference types="vite/client" />

declare global {
  interface Window {
    fransRootPath: string;
    __fransAssetUrl: (filename: string) => string;
  }
}
