export function getShareLink(shareId: string) {
  const windowLocation = new URL(window.location.href);
  windowLocation.pathname = `${window.fransRootPath}/s/${shareId}`;
  return windowLocation.toString();
}

export function getInternalFileLink(fileId: string, addDownload?: boolean) {
  const baseLink = `${window.fransRootPath}/api/v1/file/${fileId}`;
  return addDownload
    ? baseLink + "?" + new URLSearchParams({ addDownload: "1" }).toString()
    : baseLink;
}
