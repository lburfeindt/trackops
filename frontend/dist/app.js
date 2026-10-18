const openButton = document.getElementById("open-playlist");
const status = document.getElementById("playlist-status");
const titleList = document.getElementById("titles");

openButton.addEventListener("click", async () => {
  openButton.disabled = true;
  try {
    const result = await window.go.main.App.OpenPlaylist();
    if (result.cancelled) return;

    titleList.replaceChildren();
    if (result.error) {
      status.textContent = result.error;
      return;
    }

    status.textContent = result.name;
    for (const title of result.titles ?? []) {
      const item = document.createElement("li");
      item.textContent = title;
      titleList.append(item);
    }
  } catch (error) {
    status.textContent = String(error);
  } finally {
    openButton.disabled = false;
  }
});
