const openButton = document.getElementById("open-playlist");
const copyButton = document.getElementById("copy-titles");
const status = document.getElementById("playlist-status");
const titleList = document.getElementById("titles");
const sanitizeButton = document.getElementById("sanitize-mp3");
const sanitizeStatus = document.getElementById("sanitize-status");
const sanitizeErrors = document.getElementById("sanitize-errors");
let currentTitles = [];

window.runtime.EventsOn("sanitize-progress", (message) => {
  sanitizeStatus.textContent = message;
});

sanitizeButton.addEventListener("click", async () => {
  sanitizeButton.disabled = true;
  try {
    const result = await window.go.main.App.SanitizeMP3();
    if (result.cancelled) return;
    sanitizeErrors.replaceChildren();
    if (result.error) {
      sanitizeStatus.textContent = result.error;
      return;
    }
    const failures = result.failures ?? [];
    sanitizeStatus.textContent = `Updated: ${result.updated} · Unchanged: ${result.unchanged} · Failed: ${failures.length}`;
    for (const failure of failures.slice(0, 10)) {
      const item = document.createElement("li");
      item.textContent = `${failure.name}: ${failure.error}`;
      sanitizeErrors.append(item);
    }
    if (failures.length > 10) {
      const item = document.createElement("li");
      item.textContent = `…and ${failures.length - 10} more errors`;
      sanitizeErrors.append(item);
    }
  } catch (error) {
    sanitizeStatus.textContent = String(error);
  } finally {
    sanitizeButton.disabled = false;
  }
});

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
