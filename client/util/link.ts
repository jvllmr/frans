export function getShareLink(shareId: string) {
  const windowLocation = new URL(window.location.href);
  windowLocation.pathname = `${window.fransRootPath}/s/${shareId}`;
  return windowLocation.toString();
}
