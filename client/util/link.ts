export function getShareLink(shareId: string) {
  const windowLocation = new URL(window.location.href);
  windowLocation.pathname = `${window.fransRootPath}/s/${shareId}`;
  return windowLocation.toString();
}

export function getInternalFileLink(fileId: string) {
  return `${window.fransRootPath}/api/v1/file/${fileId}`;
}
