const openButton = document.getElementById("open-playlist");
const copyButton = document.getElementById("copy-titles");
const status = document.getElementById("playlist-status");
const titleList = document.getElementById("titles");
let currentTitles = [];

copyButton.addEventListener("click", async () => {
  try {
    await window.go.main.App.CopyTitles(currentTitles);
  } catch (error) {
    status.textContent = String(error);
  }
});

openButton.addEventListener("click", async () => {
  openButton.disabled = true;
  try {
    const result = await window.go.main.App.OpenPlaylist();
    if (result.cancelled) return;

    titleList.replaceChildren();
    currentTitles = [];
    copyButton.disabled = true;
    if (result.error) {
      status.textContent = result.error;
      return;
    }

    currentTitles = result.titles ?? [];
    status.textContent = currentTitles.length === 0 ? "No titles found" : result.name;
    copyButton.disabled = currentTitles.length === 0;
    for (const title of currentTitles) {
      const item = document.createElement("li");
      item.textContent = title;
      titleList.append(item);
    }
  } catch (error) {
    titleList.replaceChildren();
    currentTitles = [];
    copyButton.disabled = true;
    status.textContent = String(error);
  } finally {
    openButton.disabled = false;
  }
});
